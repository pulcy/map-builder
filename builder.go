// Copyright (c) 2016 Epracom Advies.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mapbuilder

import (
	"fmt"
	"strings"
)

// MapBuilder is used to build nested `map[string]interface{}` structures.
// To build a structure, separate key parts by a `.`.
type MapBuilder struct {
	root   Map
	errors []error
}

type Map map[string]interface{}

// Set overwrites data into the map at the given key.
func (mb *MapBuilder) Set(key string, value interface{}) *MapBuilder {
	keys := strings.Split(key, ".")
	keyLen := len(keys)
	data := mb.getMap(keys[:keyLen-1], true)
	data[keys[keyLen-1]] = value
	return mb
}

// Remove overwrites data from the map at the given key.
func (mb *MapBuilder) Remove(key string) *MapBuilder {
	keys := strings.Split(key, ".")
	keyLen := len(keys)
	data := mb.getMap(keys[:keyLen-1], false)
	if data != nil {
		delete(data, keys[keyLen-1])
	}
	return mb
}

// Get returns the data stored in the map at the given key.
func (mb *MapBuilder) Get(key string) interface{} {
	if key == "" {
		if len(mb.root) == 0 {
			return nil
		}
		return mb.root
	}
	keys := strings.Split(key, ".")
	keyLen := len(keys)
	data := mb.getMap(keys[:keyLen-1], false)
	if data != nil {
		return data[keys[keyLen-1]]
	}
	return nil
}

// Err returns all the errors that were found while modifying the map.
func (mb *MapBuilder) Err() error {
	switch len(mb.errors) {
	case 0:
		return nil
	case 1:
		return mb.errors[0]
	default:
		msgs := []string{}
		for _, e := range mb.errors {
			msgs = append(msgs, e.Error())
		}
		return maskAny(fmt.Errorf(strings.Join(msgs, ", ")))
	}
}

// getMap returns the map that belongs to the given (parent) key.
// If key is empty, it returns the root map.
// If create is set, all missing maps will be created.
// If create is not set and a map is not found, nil will be returned.
func (mb *MapBuilder) getMap(key []string, create bool) Map {
	if mb.root == nil && create {
		mb.root = make(Map)
	}
	data := mb.root
	for i, k := range key {
		if data == nil {
			if create {
				panic(fmt.Sprintf("data==nil in '%s' %d", key, i))
			}
			return nil
		}
		nextDataRaw, ok := data[k]
		if !ok {
			if create {
				nextDataRaw = make(Map)
				data[k] = nextDataRaw
			} else {
				return nil
			}
		}
		nextData, ok := nextDataRaw.(Map)
		if !ok {
			mb.errors = append(mb.errors, maskAny(fmt.Errorf("'%s' is not a rawMap", strings.Join(key[:i], "."))))
			if create {
				nextData = make(Map)
				data[k] = nextData
			} else {
				return nil
			}
		}
		data = nextData
	}
	return data
}
