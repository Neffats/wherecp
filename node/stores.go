package node

import (
	"github.com/Neffats/wherecp/core"
)

type RuleStorer interface {
	// Return every rule.
	All() []core.Rule
	Insert(rule core.Rule) error
	// Return a rule object from it's uid.
	Get(uid string) (core.Rule, error)
	// Update a rule object.
	Update(uid string, core.Rule) error
	// Delete a rule object from the store.
	Delete(uid string) error
}

type HostStorer interface {
	// Return every rule.
	All() []*core.Host
	Insert(hst *core.Host) error
	// Return a host object from it's uid.
	Get(uid string) (*core.Host, error)
	// Update a host object.
	Update(uid string, updated core.Host) error
	// Delete a host object from the store.
	Delete(uid string) error
	// Returns a list of hosts that have the given IP.
	WithIP(ip string) ([]*core.Host, error)
}

type NetworkStorer interface {
	// Return every rule.
	All() []*core.Network
	Insert(net *core.Network) error
	// Return a network object from it's uid.
	Get(uid string) (*core.Network, error)
	// Update a network object.
	Update(uid string, core.Network) error
	// Delete a network object from the store.
	Delete(uid string) error
	// Returns a list of networks that have the given IP.
	// IP in format of address/mask (i.e. 192.168.0.0/24)
	WithIP(ip string) ([]*core.Network, error)
}

type RangeStorer interface {
	// Return every rule.
	All() []*core.Range
	Create(rng *core.Range) error
	// Return a range object from it's uid.
	Get(uid string) (*core.Range, error)
	// Update a range object.
	Update(uid string, core.Range) error
	// Delete a range object from the store.
	Delete(uid string) error
	// Returns a list of ranges that have the given IP.
	// IP in format of start-end (i.e. 192.168.0.0-192.168.0.5)
	WithIP(ip string) ([]*core.Range, error)
}

type GroupStorer interface {
	// Return every rule.
	All() []*core.Group
	Insert(grp *core.Group) error
	// Return a group object from it's uid.
	Get(uid string) (*core.Group, error)
	// Update a group object.
	Update(uid string, core.Group) error
	// Delete a group object from the store.
	Delete(uid string) error
	// Returns a list of groups that have the given name.
	WithName(name string) ([]*core.Group, error)
}
