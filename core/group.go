package core

import (
	"errors"
	"sort"
)

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
		g.Networks = append(g.Networks, v)
	case *Range:
		g.Ranges = append(g.Ranges, v)
	case *Group:
		g.Groups = append(g.Groups, v)
	default:
		return errors.New("unsupported data type")
	}
	return nil
}

func (g *Group) addHost(h *Host) {
	if len(g.Hosts) == 0 {
		g.Hosts = append(g.Hosts, h)
		return
	}
	i := sort.Search(len(g.Hosts), func(i int) bool {
		return *g.Hosts[i].Address > *h.Address
	})
	copy(g.Hosts[i+1:], g.Hosts[i:])
	g.Hosts[i] = h
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
