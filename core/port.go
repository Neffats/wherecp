package core

const (
	TCP = iota
	UDP
	ICMP
	ARP
)

type Port struct {
	UID      int
	Name     string
	PortNo   int
	Protocol int
	Comment  string
}
