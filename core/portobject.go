package core

type PortObject interface {
	Value() (start int, end int, proto int)
}
