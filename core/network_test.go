package core

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewNetwork(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   bool
	}{
		{name: "Valid", input: "192.168.1.0/24", err: false},
		{name: "Different mask", input: "192.168.0.0/23", err: false},
		{name: "Different address", input: "192.168.2.0/24", err: false},
		{name: "Invalid mask", input: "192.168.2.0/33", err: true},
		{name: "Invalid address - letters", input: "lorem/ipsum", err: true},
		{name: "Invalid address - bad ip", input: "355.22.1.0/24", err: true},
		{name: "Invalid address - not network address", input: "192.168.1.2/24", err: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			_, err := NewNetwork("test", parts[0], parts[1], "test network")
			if err != nil {
				if tc.err == true {
					return
				}
				t.Fatalf("got error when not expected: %v", err)
			}
			if tc.err == true {
				t.Fatalf("expected error")
			}
		})
	}
}

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
		{name: "Different mask", input: "192.168.0.0/23", want: false},
		{name: "Different address", input: "192.168.2.0/24", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			in, err := NewNetwork("compareNet", parts[0], parts[1], "temp network for test")
			if err != nil {
				t.Fatalf("failed to create temp test network object: %v", err)
			}
			got := netA.Match(in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNetworkcontainsHost(t *testing.T) {
	netA, err := NewNetwork("netA", "192.168.1.0", "24", "test network")
	if err != nil {
		t.Fatalf("failed to create test network object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Host contained in Network", input: "192.168.1.3", want: true, err: false},
		{name: "Host outside network", input: "192.168.2.1", want: false, err: false},
		{name: "Host address same as network", input: "192.168.1.0", want: true, err: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in, err := NewHost("testHost", tc.input, "temp host object for test")
			if err != nil {
				t.Fatalf("failed to create test host object: %v", err)
			}
			got, err := netA.containsHost(in)
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("received error when not expected: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNetworkcontainsRange(t *testing.T) {
	netA, err := NewNetwork("netA", "192.168.1.0", "24", "test network")
	if err != nil {
		t.Fatalf("failed to create test network object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Range contained in Network", input: "192.168.1.3-192.168.1.6", want: true, err: false},
		{name: "Range outside network", input: "192.168.2.1-192.168.2.5", want: false, err: false},
		{name: "Range start inside finish outside network", input: "192.168.1.5-192.168.2.3", want: false, err: false},
		{name: "Range start outside finish inside network", input: "192.168.0.5-192.168.1.33", want: false, err: false},
		{name: "Range same size as network", input: "192.168.1.0-192.168.1.255", want: true, err: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "-")
			in, err := NewRange("testRange", parts[0], parts[1], "temp range object for test")
			if err != nil {
				t.Fatalf("failed to create test range object: %v", err)
			}
			got, err := netA.containsRange(in)
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("received error when not expected: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNetworkcontainsNetwork(t *testing.T) {
	netA, err := NewNetwork("netA", "192.168.1.0", "24", "test network")
	if err != nil {
		t.Fatalf("failed to create test network object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Matching network", input: "192.168.1.0/24", want: true, err: false},
		{name: "Contains network", input: "192.168.1.0/26", want: true, err: false},
		{name: "Outside of network", input: "192.168.5.0/24", want: false, err: false},
		{name: "Network that contains test network", input: "192.168.0.0/20", want: false, err: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			in, err := NewNetwork("compareNet", parts[0], parts[1], "temp network for test")
			if err != nil {
				t.Fatalf("failed to create temp test network object: %v", err)
			}
			got, err := netA.containsNetwork(in)
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("received error when not expected: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNetworkContains(t *testing.T) {
	netA, err := NewNetwork("netA", "192.168.1.0", "24", "test network")
	if err != nil {
		t.Fatalf("failed to create test network object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Network - inside test network", input: "192.168.1.0/24", want: true, err: false},
		{name: "Contains network", input: "192.168.1.0/26", want: true, err: false},
		{name: "Outside of network", input: "192.168.5.0/24", want: false, err: false},
		{name: "Network that contains test network", input: "192.168.0.0/20", want: false, err: false},
	}

}
