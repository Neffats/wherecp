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

type Haser interface {
	HasObject(obj interface{}) (bool, error)
}

func InSource() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Source
	}
}

func InDestination() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Destination
	}
}

func InService() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Port
	}
}

func (r *Rule) Has(obj interface{}, comp func(*Rule) Haser) (bool, error) {
	component := comp(r)
	has, err := component.HasObject(obj)
	if err != nil {
		return false, fmt.Errorf("failed to determine if object is in rule: %v", err)
	}

	return has, nil
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

func (r *Rule) ContainsSource(obj NetworkObject) bool {
	return r.Source.Contains(obj)
}

func (r *Rule) ContainsDestination(obj NetworkObject) bool {
	return r.Destination.Contains(obj)
}

func (r *Rule) ContainsPort(obj PortObject) bool {
	return r.Port.Contains(obj)
}
