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

package docs

import (
	"fmt"
	"testing"

	"github.com/google/note-maps/kv"
	"github.com/google/note-maps/kv/memory"
)

func TestSchemaSetScanLookup(t *testing.T) {
	store := Store{Store: memory.New()}
	e, err := store.Alloc()
	if err != nil {
		t.Error(e)
	}
	sample := Document{
		Title:   "Test Title",
		Content: "Ipsum dolor etcetera",
	}
	err = store.SetDocument(e, &sample)
	if err != nil {
		t.Error(err)
	}
	ds, err := store.GetDocumentSlice([]kv.Entity{e, e})
	if err != nil {
		t.Error(err)
	} else if len(ds) != 2 || ds[0] != sample || ds[1] != sample {
		t.Error("want", []Document{sample, sample}, "got", ds)
	}
	matches, err := store.EntitiesMatchingDocumentTitle("test title")
	if err != nil {
		t.Error(err)
	} else if len(matches) != 1 {
		t.Error("want 1 match, got", len(matches), ":", matches)
	} else {
		ds, err = store.GetDocumentSlice(matches)
		if err != nil {
			t.Error(err)
		} else if len(ds) != 1 {
			t.Error("want one documents, got", ds)
		} else if ds[0].Title != "Test Title" {
			t.Errorf("want %v, got %v",
				"Test Title", ds[0].Title)
		}
	}
}

func TestSchemaByTitle(t *testing.T) {
	store := Store{Store: memory.New()}
	for i := 0; i < 10; i++ {
		for _, name := range []string{"Foo", "Bar", "Quux"} {
			e, err := store.Alloc()
			if err != nil {
				t.Fatal(e)
			}
			err = store.SetDocument(e, &Document{
				Title:   fmt.Sprintf("%s #%v", name, i),
				Content: "Ipsum dolor etcetera",
			})
			if err != nil {
				t.Error(err)
			}
		}
	}
	var cursor kv.IndexCursor
	var docs []Document
	already := make(map[kv.Entity]bool)
	for i := 0; ; i++ {
		es, err := store.EntitiesByDocumentTitle(&cursor, 5)
		println(string(cursor.Key), cursor.Offset)
		if err != nil {
			t.Error(err)
			break
		}
		if len(es) == 0 {
			break
		}
		for _, e := range es {
			if already[e] {
				t.Error("duplicate", e)
			}
			already[e] = true
		}
		ds, err := store.GetDocumentSlice(es)
		if err != nil {
			t.Error(err)
			break
		}
		docs = append(docs, ds...)
	}
	for i := 1; i < len(docs); i++ {
		if docs[i-1].Title > docs[i].Title {
			t.Errorf("want %#v before %#v, got after",
				docs[i-1].Title, docs[i].Title)
		}
	}
}
