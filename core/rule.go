package core

import (
	"github.com/google/uuid"
)

// Rule is a representation of a firewall rule.
type Rule struct {
	uid         string
	number      int
	source      *Group
	destination *Group
	port        *PortGroup
	action      bool
	comment     string
}

// NewRule returns a pointer to a new Rule object.
func NewRule(number int, src, dst *Group, prt *PortGroup, action bool, comment string) *Rule {
	uid := uuid.New()
	return &Rule{
		uid: uid.String(),
		number:      number,
		source:      src,
		destination: dst,
		port:        prt,
		action:      action,
		comment:     comment,
	}
}

type Haser interface {
	HasObject(obj interface{}) (bool, error)
}

func HasInSource() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.source
	}
}

func HasInDestination() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.destination
	}
}

func merge(name string, g1, g2 *Group) *Group {
	merged := NewGroup(name, "")
	merged.addGroup(g1)
	merged.addGroup(g2)
	return merged
}

func HasInAny() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return merge("Any", r.source, r.destination)
	}
}

func HasInService() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.port
	}
}

type NetContainser interface {
	Contains(obj NetworkObject) bool
}

func ContainsInSource() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.source
	}
}

func ContainsInDestination() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.destination
	}
}
