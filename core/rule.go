package core

import (
	"fmt"
)

// Rule is a representation of a firewall rule.
type Rule struct {
	UID         string
	Source      *Group
	Destination *Group
	Port        *PortGroup
	Action      bool
	Comment     string
}

// NewRule returns a pointer to a new Rule object.
func NewRule(src, dst *Group, prt *PortGroup, action bool, comment string) *Rule {
	return &Rule{
		Source:      src,
		Destination: dst,
		Port:        prt,
		Action:      action,
		Comment:     comment,
	}
}

// HasSource returns true if the rule contains the object in it's source group.
func (r *Rule) HasSource(obj interface{}) (bool, error) {
	has, err := r.Source.HasObject(obj)
	if err != nil {
		return false, fmt.Errorf("couldn't determine whether rule has source because: %v", err)
	}
	return has, nil
}

// HasDestination returns true if the rule contains the object in it's destination group.
func (r *Rule) HasDestination(obj interface{}) (bool, error) {
	has, err := r.Destination.HasObject(obj)
	if err != nil {
		return false, fmt.Errorf("couldn't determine whether rule has destination because: %v", err)
	}
	return has, nil
}

// HasPort returns true if the rule contains the object in it's port group.
func (r *Rule) HasPort(obj interface{}) (bool, error) {
	has, err := r.Port.HasObject(obj)
	if err != nil {
		return false, fmt.Errorf("couldn't determine whether rule has port because: %v", err)
	}
	return has, nil
}
