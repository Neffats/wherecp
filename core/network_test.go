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
		{name: "Invalid IP address", input: "lorem ipsum", want: false, err: true},
		{name: "Host address same as network", input: "192.168.1.0", want: true, err: false},
		{name: "Invalid host format - range", input: "192.168.1.1-192.168.1.5", want: false, err: true},
		{name: "Invalid host format - network", input: "192.168.1.0/24", want: false, err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := netA.containsHost(tc.input)
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
		{name: "Invalid IP address", input: "lorem ipsum", want: false, err: true},
		{name: "Range start inside finish outside network", input: "192.168.1.5-192.168.2.3", want: false, err: false},
		{name: "Range start outside finish inside network", input: "192.168.0.5-192.168.1.33", want: false, err: false},
		{name: "Range same size as network", input: "192.168.1.0-192.168.1.255", want: true, err: false},
		{name: "Invalid range (start bigger than end)", input: "192.168.1.75-192.168.1.33", want: false, err: true},
		{name: "Invalid range format - network", input: "192.168.1.0/24", want: false, err: true},
		{name: "Invalid range format - host", input: "192.168.1.55", want: false, err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := netA.containsRange(tc.input)
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
		{name: "Invalid network mask", input: "192.168.0.0/41", want: false, err: true},
		{name: "Invalid network address", input: "999.999.999.999/20", want: false, err: true},
		{name: "Invalid network - host", input: "192.168.1.1", want: false, err: true},
		{name: "Invalid network - range", input: "192.168.1.1-192.168.5.2", want: false, err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := netA.containsNetwork(tc.input)
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
