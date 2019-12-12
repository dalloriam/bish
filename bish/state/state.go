package state

import (
	"sync"

	"github.com/dalloriam/bish/bish/hooks"
)

// State stores the entire state of the shell.
type State struct {
	data  map[string]interface{}
	hooks map[string]hooks.Hook
	mtx   sync.RWMutex
}

// New returns a new shell context store.
func New() *State {
	return &State{data: make(map[string]interface{}), hooks: make(map[string]hooks.Hook)}
}

// AddHook registers a new hook in shell state.
func (c *State) AddHook(name string, hk hooks.Hook) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.hooks[name] = hk
}

// Hooks returns all currently-defined hooks.
func (c *State) Hooks() []hooks.Hook {
	var hks []hooks.Hook
	for _, v := range c.hooks {
		hks = append(hks, v)
	}
	return hks
}

// GetKey reads & returns a key from the context store.
func (c *State) GetKey(domain string, key string) (interface{}, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	fmtedKey := domain + "_" + key
	v, ok := c.data[fmtedKey]
	return v, ok
}

// SetKey sets a key in the context store.
func (c *State) SetKey(domain string, key string, value interface{}) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	fmtedKey := domain + "_" + key
	c.data[fmtedKey] = value
}
