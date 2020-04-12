package core

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

func HasInSource() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Source
	}
}

func HasInDestination() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Destination
	}
}

func HasInService() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Port
	}
}

type NetContainser interface {
	Contains(obj NetworkObject) bool
}

func ContainsInSource() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Source
	}
}

func ContainsInDestination() func(*Rule) Haser {
	return func(r *Rule) Haser {
		return r.Destination
	}
}
