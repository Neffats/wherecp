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
	range1, err := NewRange("range1", "192.168.1.1", "192.168.1.250", "range 1")
	if err != nil {
		t.Fatalf("failed to create range1: %v", err)
	}
	range2, err := NewRange("range2", "192.168.2.1", "192.168.2.250", "range 2")
	if err != nil {
		t.Fatalf("failed to create range2: %v", err)
	}
	testGroup, err := NewGroup("testGroup", "group for testing")
	if err != nil {
		t.Fatalf("failed to create test group: %v", err)
	}

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

	tests := []struct {
		name   string
		input  interface{}
		strict bool
		want   bool
		err    bool
	}{
		{name: "Strict - Host match", input: host1, strict: true, want: true, err: false},
		{name: "Strict - Network match", input: net1, strict: true, want: true, err: false},
		{name: "Strict - Range match", input: range1, strict: true, want: true, err: false},
		{name: "Strict - Host no match", input: host2, strict: true, want: false, err: false},
		{name: "Strict - Network no match", input: net2, strict: true, want: false, err: false},
		{name: "Strict - Range no match", input: range2, strict: true, want: false, err: false},
		{name: "Strict - Unsupported type", input: "lorem ipsum", strict: true, want: false, err: true},
	}
}
