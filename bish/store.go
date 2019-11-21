package bish

import (
	"sync"
)

// ContextStore stores the current context of the shell.
type ContextStore struct {
	data map[string]interface{}
	mtx  sync.RWMutex
}

// NewContext returns a new shell context store.
func NewContext() *ContextStore {
	return &ContextStore{data: make(map[string]interface{})}
}

// GetKey reads & returns a key from the context store.
func (c *ContextStore) GetKey(domain string, key string) (interface{}, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	fmtedKey := domain + "_" + key
	v, ok := c.data[fmtedKey]
	return v, ok
}

func (c *ContextStore) SetKey(domain string, key string, value interface{}) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	fmtedKey := domain + "_" + key
	c.data[fmtedKey] = value
}
