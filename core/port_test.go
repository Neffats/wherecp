package core

import "testing"

func TestString2Proto(t *testing.T) {
	tests := []struct {
		name     string
		in       string 
		want     int
	}{
		{name: "Valid-tcp", in: "tcp", want: 0},
		{name: "Upper-case-tcp", in: "TCP", want: 0},
		{name: "Valid-icmp", in: "icmp", want: 2},
		{name: "Invalid-proto", in: "INVALID", want: -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := String2Proto(tc.in)
			if got != tc.want {
				t.Fatalf("want: %d, got: %d", tc.want, got)
			}
		})
	}
}

func TestProto2String(t *testing.T) {
	tests := []struct {
		name     string
		in       int
		want     string
	}{
		{name: "Valid-tcp", in: 0, want: "tcp"},
		{name: "Valid-icmp", in: 2, want: "icmp"},
		{name: "Invalid-proto-too-high", in: 10, want: ""},
		{name: "Invalid-proto-negative-number", in: -1, want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Proto2String(tc.in)
			if got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

func TestNewPort(t *testing.T) {
	tests := []struct {
		name     string
		number   uint
		protocol string
		err      bool
	}{
		{name: "valid tcp port", number: 443, protocol: "tcp", err: false},
		{name: "invalid protocol", number: 443, protocol: "xma", err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPort("test", tc.number, tc.protocol, "test port")
			if err != nil {
				if tc.err == true {
					return
				}
				t.Fatalf("got error when not expected: %v", err)
			}
			if tc.err == true {
				t.Fatalf("expected error, but got none")
			}
		})
	}
}

func TestMatchPort(t *testing.T) {
	port, err := NewPort("Test Port HTTPS", 443, "TCP", "Port for testing")
	if err != nil {
		t.Fatalf("failed to create test port: %v", err)
	}
	tests := []struct {
		name     string
		number   uint
		protocol string
		want     bool
	}{
		{name: "matching-port-proto", number: 443, protocol: "tcp", want: true},
		{name: "Proto-not-matching", number: 443, protocol: "udp", want: false},
		{name: "Port-number-not-matching", number: 22, protocol: "tcp", want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testPort, err := NewPort("Test Port", tc.number, tc.protocol, "Port for testing")
			if err != nil {
				t.Fatalf("failed to create test port: %v", err)
			}
			got := port.Match(testPort)
			if got != tc.want {
				t.Fatalf("want: %t got: %t", tc.want, got)
			}
		})
	}
}
