package core

type PortObject interface {
	Value() (start uint, end uint, proto int)
}
