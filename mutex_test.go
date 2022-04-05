// +build integration

package database

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var cfg = dbConnectCfg()

func dbConnectCfg() pg.Options {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	opts, err := pg.ParseURL(os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}

	return pg.Options{
		User:     opts.User,
		Password: opts.Password,
		Addr:     opts.Addr,
		PoolSize: 1,
		Database: opts.Database,
	}
}

func TestMutex_TryLock(t *testing.T) {
	t.Run("With single-connection pool", func(t *testing.T) {
		for _, res := range []bool{true, false} {
			conn := pg.Connect(&cfg)

			mu, err := NewMutex(conn, 111)
			assert.Nil(t, err)

			ok, err := mu.TryLock()
			assert.Nil(t, err)

			if ok != res {
				t.Errorf("TryLock is expected to return %v, but got %v", res, ok)
			}
		}
	})

	t.Run("With multi-connection pool", func(t *testing.T) {
		c := cfg
		c.PoolSize = 10

		conn := pg.Connect(&c)

		_, err := NewMutex(conn, 111)
		if err == nil {
			t.Error("NewMutex is expected to return error, but got nil")
		}
	})
}

func TestMutex_Unlock(t *testing.T) {
	for range make([]int, 2) {
		conn := pg.Connect(&cfg)

		mu, err := NewMutex(conn, 222)
		assert.Nil(t, err)

		ok, err := mu.TryLock()
		assert.Nil(t, err)

		if !ok {
			t.Error("TryLock is expected to return true, but got false")
		}

		err = mu.Unlock()
		assert.Nil(t, err)
	}
}
