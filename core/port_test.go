package core

import "testing"

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
