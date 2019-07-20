// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package memory provides an in-memory implementation of kv.Store.
package memory

import (
	"sort"
	"strings"
	"sync"

	"github.com/google/note-maps/kv"
)

// New returns a memory-backed implementation of the kv.Store interface
// intended exclusively for use in tests, and not in production.
func New() kv.Store {
	return &store{
		m: make(map[string][]byte),
	}
}

type store struct {
	m     map[string][]byte
	next  kv.Entity
	mutex sync.Mutex
}

func (s *store) Alloc() (kv.Entity, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.next++
	return s.next, nil
}

func (s *store) Get(k []byte, f func([]byte) error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return f(s.m[string(k)])
}

func (s *store) Set(k, v []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[string(k)] = v
	return nil
}

func (s *store) PrefixIterator(prefix []byte) kv.Iterator {
	var iter iterator
	p := string(prefix)
	for k, v := range s.m {
		if strings.HasPrefix(k, p) {
			iter.pairs = append(iter.pairs, pair{
				key:   k[len(p):],
				value: v,
			})
		}
	}
	sort.Slice(
		iter.pairs,
		func(a, b int) bool { return iter.pairs[a].key < iter.pairs[b].key })
	return &iter
}

type pair struct {
	key   string
	value []byte
}

type iterator struct {
	pairs []pair
	i     int
}

func (i *iterator) Seek(key []byte) {
	k := string(key)
	i.i = 0
	for i.i = 0; i.i < len(i.pairs) && i.pairs[i.i].key < k; i.i++ {
	}
}

func (i *iterator) Next() { i.i++ }

func (i *iterator) Valid() bool {
	return 0 <= i.i && i.i < len(i.pairs)
}

func (i *iterator) Key() []byte { return []byte(i.pairs[i.i].key) }

func (i *iterator) Value(f func([]byte) error) error {
	return f(i.pairs[i.i].value)
}

func (i *iterator) Discard() {}
