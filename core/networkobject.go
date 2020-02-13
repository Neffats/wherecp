package core

import (
	"github.com/Neffats/ip"
)

type NetworkObject interface {
	Value() (start *ip.Address, end *ip.Address)
}
