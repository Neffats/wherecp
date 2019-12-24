package core

type Group struct {
	UID      int
	Name     string
	Hosts    []*Host
	Networks []*Network
	Ranges   []*Range
	Groups   []*Group
	Comment  string
}

func NewGroup(name, comment string) *Group {
	return &Group{
		UID:      0,
		Name:     name,
		Hosts:    make([]*Host, 1),
		Networks: make([]*Network, 1),
		Ranges:   make([]*Range, 1),
		Groups:   make([]*Group, 1),
		Comment:  comment,
	}
}
