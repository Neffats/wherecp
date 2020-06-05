package core

import (
	"github.com/Neffats/ip"
)

// NetworkObject is the most basic representation of any supported object type.
// Every object can be converted to this type, this is how different types can be
// compared.
type NetworkObject struct {
	Start ip.Address
	End ip.Address
}

type NetworkUnpacker interface {
	Unpack() ([]NetworkObject)
}

type Containser interface {
	Contains(NetworkUnpacker) bool
}
