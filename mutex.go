package database

import (
	"errors"
	"sync"

	"github.com/go-pg/pg/v9"
)

// Mutex is a shared mutex, stored in database
type Mutex struct {
	db     *pg.DB
	mu     sync.Mutex
	lockID int64
}

// NewMutex creates a new mutex
func NewMutex(db *pg.DB, id int64) (*Mutex, error) {
	if db.Options().PoolSize > 1 {
		return nil, errors.New("pool size for mutex cannot be greater than 1")
	}

	return &Mutex{db: db, lockID: id}, nil
}

// TryLock tries to acquire a lock and returns true in case of success, otherwise returns false
func (m *Mutex) TryLock() (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	res, err := m.db.Exec("SELECT 1 FROM pg_locks WHERE pid=pg_backend_pid() AND locktype='advisory' AND objid=?", m.lockID)
	if err != nil {
		return false, err
	}
	if res.RowsReturned() > 0 {
		return true, nil
	}

	var isLocked bool
	_, err = m.db.Query(&isLocked, "SELECT pg_try_advisory_lock(?)", m.lockID)
	if err != nil {
		return false, err
	}

	return isLocked, nil
}

// Unlock releases a lock in database
func (m *Mutex) Unlock() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.db.Exec("SELECT pg_advisory_unlock(?) ", m.lockID)
	return err
}
