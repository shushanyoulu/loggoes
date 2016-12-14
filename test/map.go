// Package ttlmap provides a map-like interface with string keys and expirable
// items. Keys are currently limited to strings.
package main

import "errors"

// Errors returned by Set and SetNX operations.
var (
	ErrExists  = errors.New("item already exists")
	ErrDrained = errors.New("map was drained")
)

// Options for initializing a Map.
type Options struct {
	InitialCapacity int
	OnWillExpire    func(key string, item *Item)
	OnWillEvict     func(key string, item *Item)
}

// Map is the equivalent of a map[string]interface{} but with expirable Items.
type Map struct {
	store  *store
	keeper *keeper
}

// New creates a new Map with given options.
func New(opts *Options) *Map {
	// fmt.Printf("%T,%v\n", opts, &opts)
	if opts == nil {
		opts = &Options{}
	}
	store := newStore(opts)
	m := &Map{
		store:  store,
		keeper: newKeeper(store),
	}
	go m.keeper.run()
	return m
}

// Len returns the number of elements in the map.
func (m *Map) Len() int {
	m.store.RLock()
	n := len(m.store.kv)
	m.store.RUnlock()
	return n
}

// Get returns the item in the map given its key.
func (m *Map) Get(key string) *Item {
	m.store.RLock()
	if m.keeper.drained {
		m.store.RUnlock()
		return nil
	}
	pqi := m.store.kv[key]
	m.store.RUnlock()
	if pqi != nil {
		return pqi.item
	}
	return nil
}

// Set assigns an expirable Item with the specified key in the map.
// ErrDrained will be returned if the map is already drained.
func (m *Map) Set(key string, item *Item) error {
	m.store.Lock()
	if m.keeper.drained {
		m.store.Unlock()
		return ErrDrained
	}
	if pqi := m.store.kv[key]; pqi != nil {
		m.expireOrEvict(pqi)
	}
	m.set(key, item)
	m.store.Unlock()
	return nil
}

// SetNX assigns an expirable Item with the specified key in the map, only if
// the key is not already being in use.
// ErrExists will be returned if the key already exists.
// ErrDrained will be returned if the map is already drained.
func (m *Map) SetNX(key string, item *Item) error {
	m.store.Lock()
	if m.keeper.drained {
		m.store.Unlock()
		return ErrDrained
	}
	if pqi := m.store.kv[key]; pqi != nil {
		m.store.Unlock()
		return ErrExists
	}
	m.set(key, item)
	m.store.Unlock()
	return nil
}

// Delete deletes the item with the specified key in the map.
// If an item is found, it is returned.
func (m *Map) Delete(key string) *Item {
	m.store.Lock()
	if m.keeper.drained {
		m.store.Unlock()
		return nil
	}
	if pqi := m.store.kv[key]; pqi != nil {
		m.delete(pqi)
		m.store.Unlock()
		return pqi.item
	}
	m.store.Unlock()
	return nil
}

// Draining returns the channel that is closed when the map starts draining.
func (m *Map) Draining() <-chan struct{} {
	return m.keeper.drainingChan
}

// Drain evicts all remaining elements from the map and terminates the usage of
// this map.
func (m *Map) Drain() {
	m.keeper.signalDrain()
	<-m.keeper.doneChan

}

func (m *Map) expireOrEvict(pqi *pqitem) {
	if pqi.index == 0 {
		m.keeper.signalUpdate()
	}
	if !m.store.tryExpire(pqi) {
		m.store.evict(pqi)
	}
}

func (m *Map) set(key string, item *Item) {
	pqi := &pqitem{
		key:   key,
		item:  item,
		index: -1,
	}
	m.store.set(pqi)
	if pqi.index == 0 {
		m.keeper.signalUpdate()
	}
}

func (m *Map) delete(pqi *pqitem) {
	if pqi.index == 0 {
		m.keeper.signalUpdate()
	}
	m.store.delete(pqi)
}
