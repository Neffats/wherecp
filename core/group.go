package core

import (
	"errors"
	"sort"
)

const defaultGroupCapacity = 100

// Group is a structure that groups different object together.
// Acting as a container for different objects. Each of the item arrays are ordered
// for efficient searching.
// The Network objects are ordered by Address, and the Groups
// are ordered by name.
type Group struct {
	UID      int
	Name     string
	Hosts    []*Host
	Networks []*Network
	Ranges   []*Range
	Groups   []*Group
	Comment  string
}

// NewGroup returns a new empty group.
func NewGroup(name, comment string) *Group {
	return &Group{
		UID:      0,
		Name:     name,
		Hosts:    make([]*Host, 0),
		Networks: make([]*Network, 0),
		Ranges:   make([]*Range, 0),
		Groups:   make([]*Group, 0),
		Comment:  comment,
	}
}

// Add will add the specified object to the group.
// Supported types: Host/Network/Range/Group
func (g *Group) Add(obj interface{}) error {
	switch v := obj.(type) {
	case *Host:
		g.addHost(v)
	case *Network:
		g.addNetwork(v)
	case *Range:
		g.addRange(v)
	case *Group:
		g.Groups = append(g.Groups, v)
	default:
		return errors.New("unsupported data type")
	}
	return nil
}

func (g *Group) addHost(h *Host) {
	// Ordered smallest to largets by Address.
	i := sort.Search(len(g.Hosts), func(i int) bool {
		return *g.Hosts[i].Address > *h.Address
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newHosts := make([]*Host, len(g.Hosts)+1)
	// Shift the slice forward by one at the insert location.
	copy(newHosts[:i], g.Hosts[:i])
	copy(newHosts[i+1:], g.Hosts[i:])
	// Append host at the insert location.
	newHosts[i] = h
	g.Hosts = newHosts
}

func (g *Group) addNetwork(n *Network) {
	// Ordered smallest to largest by network address (first address) first
	// then by broadcast address (last address) second. Smallest networks will be in
	// front of larger networks i.e. 192.168.0.0/25 will be before 192.168.0.0/24
	i := sort.Search(len(g.Networks), func(i int) bool {
		thisStart, thisEnd := g.Networks[i].Value()
		otherStart, otherEnd := n.Value()

		addr := *thisStart >= *otherStart
		mask := *thisEnd >= *otherEnd
		return addr && mask
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newNets := make([]*Network, len(g.Networks)+1)
	// Shift the slice forward by one at the insert location.
	copy(newNets[:i], g.Networks[:i])
	copy(newNets[i+1:], g.Networks[i:])
	// Append network at the insert location.
	newNets[i] = n
	g.Networks = newNets
}

func (g *Group) addRange(r *Range) {
	// Ordered smallest to largest by start address (first address) first
	// then by end address (last address) second. Smaller ranges will come before
	// larger ranges i.e. 192.168.0.0-192.168.0.10 will be in front of 192.168.0.0-192.168.0.200
	i := sort.Search(len(g.Ranges), func(i int) bool {
		thisStart, thisEnd := g.Ranges[i].Value()
		otherStart, otherEnd := r.Value()

		start := *thisStart >= *otherStart
		end := *thisEnd >= *otherEnd
		return start && end
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newRange := make([]*Range, len(g.Ranges)+1)
	// Shift the slice forward by one at the insert location.
	copy(newRange[:i], g.Ranges[:i])
	copy(newRange[i+1:], g.Ranges[i:])
	// Append network at the insert location.
	newRange[i] = r
	g.Ranges = newRange
}

func (g *Group) addGroup(grp *Group) {
	// not implemented
}

// HasObject returns true if the group has a members object whose type and address matches the supplied object.
func (g *Group) HasObject(obj interface{}) (bool, error) {
	switch v := obj.(type) {
	case *Host:
		for _, hst := range g.Hosts {
			if hst.Match(v) {
				return true, nil
			}
		}
	case *Network:
		for _, net := range g.Networks {
			if net.Match(v) {
				return true, nil
			}
		}
	case *Range:
		for _, rng := range g.Ranges {
			if rng.Match(v) {
				return true, nil
			}
		}
	case *Group:
		for _, grp := range g.Groups {
			has, err := grp.HasObject(v)
			if err != nil {
				return false, err
			}
			return has, nil
		}
	default:
		return false, errors.New("unsupported data type")
	}
	return false, nil
}

func (g *Group) Contains(obj NetworkObject) bool {
	for _, h := range g.Hosts {
		if h.Match(obj) {
			return true
		}
	}
	for _, n := range g.Networks {
		if n.Contains(obj) {
			return true
		}
	}
	for _, r := range g.Ranges {
		if r.Contains(obj) {
			return true
		}
	}
	for _, grp := range g.Groups {
		if grp.Contains(obj) {
			return true
		}
	}
	return false
}
