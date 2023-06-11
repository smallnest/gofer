// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syncx

import (
	"sync"
	"sync/atomic"
)

// This file contains reference map implementations for unit-tests.

// mapInterface is the interface Map implements.
type mapInterface interface {
	Load(int) (int, bool)
	Store(key, value int)
	LoadOrStore(key, value int) (actual int, loaded bool)
	LoadAndDelete(key int) (value int, loaded bool)
	Delete(int)
	Swap(key, value int) (previous int, loaded bool)
	CompareAndSwap(key, old, new int) (swapped bool)
	CompareAndDelete(key, old int) (deleted bool)
	Range(func(key, value int) (shouldContinue bool))
}

var (
	_ mapInterface = &RWMutexMap{}
	_ mapInterface = &DeepCopyMap{}
)

// RWMutexMap is an implementation of mapInterface using a sync.RWMutex.
type RWMutexMap struct {
	mu    sync.RWMutex
	dirty map[int]int
}

func (m *RWMutexMap) Load(key int) (value int, ok bool) {
	m.mu.RLock()
	value, ok = m.dirty[key]
	m.mu.RUnlock()
	return
}

func (m *RWMutexMap) Store(key, value int) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[int]int)
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

func (m *RWMutexMap) LoadOrStore(key, value int) (actual int, loaded bool) {
	m.mu.Lock()
	actual, loaded = m.dirty[key]
	if !loaded {
		actual = value
		if m.dirty == nil {
			m.dirty = make(map[int]int)
		}
		m.dirty[key] = value
	}
	m.mu.Unlock()
	return actual, loaded
}

func (m *RWMutexMap) Swap(key, value int) (previous int, loaded bool) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[int]int)
	}

	previous, loaded = m.dirty[key]
	m.dirty[key] = value
	m.mu.Unlock()
	return
}

func (m *RWMutexMap) LoadAndDelete(key int) (value int, loaded bool) {
	m.mu.Lock()
	value, loaded = m.dirty[key]
	if !loaded {
		m.mu.Unlock()
		return 0, false
	}
	delete(m.dirty, key)
	m.mu.Unlock()
	return value, loaded
}

func (m *RWMutexMap) Delete(key int) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

func (m *RWMutexMap) CompareAndSwap(key, old, new int) (swapped bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dirty == nil {
		return false
	}

	value, loaded := m.dirty[key]
	if loaded && value == old {
		m.dirty[key] = new
		return true
	}
	return false
}

func (m *RWMutexMap) CompareAndDelete(key, old int) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dirty == nil {
		return false
	}

	value, loaded := m.dirty[key]
	if loaded && value == old {
		delete(m.dirty, key)
		return true
	}
	return false
}

func (m *RWMutexMap) Range(f func(key, value int) (shouldContinue bool)) {
	m.mu.RLock()
	keys := make([]int, 0, len(m.dirty))
	for k := range m.dirty {
		keys = append(keys, k)
	}
	m.mu.RUnlock()

	for _, k := range keys {
		v, ok := m.Load(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

// DeepCopyMap is an implementation of mapInterface using a Mutex and
// atomic.Value.  It makes deep copies of the map on every write to avoid
// acquiring the Mutex in Load.
type DeepCopyMap struct {
	mu    sync.Mutex
	clean atomic.Value
}

func (m *DeepCopyMap) Load(key int) (value int, ok bool) {
	clean, _ := m.clean.Load().(map[int]int)
	value, ok = clean[key]
	return value, ok
}

func (m *DeepCopyMap) Store(key, value int) {
	m.mu.Lock()
	dirty := m.dirty()
	dirty[key] = value
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) LoadOrStore(key, value int) (actual int, loaded bool) {
	clean, _ := m.clean.Load().(map[int]int)
	actual, loaded = clean[key]
	if loaded {
		return actual, loaded
	}

	m.mu.Lock()
	// Reload clean in case it changed while we were waiting on m.mu.
	clean, _ = m.clean.Load().(map[int]int)
	actual, loaded = clean[key]
	if !loaded {
		dirty := m.dirty()
		dirty[key] = value
		actual = value
		m.clean.Store(dirty)
	}
	m.mu.Unlock()
	return actual, loaded
}

func (m *DeepCopyMap) Swap(key, value int) (previous int, loaded bool) {
	m.mu.Lock()
	dirty := m.dirty()
	previous, loaded = dirty[key]
	dirty[key] = value
	m.clean.Store(dirty)
	m.mu.Unlock()
	return
}

func (m *DeepCopyMap) LoadAndDelete(key int) (value int, loaded bool) {
	m.mu.Lock()
	dirty := m.dirty()
	value, loaded = dirty[key]
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
	return
}

func (m *DeepCopyMap) Delete(key int) {
	m.mu.Lock()
	dirty := m.dirty()
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) CompareAndSwap(key, old, new int) (swapped bool) {
	clean, _ := m.clean.Load().(map[int]int)
	if previous, ok := clean[key]; !ok || previous != old {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	dirty := m.dirty()
	value, loaded := dirty[key]
	if loaded && value == old {
		dirty[key] = new
		m.clean.Store(dirty)
		return true
	}
	return false
}

func (m *DeepCopyMap) CompareAndDelete(key, old int) (deleted bool) {
	clean, _ := m.clean.Load().(map[int]int)
	if previous, ok := clean[key]; !ok || previous != old {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	dirty := m.dirty()
	value, loaded := dirty[key]
	if loaded && value == old {
		delete(dirty, key)
		m.clean.Store(dirty)
		return true
	}
	return false
}

func (m *DeepCopyMap) Range(f func(key, value int) (shouldContinue bool)) {
	clean, _ := m.clean.Load().(map[int]int)
	for k, v := range clean {
		if !f(k, v) {
			break
		}
	}
}

func (m *DeepCopyMap) dirty() map[int]int {
	clean, _ := m.clean.Load().(map[int]int)
	dirty := make(map[int]int, len(clean)+1)
	for k, v := range clean {
		dirty[k] = v
	}
	return dirty
}

// sync.Map

type SyncMap[K comparable, V any] struct {
	dirty sync.Map
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.dirty.Load(key)
	if !ok {
		var v V
		return v, false
	}

	return v.(V), true
}

func (m *SyncMap[K, V]) Store(key, value int) {
	m.dirty.Store(key, value)
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.dirty.LoadOrStore(key, value)

	return v.(V), loaded
}

func (m *SyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, loaded := m.dirty.Swap(key, value)

	if !loaded {
		var v V
		return v, false
	}

	return v.(V), true
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.dirty.LoadAndDelete(key)
	if !loaded {
		var v V
		return v, false
	}

	return v.(V), true
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.dirty.Delete(key)
}

func (m *SyncMap[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	return m.dirty.CompareAndSwap(key, old, new)
}

func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.dirty.CompareAndDelete(key, old)
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) (shouldContinue bool)) {
	m.dirty.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}
