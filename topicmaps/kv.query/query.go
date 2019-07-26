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
	"github.com/google/note-maps/topicmaps"
	"github.com/google/note-maps/topicmaps/kv.models"
)

// Store adds some query logic to models.Store.
type Store struct{ models.Store }

func (s *Store) TopicsByName(c *kv.IndexCursor, n int) ([]kv.Entity, error) {
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

type Flag int

const (
	Refs Flag = 1 << iota
	TopicMaps
	Topics
	Names
	Occurrences
)

func (s *Store) LoadTopic(t kv.Entity, f Flag) (*topicmaps.Topic, error) {
	var topic topicmaps.Topic

	if (f & Refs) != 0 {
		panic("loading refs is not yet implemented")
	}

	if (f & Names) != 0 {
		if ns, err := s.GetTopicNamesSlice([]kv.Entity{t}); err != nil {
			return nil, err
		} else if names, err := s.GetNameSlice(ns[0]); err != nil {
			return nil, err
		} else if len(names) > 0 {
			topic.Names = make([]*topicmaps.Name, 0, len(names))
			for _, stored := range names {
				var loaded topicmaps.Name
				loaded.Value = stored.Value
				topic.Names = append(topic.Names, &loaded)
			}
		}
	}

	if (f & Occurrences) != 0 {
		if os, err := s.GetTopicOccurrencesSlice([]kv.Entity{t}); err != nil {
			return nil, err
		} else if occurrences, err := s.GetOccurrenceSlice(os[0]); err != nil {
			return nil, err
		} else if len(occurrences) > 0 {
			topic.Occurrences = make([]*topicmaps.Occurrence, 0, len(occurrences))
			for _, stored := range occurrences {
				var loaded topicmaps.Occurrence
				loaded.Value = stored.Value
				topic.Occurrences = append(topic.Occurrences, &loaded)
			}
		}
	}

	return &topic, nil
}
