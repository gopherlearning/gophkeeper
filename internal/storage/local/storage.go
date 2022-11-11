package local

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/gopherlearning/gophkeeper/internal/model"
	"github.com/gopherlearning/gophkeeper/pkg/cryptor"
)

var (
	ErrNilCryptor = errors.New("empty cryptor")
)

type Storage struct {
	cancelFunc context.CancelFunc
	db         *badger.DB
	c          Cryptor
}

type Cryptor interface {
	Decrypt(chipher []byte) ([]byte, error)
	Encrypt(plaindata []byte, recipient ...string) ([]byte, error)
	GetKey() []byte
}

func NewLocalStorage(secret, path string, cryp Cryptor) (*Storage, error) {
	var key []byte

	_, err := os.Stat(path)
	if err != nil {
		if cryp == nil {
			return nil, ErrNilCryptor
		}

		key = cryp.GetKey()
	}

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

	if len(key) != 0 {
		err = db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(secret), key)
		})
	}

	if cryp == nil {
		err = db.View(func(txn *badger.Txn) error {
			var r *badger.Item
			r, err = txn.Get([]byte(secret))
			if err != nil {
				return err
			}
			return r.Value(func(val []byte) error {
				cryp, err = cryptor.NewCryptor(val)
				if err != nil {
					return err
				}

				return nil
			})
		})

		if err != nil {
			cancel()
			return nil, err
		}
	}

	return &Storage{db: db, cancelFunc: cancel, c: cryp}, nil
}

func (s *Storage) Close() error {
	s.cancelFunc()
	return s.db.Close()
}

func (s *Storage) ListKeys() []string {
	result := make([]string, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{Reverse: false, AllVersions: false, Prefix: []byte("secret")})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			result = append(result, fmt.Sprint(it.Item().Key()))
		}
		return nil
	})

	if err != nil {
		return nil
	}

	return result
}

func (s *Storage) Get(key []byte) []byte {
	var result []string
	err := s.db.View(func(txn *badger.Txn) error {
		var r *badger.Item
		r, err := txn.Get(key)
		if err != nil {
			return err
		}
		return r.Value(func(val []byte) error {
			cryp, err = cryptor.NewCryptor(val)
			if err != nil {
				return err
			}

			return nil
		})
	})

	if err != nil {
		return nil
	}

	return result
}

func (s *Storage) Update(m model.Secret) (err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err = encoder.Encode(m)

	if err != nil {
		return err
	}

	var enc []byte
	enc, err = s.c.Encrypt(buf.Bytes())

	if err != nil {
		return err
	}

	key := sha256.Sum256([]byte(m.Name))

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(&badger.Entry{
			Key:   key[:],
			Value: enc,
		})
	})

	return err
}
