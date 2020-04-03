package core

import (
	"errors"
	"fmt"
	"reflect"
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

// Match will return true if the two groups are identical.
func (g *Group) Match(grp *Group) bool {
	return reflect.DeepEqual(g, grp)
}

// MatchContent will return true if both groups contain the same members.
func (g *Group) MatchContent(grp *Group) bool {
	// Check if the lengths of the groups match.
	// If they don't then the two groups must be different.
	if len(g.Hosts) != len(grp.Hosts) {
		return false
	}
	if len(g.Networks) != len(grp.Networks) {
		return false
	}
	if len(g.Ranges) != len(grp.Ranges) {
		return false
	}
	if len(g.Groups) != len(grp.Groups) {
		return false
	}

	var match bool

	// Compare Hosts of groups.
	// All group members are sorted, so all members should be in the same location.
	for i := 0; i < len(g.Hosts); i++ {
		match = g.Hosts[i].Match(grp.Hosts[i])
		if !match {
			return false
		}
	}

	// Compare Networks of groups.
	for i := 0; i < len(g.Networks); i++ {
		match = g.Networks[i].Match(grp.Networks[i])
		if !match {
			return false
		}
	}

	// Compare Ranges of groups.
	for i := 0; i < len(g.Ranges); i++ {
		match = g.Ranges[i].Match(grp.Ranges[i])
		if !match {
			return false
		}
	}

	// Compare Groups of groups.
	for i := 0; i < len(g.Groups); i++ {
		match = g.Groups[i].Match(grp.Groups[i])
		if !match {
			return false
		}
	}

	return true
}

// Add will add the specified object to the group.
// Supported types: Host/Network/Range/Group
func (g *Group) Add(obj interface{}) error {
	present, err := g.HasObject(obj)
	if err != nil {
		return fmt.Errorf("failed to check if object is already a group member: %v", err)
	}
	if present {
		return fmt.Errorf("object is already a member of this group: %s", obj)
	}

	switch v := obj.(type) {
	case *Host:
		g.addHost(v)
	case *Network:
		g.addNetwork(v)
	case *Range:
		g.addRange(v)
	case *Group:
		g.addGroup(v)
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
	// Ordered alphabetically by Group name.
	i := sort.Search(len(g.Groups), func(i int) bool {
		return g.Groups[i].Name >= grp.Name
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newGroup := make([]*Group, len(g.Groups)+1)
	// Shift the slice forward by one at the insert location.
	copy(newGroup[:i], g.Groups[:i])
	copy(newGroup[i+1:], g.Groups[i:])
	// Append group at the insert location.
	newGroup[i] = grp
	g.Groups = newGroup
}

// HasObject returns true if the group has a members object whose type and address matches the supplied object.
func (g *Group) HasObject(obj interface{}) (bool, error) {
	// TODO: Make more efficient since lists are now ordered.
	switch v := obj.(type) {
	case *Host:
		has := g.HasHost(v)
		if has {
			return has, nil
		}
	case *Network:
		has := g.HasNetwork(v)
		if has {
			return has, nil
		}
	case *Range:
		has := g.HasRange(v)
		if has {
			return has, nil
		}
	case *Group:
		has := g.HasGroup(v)
		if has {
			return has, nil
		}
	default:
		return false, fmt.Errorf("unsupported data type: %T", v)
	}

	// Check if any of it's group members contain the object.
	for _, grp := range g.Groups {
		has, err := grp.HasObject(obj)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

func (g *Group) HasHost(h *Host) bool {
	if len(g.Hosts) < 1 {
	    return false
    }
    var i int
    // Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
    if len(g.Hosts) == 1 {
	    i = 0
    } else {
	    i = sort.Search(len(g.Hosts), func(i int) bool {
		    keySt, keyEnd := h.Value()
		    midSt, midEnd := g.Hosts[i].Value()

		    return *keySt == *midSt && *keyEnd == *midEnd
	    })
    }

    // Check that what we go makes sense.
    if i == -1 || i >= len(g.Hosts) {
	    return false
    }

    // Double check that objects match.
    return g.Hosts[i].Match(h)
}

func (g *Group) HasNetwork(n *Network) bool {
	if len(g.Networks) < 1 {
		return false
	}

	var i int
    // Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
    if len(g.Networks) == 1 {
	    i = 0
    } else {
	    i = sort.Search(len(g.Networks), func(i int) bool {
		    keySt, keyEnd := n.Value()
		    midSt, midEnd := g.Networks[i].Value()

		    return *keySt == *midSt && *keyEnd == *midEnd
	    })
    }

    // Check that what we go makes sense.
    if i == -1 || i >= len(g.Networks) {
	    return false
    }

    // Double check that objects match.
    return g.Networks[i].Match(n)
}

func (g *Group) HasRange(r *Range) bool {
    var i int
    // Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
    if len(g.Ranges) == 1 {
	    i = 0
    } else {
	    i = sort.Search(len(g.Ranges), func(i int) bool {
		    keySt, keyEnd := r.Value()
		    midSt, midEnd := g.Ranges[i].Value()

		    return *keySt == *midSt && *keyEnd == *midEnd
	    })
    }

    // Check that what we go makes sense.
    if i == -1 || i >= len(g.Ranges) {
	    return false
    }

    // Double check that objects match.
    return g.Ranges[i].Match(r)


}

func (g *Group) HasGroup(grp *Group) bool {
	if len(g.Groups) < 1 {
		return false
	}

	var i int
	// Edge case handling. When len() == 0, sort.Search() was return index of 1 with is oob.
	if len(g.Groups) == 1 {
		i = 0
	} else {
		i = sort.Search(len(g.Groups), func(i int) bool {
			return g.Groups[i].Name == grp.Name
		})
	}

	// Check that what we go makes sense.
	if i == -1 || i >= len(g.Groups) {
		return false
	}

	// Double check that objects match.
	return g.Groups[i].Match(grp)
}

// this doesn't work because of the []NetworkObject, but keeping it in for now.
func binarySearch(left, right int, list []NetworkObject, key NetworkObject) int {
	if right >= left {
		mid := left - (right-left)/2

		keySt, keyEnd := key.Value()
		midSt, midEnd := list[mid].Value()

		if *keySt == *midSt {
			if *keyEnd == *midEnd {
				return mid
			}
			if *keyEnd > *midEnd {
				return binarySearch(mid+1, right, list, key)
			}
			return binarySearch(left, mid-1, list, key)
		}

		if *keySt > *midSt {
			return binarySearch(mid+1, right, list, key)
		}

		return binarySearch(left, mid-1, list, key)

	}

	return -1
}

func (g *Group) Contains(obj NetworkObject) bool {
	for _, h := range g.Hosts {
		if h.Contains(obj) {
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
