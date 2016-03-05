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
	"reflect"
	"testing"
)

func assertGetToBe(t *testing.T, mb MapBuilder, key string, value interface{}) {
	found := mb.Get(key)
	if !reflect.DeepEqual(found, value) {
		t.Errorf("diff at '%s': expected '%v' got '%v'", key, value, found)
	}
}

func TestSet(t *testing.T) {
	mb := MapBuilder{}
	mb.Set("key", 1)
	assertGetToBe(t, mb, "key", 1)

	mb.Set("parent.key", 5)
	assertGetToBe(t, mb, "parent.key", 5)

	mb.Set("parent.foo.key", 15)
	assertGetToBe(t, mb, "parent.foo.key", 15)

	expected := Map{
		"key": 1,
		"parent": Map{
			"key": 5,
			"foo": Map{
				"key": 15,
			},
		},
	}
	assertGetToBe(t, mb, "", expected)

	if err := mb.Err(); err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
}

func TestRemove(t *testing.T) {
	mb := MapBuilder{}
	mb.Set("key", "xyz")
	assertGetToBe(t, mb, "key", "xyz")
	mb.Remove("key")
	assertGetToBe(t, mb, "key", nil)

	mb.Set("parent.key", "xyz")
	assertGetToBe(t, mb, "parent.key", "xyz")
	mb.Remove("parent.key")
	assertGetToBe(t, mb, "parent.key", nil)
	assertGetToBe(t, mb, "parent", Map{})

	mb.Set("parent.key", "xyz")
	assertGetToBe(t, mb, "parent.key", "xyz")
	mb.Remove("parent")
	assertGetToBe(t, mb, "parent.key", nil)
	assertGetToBe(t, mb, "parent", nil)
	assertGetToBe(t, mb, "", nil)

	if err := mb.Err(); err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
}

func TestEmpty(t *testing.T) {
	mb := MapBuilder{}
	assertGetToBe(t, mb, "", nil)
}

func TestError(t *testing.T) {
	mb := MapBuilder{}
	mb.Set("key", "xyz")
	assertGetToBe(t, mb, "key", "xyz")

	mb.Set("key.child", "xyz")
	assertGetToBe(t, mb, "key.child", "xyz")

	if err := mb.Err(); err == nil {
		t.Errorf("expected error, got nil")
	}
}
