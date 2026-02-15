package db

import (
	"database/sql"
	"sync"
)

type StatementCache struct {
	db    *sql.DB
	cache map[string]*sql.Stmt
	mu    sync.RWMutex
}

func NewStatementCache(db *sql.DB) *StatementCache {
	return &StatementCache{
		db:    db,
		cache: make(map[string]*sql.Stmt),
	}
}

func (sc *StatementCache) Prepare(query string) (*sql.Stmt, error) {
	sc.mu.RLock()
	stmt, ok := sc.cache[query]
	sc.mu.RUnlock()

	if ok {
		return stmt, nil
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Double check
	if stmt, ok := sc.cache[query]; ok {
		return stmt, nil
	}

	stmt, err := sc.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	sc.cache[query] = stmt
	return stmt, nil
}

func (sc *StatementCache) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for _, stmt := range sc.cache {
		stmt.Close()
	}
	sc.cache = make(map[string]*sql.Stmt)
}
