package core

import (
	"errors"
	"sort"
)

const defaultGroupCapacity = 100

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
	i := sort.Search(len(g.Hosts), func(i int) bool {
		return *g.Hosts[i].Address > *h.Address
	})

	// TODO: Is there a nicer way of doing this?
	newHosts := make([]*Host, len(g.Hosts)+1)
	copy(newHosts[:i], g.Hosts[:i])
	copy(newHosts[i+1:], g.Hosts[i:])
	newHosts[i] = h
	g.Hosts = newHosts
}

func (g *Group) addNetwork(n *Network) {
	i := sort.Search(len(g.Networks), func(i int) bool {
		thisStart, thisEnd := g.Networks[i].Value()
		otherStart, otherEnd := n.Value()

		addr := *thisStart >= *otherStart
		mask := *thisEnd >= *otherEnd
		return addr && mask
	})

	// TODO: Is there a nicer way of doing this?
	newNets := make([]*Network, len(g.Networks)+1)
	copy(newNets[:i], g.Networks[:i])
	copy(newNets[i+1:], g.Networks[i:])
	newNets[i] = n
	g.Networks = newNets
}

func (g *Group) addRange(r *Range) {
	i := sort.Search(len(g.Ranges), func(i int) bool {
		thisStart, thisEnd := g.Ranges[i].Value()
		otherStart, otherEnd := r.Value()

		start := *thisStart >= *otherStart
		end := *thisEnd >= *otherEnd
		return start && end
	})

	// TODO: Is there a nicer way of doing this?
	newRange := make([]*Range, len(g.Ranges)+1)
	copy(newRange[:i], g.Ranges[:i])
	copy(newRange[i+1:], g.Ranges[i:])
	newRange[i] = r
	g.Ranges = newRange
}

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
