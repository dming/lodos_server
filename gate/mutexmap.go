// Copyright 2014 mqant Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gate

import (
	"sync"
)

// BeeMap is a map with lock
type MutexMap struct {
	lock *sync.RWMutex
	bm   map[string]map[string]string
}

// NewBeeMap return new safemap
func NewMutexMap() *MutexMap {
	return &MutexMap{
		lock: new(sync.RWMutex),
		bm:   map[string]map[string]string{},
	}
}

// Get from maps return the k's value
func (m *MutexMap) Get(k string) map[string]string {
	m.lock.RLock()
	//defer m.lock.RUnlock()

	if val, ok := m.bm[k]; ok {
		m.lock.RUnlock()
		return val
	}
	m.lock.RUnlock()
	return nil
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *MutexMap) Set(k string, v map[string]string) bool {
	m.lock.Lock()
	//defer m.lock.Unlock()

	if _, ok := m.bm[k]; !ok {
		m.bm[k] = v
		m.lock.Unlock()
	} else if true {
		m.bm[k] = v
		m.lock.Unlock()
	} else {
		m.lock.Unlock()
		return false
	}
	return true
}

// Check Returns true if k is exist in the map.
func (m *MutexMap) Check(k string) bool {
	m.lock.RLock()
	//defer m.lock.RUnlock()

	if _, ok := m.bm[k]; ok {
		m.lock.RUnlock()
		return true
	}
	m.lock.RUnlock()
	return false
}

// Delete the given key and value.
func (m *MutexMap) Delete(k string) {
	m.lock.Lock()
	//defer m.lock.Unlock()

	delete(m.bm, k)
	m.lock.Unlock()
}

// Items returns all items in safemap.
func (m *MutexMap) Items() map[string]map[string]string {
	m.lock.RLock()
	//defer m.lock.RUnlock()

	r := make(map[string]map[string]string)
	for k, v := range m.bm {
		r[k] = v
	}
	m.lock.RUnlock()
	return r
}

// return true when a != b
func  (m *MutexMap) isDiff(a map[string]string, b map[string]string) bool {
	//m.lock.RLock()
	//defer m.lock.RUnlock()

	if len(a) != len(b) {
		return true
	}
	for k, v := range a {
		if _, ok := b[k]; !ok {
			return true
		} else if b[k] != v {
			return true
		}
	}
	return false
}