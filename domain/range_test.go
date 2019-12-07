package core

import (
	"reflect"
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
		{name: "Invalid range", input: "lorem ipsum", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := rangeA.Match(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}
