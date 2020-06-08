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
