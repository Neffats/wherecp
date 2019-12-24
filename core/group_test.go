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
		err := testGroup.Add(testHost)
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
		err := testGroup.Add(testNetwork)
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testNetwork, testGroup.Networks[0]) {
			t.Fatalf("host objects don't match")
		}
	})

	t.Run("Add range", func(t *testing.T) {
		testRange, err := NewRange("testRange", "192.168.2.128", "192.168.2.150", "test range")
		if err != nil {
			t.Fatalf("failed to create test range: %v", err)
		}
		err := testGroup.Add(testRange)
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testRange, testGroup.Networks[0]) {
			t.Fatalf("host objects don't match")
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
