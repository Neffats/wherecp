package core

import (
	"reflect"
	"testing"
)

func TestNetworkMatch(t *testing.T) {
	netA, err := NewNetwork("netA", "192.168.1.0", "24", "test network")
	if err != nil {
		t.Fatalf("failed to create test network object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "Positive Match", input: "192.168.1.0/24", want: true},
		{name: "Different mask", input: "192.168.1.0/23", want: false},
		{name: "Different address", input: "192.168.2.0/24", want: false},
		{name: "Not a network", input: "192.168.1.1", want: false},
		{name: "Invalid format", input: "lorem ipsum", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := netA.Match(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
