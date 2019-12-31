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

// Contains will return true if the group contains the specified object.
// The strict parameter specified whether we want to look for strict matches.
// A strict match is where the objects are the same type and content. We use the Match() function.
// A non-strict match is where an object can be contained by another.
// i.e. will return true for host 192.168.1.1 if network 192.168.1.0/24 is in the group.
// There's probably a better way to do this, but will keep for now.
func (g *Group) Contains(obj interface{}, strict bool) (bool, error) {
	var (
		result bool
		err    error
	)

	if strict {
		result, err = g.containsStrict(obj)
	} else {
		result, err = g.containsNotStrict(obj)
	}

	return result, err
}

func (g *Group) containsStrict(obj interface{}) (bool, error) {
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
			contains, err := grp.containsStrict(v)
			if err != nil {
				return false, err
			}
			return contains, nil
		}
	default:
		return false, errors.New("unsupported data type")
	}
	return false, nil
}

func (g *Group) containsNotStrict(obj interface{}) (bool, error) {
	switch v := obj.(type) {
	case *Host:
		for _, hst := range g.Hosts {
			if hst.Match(v) {
				return true, nil
			}
		}
		for _, net := range g.Networks {
			match, err := net.Contains(obj)
			if err != nil {
				return false, err
			}
			if match {
				return match, nil
			}
		}
		for _, rng := range g.Ranges {
			match, err := rng.Contains(obj)
			if err != nil {
				return false, err
			}
			if match {
				return true, nil
			}
		}
	case *Network, *Range:
		for _, net := range g.Networks {
			match, err := net.Contains(obj)
			if err != nil {
				return false, err
			}
			if match {
				return match, nil
			}
		}
		for _, rng := range g.Ranges {
			match, err := rng.Contains(obj)
			if err != nil {
				return false, err
			}
			if match {
				return true, nil
			}
		}
	default:
		return false, errors.New("unsupported data type")
	}

	return false, nil
}
