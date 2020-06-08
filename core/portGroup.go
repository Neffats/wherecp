package core

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/google/uuid"
)

// PortGroup groups together different Port, PortRanges and other PortGroup objects.
type PortGroup struct {
	uid     string
	name    string
	ports   []*Port
	ranges  []*PortRange
	groups  []*PortGroup
	comment string
}

// NewPortGroup returns a new empty oort group.
func NewPortGroup(name, comment string) *PortGroup {
	uid := uuid.New()
	return &PortGroup{
		uid:     uid.String(),
		name:    name,
		ports:   make([]*Port, 0),
		ranges:  make([]*PortRange, 0),
		groups:  make([]*PortGroup, 0),
		comment: comment,
	}
}

// Add will add the specified object to the group.
// Supported types: Port/Port Range/Port Group
func (pg *PortGroup) Add(obj interface{}) error {
	present, err := pg.HasObject(obj)
	if err != nil {
		return fmt.Errorf("failed to check if object is already a group member: %v", err)
	}
	if present {
		return fmt.Errorf("object is already a member of this group: %s", obj)
	}

	switch v := obj.(type) {
	case *Port:
		pg.addPort(v)
	case *PortRange:
		pg.addPortRange(v)
	case *PortGroup:
		pg.addPortGroup(v)
	default:
		return errors.New("unsupported data type")
	}
	return nil
}

func (pg *PortGroup) addPort(p *Port) {
	// Ordered smallest to largets, first by port number then by protocol number.
	i := sort.Search(len(pg.ports), func(i int) bool {
		return pg.ports[i].number > p.number && pg.ports[i].protocol >= p.protocol
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newPorts := make([]*Port, len(pg.ports)+1)
	// Shift the slice forward by one at the insert location.
	copy(newPorts[:i], pg.ports[:i])
	copy(newPorts[i+1:], pg.ports[i:])
	// Append port at the insert location.
	newPorts[i] = p
	pg.ports = newPorts
}

func (pg *PortGroup) addPortRange(pr *PortRange) {
	// Ordered smallest to largest by start range start (first port) first
	// then by range end (last port) second, then by protocol number. Smaller ranges will come before
	// larger ranges i.e. 1-2, 1-5, 2-4
	i := sort.Search(len(pg.ranges), func(i int) bool {
		thisStart, thisEnd, thisProto := pg.ranges[i].Value()
		otherStart, otherEnd, otherProto := pr.Value()

		start := thisStart >= otherStart
		end := thisEnd >= otherEnd
		proto := thisProto >= otherProto
		return start && end && proto
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newRanges := make([]*PortRange, len(pg.ranges)+1)
	// Shift the slice forward by one at the insert location.
	copy(newRanges[:i], pg.ranges[:i])
	copy(newRanges[i+1:], pg.ranges[i:])
	// Append port at the insert location.
	newRanges[i] = pr
	pg.ranges = newRanges
}

func (pg *PortGroup) addPortGroup(grp *PortGroup) {
	// Ordered alphabetically by Group name.
	i := sort.Search(len(pg.groups), func(i int) bool {
		return pg.groups[i].name >= grp.name
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newGroup := make([]*PortGroup, len(pg.groups)+1)
	// Shift the slice forward by one at the insert location.
	copy(newGroup[:i], pg.groups[:i])
	copy(newGroup[i+1:], pg.groups[i:])
	// Append group at the insert location.
	newGroup[i] = grp
	pg.groups = newGroup
}

// HasObject returns true if the group has a members object whose type and address matches the supplied object.
func (pg *PortGroup) HasObject(obj interface{}) (bool, error) {
	// TODO: Make more efficient since lists are now ordered.
	switch v := obj.(type) {
	case *Port:
		if len(pg.ports) < 1 {
			return false, nil
		}
		var i int
		// Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
		if len(pg.ports) == 1 {
			i = 0
		} else {
			i = sort.Search(len(pg.ports), func(i int) bool {
				keySt, keyEnd, keyProto := v.Value()
				midSt, midEnd, midProto := pg.ports[i].Value()

				return keySt == midSt && keyEnd == midEnd && keyProto == midProto
			})
		}

		// Check that what we go makes sense.
		if i == -1 || i >= len(pg.ports) {
			return false, nil
		}

		// Double check that objects match.
		if pg.ports[i].Match(v) {
			return true, nil
		}
	case *PortRange:
		var i int
		// Edge case handling. When len() == 0, sort.Search() was return index of 1 with is oob.
		if len(pg.ranges) == 1 {
			i = 0
		} else {
			i = sort.Search(len(pg.ranges), func(i int) bool {
				keySt, keyEnd, keyProto := v.Value()
				midSt, midEnd, midProto := pg.ranges[i].Value()

				return keySt == midSt && keyEnd == midEnd && keyProto == midProto
			})
		}

		// Check that what we go makes sense.
		if i == -1 || i >= len(pg.ranges) {
			return false, nil
		}

		// Double check that objects match.
		if pg.ranges[i].Match(v) {
			return true, nil
		}
	case *PortGroup:
		var i int
		// Edge case handling. When len() == 0, sort.Search() was return index of 1 with is oob.
		if len(pg.groups) == 1 {
			i = 0
		} else {
			i = sort.Search(len(pg.groups), func(i int) bool {
				return pg.groups[i].name == v.name
			})
		}

		// Check that what we go makes sense.
		if i == -1 || i >= len(pg.groups) {
			return false, nil
		}

		// Double check that objects match.
		if pg.groups[i].Match(v) {
			return true, nil
		}
	default:
		return false, errors.New("unsupported data type")
	}

	// Check if any of it's group members contain the object.
	for _, grp := range pg.groups {
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

// Contains will return true if a port object is contained by a member in the group. 
func (pg *PortGroup) Contains(obj PortObject) bool {
	for _, p := range pg.ports {
		if p.Contains(obj) {
			return true
		}
	}
	for _, r := range pg.ranges {
		if r.Contains(obj) {
			return true
		}
	}
	for _, g := range pg.groups {
		if g.Contains(obj) {
			return true
		}
	}
	return false
}

// Match will return true if both groups are identical.
func (pg *PortGroup) Match(grp *PortGroup) bool {
	return reflect.DeepEqual(pg, grp)
}

// MatchContent returns true if both groups contain the exact same members.
func (pg *PortGroup) MatchContent(grp *PortGroup) bool {
	// Check if the lengths of the groups match.
	// If they don't then the two groups must be different.
	if len(pg.ports) != len(grp.ports) {
		return false
	}
	if len(pg.ranges) != len(grp.ranges) {
		return false
	}
	if len(pg.groups) != len(grp.groups) {
		return false
	}

	var match bool

	// Compare Hosts of groups.
	// All group members are sorted, so all members should be in the same location.
	for i := 0; i < len(pg.ports); i++ {
		match = pg.ports[i].Match(grp.ports[i])
		if !match {
			return false
		}
	}

	// Compare Ranges of both groups.
	for i := 0; i < len(pg.ranges); i++ {
		match = pg.ranges[i].Match(grp.ranges[i])
		if !match {
			return false
		}
	}

	// Compare Groups of groups.
	for i := 0; i < len(pg.groups); i++ {
		match = pg.groups[i].Match(grp.groups[i])
		if !match {
			return false
		}
	}

	return true
}
