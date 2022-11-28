package local

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gopherlearning/gophkeeper/internal/model"
	proto "github.com/gopherlearning/gophkeeper/pkg/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	ErrNilCryptor = errors.New("empty cryptor")
)

type Storage struct {
	cancelFunc   context.CancelFunc
	db           *badger.DB
	remoteStatus atomic.Value
	remote       proto.PublicClient
	serverURL    atomic.Value
	owner        string
}

func NewLocalStorage(secret, path, serverURL string) (*Storage, error) {
	db, err := badger.Open(badger.DefaultOptions(path).
		WithEncryptionKey([]byte(secret)).
		WithIndexCacheSize(20 * 1024 * 1024).
		WithSyncWrites(true).WithLogger(nil))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		db.Close()
	}()

	if len(serverURL) != 0 {
		err = db.Update(func(txn *badger.Txn) error {
			return txn.SetEntry(&badger.Entry{Key: []byte("serverURL"), Value: []byte(serverURL)})
		})
		if err != nil {
			cancel()
			return nil, err
		}
	}

	fmt.Println(hex.EncodeToString(sha256.New().Sum([]byte(secret))))

	store := &Storage{db: db, cancelFunc: cancel, owner: hex.EncodeToString(sha256.New().Sum([]byte(secret)))}
	err = store.getServerURL()

	if err != nil && err != badger.ErrKeyNotFound {
		cancel()
		return nil, err
	}

	if store.serverURL.Load() == nil {
		store.serverURL.Store("")
	}

	log.Debug().Msg(store.serverURL.Load().(string))

	if len(store.serverURL.Load().(string)) != 0 {
		store.connectToServer(ctx)

		var updated uint64

		err = store.db.View(func(txn *badger.Txn) error {
			var i *badger.Item
			i, err = txn.Get([]byte("lastUpdated"))
			if err != nil {
				return err
			}

			return i.Value(func(val []byte) error {
				updated, err = strconv.ParseUint(string(val), 10, 64)
				if err != nil {
					return err
				}
				return nil
			})
		})
		if err != nil {
			log.Debug().Err(err)
		}

		go store.getUpdatesFromServer(ctx, updated)
	}

	// }
	// if len(serverURL) != 0 {
	// go func() {
	// 	store.remoteStatus.Store("❌ >")

	// 	for {
	// 		time.Sleep(time.Second * 2)
	// 		v := store.remoteStatus.Load().(string)
	// 		switch v {
	// 		case "✅ >":
	// 			store.remoteStatus.Store("❌ >")
	// 		case "❌ >":
	// 			store.remoteStatus.Store("✅ >")
	// 		}
	// 	}
	// }()
	// }

	return store, nil
}

func (s *Storage) getServerURL() error {
	return s.db.View(func(txn *badger.Txn) error {
		v, err := txn.Get([]byte("serverURL"))
		if err != nil {
			return err
		}
		return v.Value(func(val []byte) error {
			s.serverURL.Store(string(val))
			return nil
		})
	})
}

func (s *Storage) getUpdatesFromServer(ctx context.Context, updated uint64) {
	s.remoteStatus.Store("❌ >")

	ticker := time.NewTicker(6 * time.Second)

	go func() {
		for range ticker.C {
			log.Debug().Msg("сервер недоступен")
			s.remoteStatus.Store("❌ >")

		}
	}()

	stream, err := s.remote.Updates(ctx, &proto.Request{Owner: s.owner, Updated: updated})

	for err == nil {
		select {
		case <-ctx.Done():
			return
		default:
			secret, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					log.Err(err)
					return
				}

				log.Err(err)
				s.remoteStatus.Store("❌ >")

				continue
			}

			if secret.Name == "__ping" {
				log.Debug().Msg("ping")
				ticker.Reset(6 * time.Second)

				if s.remoteStatus.Load().(string) != "✅ >" {
					s.remoteStatus.Store("✅ >")
				}
			}

			err = s.db.Update(func(txn *badger.Txn) error {
				return txn.Set([]byte(secret.Name), secret.GetData())
			})

			if err != nil {
				log.Err(err)
			}
		}
	}
}

func (s *Storage) connectToServer(ctx context.Context) {
	var (
		con grpc.ClientConnInterface
		err error
	)

	con, err = grpc.DialContext(ctx, strings.TrimPrefix(s.serverURL.Load().(string), "grpc://"), grpc.WithInsecure())
	if err != nil {
		log.Debug().Err(err)
		return
	}

	s.remote = proto.NewPublicClient(con)
}

func (s *Storage) Close() error {
	s.cancelFunc()
	return s.db.Close()
}

func (s *Storage) Status() func() (string, bool) {
	return func() (string, bool) {
		prefix, _ := s.remoteStatus.Load().(string)
		return prefix, true
	}
}

func (s *Storage) ListKeys(types ...model.SecretType) []string {
	result := make([]string, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{Reverse: false, AllVersions: false, Prefix: []byte("_secret:")})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			if len(types) == 0 {
				result = append(result, strings.TrimPrefix(string(it.Item().Key()), "_secret:"))
			}
			var value []byte
			err := it.Item().Value(func(val []byte) error {
				value = val
				return nil
			})
			if err != nil {
				fmt.Println("ошибка отображения секрета")
				return nil
			}
			buf := bytes.NewReader(value)
			secret := model.Secret{}
			err = gob.NewDecoder(buf).Decode(&secret)
			if err != nil {
				fmt.Println("ошибка отображения секрета")
				return nil
			}

			for _, v := range types {
				if v == secret.Type {
					result = append(result, strings.TrimPrefix(string(it.Item().Key()), "_secret:"))
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil
	}

	return result
}

func (s *Storage) Get(m model.Secret) *model.Secret {
	var result []byte

	err := s.db.View(func(txn *badger.Txn) error {
		var r *badger.Item
		r, err := txn.Get([]byte("_secret:" + m.Name))
		if err != nil {
			return err
		}

		return r.Value(func(val []byte) error {
			result = val
			return nil
		})
	})

	if err != nil {
		return nil
	}

	buf := bytes.NewReader(result)
	secret := model.Secret{}
	err = gob.NewDecoder(buf).Decode(&secret)

	if err != nil {
		fmt.Println("ошибка отображения секрета")
		return nil
	}

	return &secret
}

func (s *Storage) Update(m model.Secret) (err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err = encoder.Encode(m)

	if err != nil {
		return err
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(&badger.Entry{
			Key:   []byte("_secret:" + m.Name),
			Value: buf.Bytes(),
		})
	})

	if s.remote != nil {
		_, err = s.remote.Update(context.Background(), &proto.Secret{
			Data:    m.Data,
			Name:    m.Name,
			Owner:   s.owner,
			Type:    proto.SecretType(m.Type),
			Updated: uint64(time.Now().Unix()),
		})

		return err
	}

	return err
}

func (s *Storage) Remove(m model.Secret) (err error) {
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("_secret:" + m.Name))
	})

	if err != nil {
		return err
	}

	if s.serverURL.Load() != nil {
		_, err = s.remote.Update(context.Background(), &proto.Secret{
			Name:  m.Name,
			Owner: s.owner,
		})

		return err
	}

	return nil
}
