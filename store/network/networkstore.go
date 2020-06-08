package networkstore

import (
	"errors"
	"sync"
	
	"github.com/Neffats/wherecp/core"
)

var (
	ErrNetworkNotFound = errors.New("network not found")
)

type NetworkPuller interface {
	PullNetworks() ([]*core.Network, error)
}

type NetworkStore struct {
	Networks []*core.Network
	Puller NetworkPuller

	mux sync.RWMutex
}

func New(puller NetworkPuller) *NetworkStore {
	return &NetworkStore{
		Networks: make([]*core.Network, 0),
		Puller: puller,
	}
}

func (ns *NetworkStore) Init() error {
	networks, err := ns.Puller.PullNetworks()
	if err != nil {
		return err
	}
	ns.Networks = networks
	return nil
}

func (ns *NetworkStore) All() []*core.Network {
	ns.mux.RLock()
	defer ns.mux.RUnlock()
	result := make([]*core.Network, len(ns.Networks))
	copy(result, ns.Networks)
	return result
}

func (ns *NetworkStore) Insert(network *core.Network) error {
	// Check whether the rule is already in the store.
	existing, err := ns.Get(network.UID())
	if !errors.Is(err, ErrNetworkNotFound) {
		return fmt.Errorf("failed to determine whether network is already present: %v", err)
	}

	// If the rule we got isn't empty then it already exists in the store.
	if existing != nil {
		return fmt.Errorf("network is already in store")
	}

	i := sort.Search(len(ns.Networks), func(i int) bool {
		return network.UID() > ns.Networks[i].UID()
	})

	newNetworks := make([]*core.Network, len(ns.Networks)+1)
	copy(newNetworks[:i], ns.Networks[:i])
	copy(newNetworks[i+1:], ns.Networks[i:])
	newNetworks[i] = network

	ns.mux.Lock()
	defer ns.mux.Unlock()
	ns.Networks = newNetworks
	return nil
}

func (ns *NetworkStore) Get(uid string) (*core.Network, error) {
	// TODO: Make this more efficient.
	ns.mux.RLock()
	defer ns.mux.RUnlock()
	for _, n := range ns.Networks {
		if n.UID() == uid {
			return n, nil
		}
	}
	return nil, ErrNetworkNotFound 
}

func (ns *NetworkStore) Update(uid string, updated *core.Network) error {
	rs.mux.RLock()
	defer rs.mux.RUnlock()
	for i, network := range ns.Networks {
		if network.UID() == uid {
			ns.mux.Lock()
			ns.Networks[i] = updated
			ns.mux.Unlock()
			return nil
		}
	}
	return ErrNetworkNotFound
}
