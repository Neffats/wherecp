package core

import (
	"reflect"
	"testing"
)

func TestNewHost(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   bool
	}{
		{name: "Valid host", input: "10.10.10.10", err: false},
		{name: "Invalid host - bad adddress", input: "355.32.1.3", err: true},
		{name: "Invalid host - letters", input: "lorem ipsum", err: true},
		{name: "Invalid host - network", input: "10.10.10.0/24", err: true},
		{name: "Invalid host - range", input: "10.10.10.10-10.10.10.20", err: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewHost("test", tc.input, "test host")
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

func TestHostMatch(t *testing.T) {
	hostA, err := NewHost("hostA", "10.10.10.10", "test host")
	if err != nil {
		t.Fatalf("failed to create test host object")
	}
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "True match", input: "10.10.10.10", want: true},
		{name: "No match", input: "10.10.0.0", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testHost, err := NewHost("testHost", tc.input, "test host")
			if err != nil {
				t.Fatalf("failed to create test host: %v", err)
			}
			got := hostA.Match(testHost)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
