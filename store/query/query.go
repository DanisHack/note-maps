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

package query

import (
	"github.com/google/note-maps/kv"
	"github.com/google/note-maps/store/models"
	"github.com/google/note-maps/topicmaps"
)

// Txn adds some query logic to models.Txn.
type Txn struct{ models.Txn }

// TopicsByName returns a page of topics ordered by name.
func (s *Txn) TopicsByName(c *kv.IndexCursor, n int) ([]kv.Entity, error) {
	ns, err := s.EntitiesByNameValue(c, n)
	if err != nil {
		return nil, err
	}
	names, err := s.GetNameSlice(ns)
	if err != nil {
		return nil, err
	}
	ts := make([]kv.Entity, len(names))
	for i := range names {
		ts[i] = kv.Entity(names[i].Topic)
	}
	return ts, nil
}

// Mask describes which fields should be included in a response.
type Mask int

const (
	// Refs indicates IIs, SIs, and SLs should be included in a response.
	Refs Mask = 1 << iota
	// TopicMaps indicates TopicMaps should be included in a response.
	TopicMaps
	// Topics indicates Topics should be included in a response.
	Topics
	// Names indicates Names should be included in a response.
	Names
	// Occurrences indicates Occurrences should be included in a response.
	Occurrences
)

// LoadTopic retrieves f fields of te into a topicmaps.Topic.
func (s *Txn) LoadTopic(te kv.Entity, f Mask) (*topicmaps.Topic, error) {
	var topic topicmaps.Topic

	if (f & Refs) != 0 {
		panic("loading refs is not yet implemented")
	}

	if (f & Names) != 0 {
		if nes, err := s.GetTopicNames(te); err != nil {
			return nil, err
		} else if ns, err := s.GetNameSlice(nes); err != nil {
			return nil, err
		} else if len(ns) > 0 {
			topic.Names = make([]*topicmaps.Name, 0, len(ns))
			for _, n := range ns {
				var loaded topicmaps.Name
				loaded.Value = n.Value
				topic.Names = append(topic.Names, &loaded)
			}
		}
	}

	if (f & Occurrences) != 0 {
		if oes, err := s.GetTopicOccurrences(te); err != nil {
			return nil, err
		} else if os, err := s.GetOccurrenceSlice(oes); err != nil {
			return nil, err
		} else if len(os) > 0 {
			topic.Occurrences = make([]*topicmaps.Occurrence, 0, len(os))
			for _, o := range os {
				var loaded topicmaps.Occurrence
				loaded.Value = o.Value
				topic.Occurrences = append(topic.Occurrences, &loaded)
			}
		}
	}

	return &topic, nil
}
