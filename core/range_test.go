package core

import (
	"reflect"
	"strings"
	"testing"
)

func TestRangeMatch(t *testing.T) {
	rangeA, err := NewRange("rangeA", "192.168.1.1", "192.168.254.1", "test range object")
	if err != nil {
		t.Fatalf("failed to create test range object: %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "Valid match", input: "192.168.1.1-192.168.254.1", want: true},
		{name: "Different range", input: "10.10.10.10-10.10.20.10", want: false},
		{name: "Smaller range", input: "192.168.1.240-192.168.10.254", want: false},
		{name: "Bigger range", input: "192.167.40.1-192.169.1.1", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "-")
			in, err := NewRange("testRange", parts[0], parts[1], "temp range object for test")
			if err != nil {
				t.Fatalf("failed to create test range object: %v", err)
			}
			got := rangeA.Match(in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestRangecontainsHost(t *testing.T) {
	rangeA, err := NewRange("rangeA", "192.168.1.1", "192.168.1.254", "test range")
	if err != nil {
		t.Fatalf("failed to create test range object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Host contained in range", input: "192.168.1.3", want: true, err: false},
		{name: "Host outside range", input: "192.168.2.1", want: false, err: false},
		{name: "Invalid IP address", input: "lorem ipsum", want: false, err: true},
		{name: "Host address same as range", input: "192.168.1.1", want: true, err: false},
		{name: "Invalid host format - range", input: "192.168.1.1-192.168.1.5", want: false, err: true},
		{name: "Invalid host format - network", input: "192.168.1.0/24", want: false, err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in, err := NewHost("testHost", tc.input, "temp host object for test")
			if err != nil {
				t.Fatalf("failed to create test host object: %v", err)
			}
			got, err := rangeA.containsHost(in)
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

func TestRangecontainsRange(t *testing.T) {
	rangeA, err := NewRange("rangeA", "192.168.1.1", "192.168.1.254", "test range")
	if err != nil {
		t.Fatalf("failed to create test range object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Range contained in range", input: "192.168.1.3-192.168.1.6", want: true, err: false},
		{name: "Range outside range", input: "192.168.2.1-192.168.2.5", want: false, err: false},
		{name: "Invalid IP address", input: "lorem ipsum", want: false, err: true},
		{name: "Range start inside finish outside range", input: "192.168.1.5-192.168.2.3", want: false, err: false},
		{name: "Range start outside finish inside range", input: "192.168.0.5-192.168.1.33", want: false, err: false},
		{name: "Range same size as range", input: "192.168.1.1-192.168.1.254", want: true, err: false},
		{name: "Invalid range (start bigger than end)", input: "192.168.1.75-192.168.1.33", want: false, err: true},
		{name: "Invalid range format - network", input: "192.168.1.0/24", want: false, err: true},
		{name: "Invalid range format - host", input: "192.168.1.55", want: false, err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "-")
			in, err := NewRange("testRange", parts[0], parts[1], "temp range object for test")
			if err != nil {
				t.Fatalf("failed to create test range object: %v", err)
			}
			got, err := rangeA.containsRange(in)
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

func TestRangecontainsNetwork(t *testing.T) {
	rangeA, err := NewRange("rangeA", "192.168.1.0", "192.168.1.255", "test range")
	if err != nil {
		t.Fatalf("failed to create test range object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Matching network", input: "192.168.1.0/24", want: true, err: false},
		{name: "Contains network", input: "192.168.1.0/26", want: true, err: false},
		{name: "Outside of range", input: "192.168.5.0/24", want: false, err: false},
		{name: "Network that contains test range", input: "192.168.0.0/20", want: false, err: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			in, err := NewNetwork("compareNet", parts[0], parts[1], "temp network for test")
			if err != nil {
				t.Fatalf("failed to create temp test network object: %v", err)
			}
			got, err := rangeA.containsNetwork(in)
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
