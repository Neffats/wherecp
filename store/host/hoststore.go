package hoststore

import (
	"errors"
	"fmt"
	"sync"
	
	"github.com/Neffats/wherecp/core"
)

var (
	ErrHostNotFound = errors.New("host not found")
)

type HostPuller interface {
	PullHosts() ([]*core.Host, error)
}

type HostStore struct {
	Hosts []*core.Host
	Puller HostPuller

	mux sync.RWMutex
}

func New(puller HostPuller) *HostStore {
	return &HostStore{
		Hosts: make([]*core.Host, 0),
		Puller: puller,
	}
}

func (hs *HostStore) Init() error {
	hosts, err := hs.Puller.PullHosts()
	if err != nil {
		return fmt.Errorf("failed to pull hosts from source: %v", err)
	}
	hs.Hosts = hosts
	return nil
}

func (hs *HostStore) All() []*core.Host {
	hs.mux.RLock()
	defer hs.mux.RUnlock()
	h := make([]*core.Host, len(hs.Hosts))
	copy(h, hs.Hosts)
	return h
}

func (hs *HostStore) Create(host *core.Host) error {
	existing, err := hs.Get(host.UID)
	if !errors.Is(err, ErrHostNotFound) {
		return fmt.Errorf("failed to determine whether host is already present: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("host already in store")
	}
	hs.mux.Lock()
	defer hs.mux.Unlock()
	hs.Hosts = append(hs.Hosts, host)
	return nil
}

func (hs *HostStore) Get(uid string) (*core.Host, error) {
	hs.mux.RLock()
	defer hs.mux.RUnlock()

	for _, h := range hs.Hosts {
		if uid == h.UID {
			return h, nil
		}
	}

	return nil, ErrHostNotFound
}

func (hs *HostStore) Update(uid string, updated *core.Host) error {
	hs.mux.RLock()
	defer hs.mux.RUnlock()
	for i, h := range hs.Hosts {
		if h.UID == uid {
			hs.mux.Lock()
			hs.Hosts[i] = updated
			hs.mux.Unlock()
			return nil
		}
	}
	return ErrHostNotFound
}

func (hs *HostStore) Delete(uid string) error {
	for i, h := range hs.Hosts {
		if h.UID == uid {
			hs.mux.Lock()
			newHosts := make([]*core.Host, len(hs.Hosts)-1)
			copy(newHosts[:i], hs.Hosts[:i])
			copy(newHosts[i:], hs.Hosts[i+1:])
			hs.Hosts = newHosts
			hs.mux.Unlock()
			return nil
		}
	}
	return ErrHostNotFound
}

func (hs *HostStore) WithIP(ip string) ([]*core.Host, error) {
	matched := make([]*core.Host, 0)
	matcher, err := core.NewHost("", ip, "")
	if err != nil {
		return matched, fmt.Errorf("failed to create host to compare: %v", err) 
	}

	hs.mux.RLock()
	defer hs.mux.RUnlock()
	for _, h := range hs.Hosts {
		if matcher.Match(h) {
			matched = append(matched, h)
		}
	}
	return matched, nil
}
