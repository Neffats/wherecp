package core

import (
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

const defaultGroupCapacity = 100

// Group is a structure that groups different object together.
// Acting as a container for different objects. Each of the item arrays are ordered
// for efficient searching.
// The Network objects are ordered by Address, and the Groups
// are ordered by name.
type Group struct {
	uid      string
	name     string
	hosts    []*Host
	networks []*Network
	ranges   []*Range
	groups   []*Group
	comment  string
}

// NewGroup returns a new empty group.
func NewGroup(name, comment string) *Group {
	uid := uuid.New()
	return &Group{
		uid:      uid.String(),
		name:     name,
		hosts:    make([]*Host, 0),
		networks: make([]*Network, 0),
		ranges:   make([]*Range, 0),
		groups:   make([]*Group, 0),
		comment:  comment,
	}
}

// Match will return true if the two groups are identical.
func (g *Group) Match(grp *Group) bool {
	if grp.name == g.name && g.MatchContent(grp) {
		return true
	}
	return false
}

// MatchContent will return true if both groups contain the same members.
func (g *Group) MatchContent(grp *Group) bool {
	// Check if the lengths of the groups match.
	// If they don't then the two groups must be different.
	if len(g.hosts) != len(grp.hosts) {
		return false
	}
	if len(g.networks) != len(grp.networks) {
		return false
	}
	if len(g.ranges) != len(grp.ranges) {
		return false
	}
	if len(g.groups) != len(grp.groups) {
		return false
	}

	var match bool

	// Compare Hosts of groups.
	// All group members are sorted, so all members should be in the same location.
	for i := 0; i < len(g.hosts); i++ {
		match = g.hosts[i].Match(grp.hosts[i])
		if !match {
			return false
		}
	}

	// Compare Networks of groups.
	for i := 0; i < len(g.networks); i++ {
		match = g.networks[i].Match(grp.networks[i])
		if !match {
			return false
		}
	}

	// Compare Ranges of groups.
	for i := 0; i < len(g.ranges); i++ {
		match = g.ranges[i].Match(grp.ranges[i])
		if !match {
			return false
		}
	}

	// Compare Groups of groups.
	for i := 0; i < len(g.groups); i++ {
		match = g.groups[i].Match(grp.groups[i])
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
	i := sort.Search(len(g.hosts), func(i int) bool {
		return *g.hosts[i].Address > *h.Address
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newHosts := make([]*Host, len(g.hosts)+1)
	// Shift the slice forward by one at the insert location.
	copy(newHosts[:i], g.hosts[:i])
	copy(newHosts[i+1:], g.hosts[i:])
	// Append host at the insert location.
	newHosts[i] = h
	g.hosts = newHosts
}

func (g *Group) addNetwork(n *Network) {
	// Ordered smallest to largest by network address (first address) first
	// then by broadcast address (last address) second. Smallest networks will be in
	// front of larger networks i.e. 192.168.0.0/25 will be before 192.168.0.0/24
	i := sort.Search(len(g.networks), func(i int) bool {
		this := g.networks[i].Unpack()
		other := n.Unpack()

		addr := this[0].Start >= other[0].Start
		mask := this[0].End >= other[0].End
		return addr && mask
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newNets := make([]*Network, len(g.networks)+1)
	// Shift the slice forward by one at the insert location.
	copy(newNets[:i], g.networks[:i])
	copy(newNets[i+1:], g.networks[i:])
	// Append network at the insert location.
	newNets[i] = n
	g.networks = newNets
}

func (g *Group) addRange(r *Range) {
	// Ordered smallest to largest by start address (first address) first
	// then by end address (last address) second. Smaller ranges will come before
	// larger ranges i.e. 192.168.0.0-192.168.0.10 will be in front of 192.168.0.0-192.168.0.200
	i := sort.Search(len(g.ranges), func(i int) bool {
		this := g.ranges[i].Unpack()
		other := r.Unpack()

		start := this[0].Start >= other[0].Start
		end := this[0].End >= other[0].End
		return start && end
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newRange := make([]*Range, len(g.ranges)+1)
	// Shift the slice forward by one at the insert location.
	copy(newRange[:i], g.ranges[:i])
	copy(newRange[i+1:], g.ranges[i:])
	// Append network at the insert location.
	newRange[i] = r
	g.ranges = newRange
}

func (g *Group) addGroup(grp *Group) {
	// Ordered alphabetically by Group name.
	i := sort.Search(len(g.groups), func(i int) bool {
		return g.groups[i].name >= grp.name
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newGroup := make([]*Group, len(g.groups)+1)
	// Shift the slice forward by one at the insert location.
	copy(newGroup[:i], g.groups[:i])
	copy(newGroup[i+1:], g.groups[i:])
	// Append group at the insert location.
	newGroup[i] = grp
	g.groups = newGroup
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
	for _, grp := range g.groups {
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
	if len(g.hosts) < 1 {
	    return false
    }
    var i int
    // Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
    if len(g.hosts) == 1 {
	    i = 0
    } else {
	    i = sort.Search(len(g.hosts), func(i int) bool {
		    return h.Match(g.hosts[i])
	    })
    }

    // Check that what we go makes sense.
    if i == -1 || i >= len(g.hosts) {
	    return false
    }

    // Double check that objects match.
    return g.hosts[i].Match(h)
}

func (g *Group) HasNetwork(n *Network) bool {
	if len(g.networks) < 1 {
		return false
	}

	var i int
    // Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
    if len(g.networks) == 1 {
	    i = 0
    } else {
	    i = sort.Search(len(g.networks), func(i int) bool {
		    return n.Match(g.networks[i])
	    })
    }

    // Check that what we go makes sense.
    if i == -1 || i >= len(g.networks) {
	    return false
    }

    // Double check that objects match.
    return g.networks[i].Match(n)
}

func (g *Group) HasRange(r *Range) bool {
    var i int
    // Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
    if len(g.ranges) == 1 {
	    i = 0
    } else {
	    i = sort.Search(len(g.ranges), func(i int) bool {
		    return r.Match(g.ranges[i])
	    })
    }

    // Check that what we go makes sense.
    if i == -1 || i >= len(g.ranges) {
	    return false
    }

    // Double check that objects match.
    return g.ranges[i].Match(r)


}

func (g *Group) HasGroup(grp *Group) bool {
	if len(g.groups) < 1 {
		return false
	}

	var i int
	// Edge case handling. When len() == 0, sort.Search() was return index of 1 with is oob.
	if len(g.groups) == 1 {
		i = 0
	} else {
		i = sort.Search(len(g.groups), func(i int) bool {
			return g.groups[i].name == grp.name
		})
	}

	// Check that what we go makes sense.
	if i == -1 || i >= len(g.groups) {
		return false
	}

	// Double check that objects match.
	return g.groups[i].Match(grp)
}

func (g *Group) Contains(obj NetworkUnpacker) bool {
	for _, h := range g.hosts {
		if h.Contains(obj) {
			return true
		}
	}
	for _, n := range g.networks {
		if n.Contains(obj) {
			return true
		}
	}
	for _, r := range g.ranges {
		if r.Contains(obj) {
			return true
		}
	}
	for _, grp := range g.groups {
		if grp.Contains(obj) {
			return true
		}
	}
	return false
}
