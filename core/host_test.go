package core

import (
	"reflect"
	"testing"
)

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
			got := hostA.Match(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
