package core

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	testGroup := NewGroup("testGroup", "group for testing")

	t.Run("Add host", func(t *testing.T) {
		testHost, err := NewHost("testHost", "192.168.2.128", "test host")
		if err != nil {
			t.Fatalf("failed to create test host: %v", err)
		}
		err = testGroup.Add(testHost)
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testHost, testGroup.Hosts[0]) {
			t.Fatalf("host objects don't match")
		}
	})

	t.Run("Add network", func(t *testing.T) {
		testNetwork, err := NewNetwork("testNetwork", "192.168.2.128", "25", "test network")
		if err != nil {
			t.Fatalf("failed to create test network: %v", err)
		}
		err = testGroup.Add(testNetwork)
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testNetwork, testGroup.Networks[0]) {
			t.Fatalf("network objects don't match")
		}
	})

	t.Run("Add range", func(t *testing.T) {
		testRange, err := NewRange("testRange", "192.168.2.128", "192.168.2.150", "test range")
		if err != nil {
			t.Fatalf("failed to create test range: %v", err)
		}
		err = testGroup.Add(testRange)
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testRange, testGroup.Ranges[0]) {
			t.Fatalf("range objects don't match")
		}
	})

	t.Run("Unsupported type", func(t *testing.T) {
		invalid := "lorem ipsum"
		err := testGroup.Add(invalid)
		if err == nil {
			t.Fatalf("didn't receive error when expected")
		}
	})
}

func TestGroupContains(t *testing.T) {
	// Set up the objects we'll need, better to move to own function?
	host1, err := NewHost("host1", "192.168.1.1", "host 1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := NewHost("host2", "192.168.2.1", "host 2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	net1, err := NewNetwork("net1", "192.168.1.0", "24", "net 1")
	if err != nil {
		t.Fatalf("failed to create net1: %v", err)
	}
	net2, err := NewNetwork("net2", "192.168.1.128", "25", "net 1")
	if err != nil {
		t.Fatalf("failed to create net2: %v", err)
	}
	net3, err := NewNetwork("net2", "192.168.2.128", "25", "net 1")
	if err != nil {
		t.Fatalf("failed to create net3: %v", err)
	}
	range1, err := NewRange("range1", "192.168.1.1", "192.168.1.250", "range 1")
	if err != nil {
		t.Fatalf("failed to create range1: %v", err)
	}
	range2, err := NewRange("range2", "192.168.1.50", "192.168.1.150", "range 2")
	if err != nil {
		t.Fatalf("failed to create range2: %v", err)
	}
	range3, err := NewRange("range2", "192.168.2.1", "192.168.2.250", "range 2")
	if err != nil {
		t.Fatalf("failed to create range3: %v", err)
	}
	testGroup := NewGroup("testGroup", "group for testing")
	testGroup2 := NewGroup("testGroup2", "group 2")

	err = testGroup.Add(host1)
	if err != nil {
		t.Fatalf("failed to add host1 to group: %v", err)
	}
	err = testGroup.Add(range1)
	if err != nil {
		t.Fatalf("failed to add range1 to group: %v", err)
	}
	err = testGroup.Add(net1)
	if err != nil {
		t.Fatalf("failed to add net1 to group: %v", err)
	}
	err = testGroup2.Add(net2)
	if err != nil {
		t.Fatalf("failed to add net2 to group2: %v", err)
	}

	tests := []struct {
		name   string
		input  interface{}
		strict bool
		want   bool
	}{
		{name: "Strict - Host match", input: host1, want: true},
		{name: "Strict - Network match", input: net1, want: true},
		{name: "Strict - Range match", input: range1, want: true},
		//{name: "Strict - Group match", input: testGroup,  want: true, },
		{name: "Strict - Host no match", input: host2, want: false},
		{name: "Strict - Network no match", input: net2, want: false},
		{name: "Strict - Range no match", input: range2, want: false},
		{name: "Strict - Unsupported type", input: "lorem ipsum", want: false},
		{name: "Not Strict - Host match", input: host1, want: true},
		{name: "Not Strict - Network match", input: net2, want: true},
		{name: "Not Strict - Range match", input: range2, want: true},
		{name: "Not Strict - Host no match", input: host2, want: false},
		{name: "Not Strict - Network no match", input: net3, want: false},
		{name: "Not Strict - Range no match", input: range3, want: false},
		{name: "Strict - Unsupported type", input: "lorem ipsum", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := testGroup.HasObject(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
