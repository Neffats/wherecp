package core

import (
	"errors"
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
		g.Hosts = append(g.Hosts, v)
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
