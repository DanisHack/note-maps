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

package p2p

import (
	"fmt"
	"sync"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	//"github.com/ipfs/go-cid"
	//"github.com/mr-tron/base58"
	//"github.com/multiformats/go-multiaddr"
	//core "github.com/textileio/go-threads/core/store"
	//"github.com/textileio/go-threads/core/service"
	"github.com/textileio/go-threads/store"

	tmpb "github.com/google/note-maps/tmaps/pb"
)

type config struct {
	baseDir string
}

type PeerOption interface {
	apply(*config)
}

type peerOptionFunc func(*config)

func (of peerOptionFunc) apply(c *config) { of(c) }

func WithBaseDir(baseDir string) PeerOption {
	return peerOptionFunc(func(c *config) { c.baseDir = baseDir })
}

type Peer struct {
	ipfs    *ipfslite.Peer
	store   *store.Store
	service store.ServiceBoostrapper
	wg      sync.WaitGroup
	mx      sync.Mutex
}

func NewPeer(repoPath string, opts ...PeerOption) (*Peer, error) {
	c := config{baseDir: repoPath}
	for _, o := range opts {
		o.apply(&c)
	}

	var p *Peer

	service, err := store.DefaultService(c.baseDir, store.WithServiceDebug(true))
	if err != nil {
		return nil, fmt.Errorf("error while bootstrapping: %s", err)
	}
	defer func() {
		if p == nil {
			service.Close()
		}
	}()

	s, err := store.NewStore(service, store.WithRepoPath(c.baseDir))
	if err != nil {
		return nil, fmt.Errorf("error when creating event store: %s", err)
	}
	defer func() {
		if p == nil {
			s.Close()
		}
	}()

	if _, err = s.Register("notemapName", &name{}); err != nil {
		return nil, err
	}
	if _, err = s.Register("notemapOccurrence", &occurrence{}); err != nil {
		return nil, err
	}
	if _, err = s.Register("notemapTopic", &topic{}); err != nil {
		return nil, err
	}
	if _, err = s.Register("notemapTopicMap", &topicMap{}); err != nil {
		return nil, err
	}

	p = &Peer{
		service: service,
		store:   s,
	}
	return p, nil
}

func (p *Peer) Close() error {
	p.mx.Lock()
	defer p.mx.Unlock()

	return p.service.Close()
}

type name struct {
	tmpb.Name
	ID string
}
type occurrence struct {
	tmpb.Occurrence
	ID string
}
type topic struct {
	tmpb.Topic
	ID string
}
type topicMap struct {
	tmpb.TopicMap
	ID string
}
