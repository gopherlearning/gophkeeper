package local

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/gopherlearning/gophkeeper/internal/model"
)

var (
	ErrNilCryptor = errors.New("empty cryptor")
)

type Storage struct {
	cancelFunc context.CancelFunc
	db         *badger.DB
}

func NewLocalStorage(secret, path string) (*Storage, error) {
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

	return &Storage{db: db, cancelFunc: cancel}, nil
}

func (s *Storage) Close() error {
	s.cancelFunc()
	return s.db.Close()
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

	return err
}

func (s *Storage) Remove(m model.Secret) (err error) {
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("_secret:" + m.Name))
	})

	if err != nil {
		return err
	}

	return nil
}
