package core

import (
	"reflect"
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	testGroup := NewGroup("testGroup", "group for testing")
	var mux sync.Mutex

	t.Run("Add host", func(t *testing.T) {
		testHost, err := NewHost("testHost", "192.168.2.128", "test host")
		if err != nil {
			t.Fatalf("failed to create test host: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testHost)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testHost, testGroup.Hosts[0]) {
			t.Fatalf("host objects don't match")
		}
	})

	t.Run("Add host higher", func(t *testing.T) {
		testHost, err := NewHost("testHost", "192.168.2.129", "test host")
		if err != nil {
			t.Fatalf("failed to create test host: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testHost)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testHost, testGroup.Hosts[1]) {
			t.Fatalf("host objects don't match")
		}
	})

	t.Run("Add host lower", func(t *testing.T) {
		testHost, err := NewHost("testHost", "192.168.2.127", "test host")
		if err != nil {
			t.Fatalf("failed to create test host: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testHost)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		// Since list is sorted, this host should be at the start.
		if !reflect.DeepEqual(testHost, testGroup.Hosts[0]) {
			t.Fatalf("host objects don't match")
		}
	})

	t.Run("Add network", func(t *testing.T) {
		testNetwork, err := NewNetwork("testNetwork", "192.168.2.128", "255.255.255.128", "test network")
		if err != nil {
			t.Fatalf("failed to create test network: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testNetwork)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testNetwork, testGroup.Networks[0]) {
			t.Fatalf("network objects don't match")
		}
	})

	t.Run("Add network higher address", func(t *testing.T) {
		testNetwork, err := NewNetwork("testNetwork", "192.168.3.128", "255.255.255.128", "test network")
		if err != nil {
			t.Fatalf("failed to create test network: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testNetwork)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testNetwork, testGroup.Networks[1]) {
			t.Fatalf("network objects don't match")
		}
	})

	t.Run("Add network lower address", func(t *testing.T) {
		testNetwork, err := NewNetwork("testNetwork", "192.168.1.128", "255.255.255.128", "test network")
		if err != nil {
			t.Fatalf("failed to create test network: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testNetwork)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testNetwork, testGroup.Networks[0]) {
			t.Fatalf("network objects don't match")
		}
	})

	t.Run("Add network lower mask", func(t *testing.T) {
		testNetwork, err := NewNetwork("testNetwork", "192.168.1.128", "255.255.255.192", "test network")
		if err != nil {
			t.Fatalf("failed to create test network: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testNetwork)
		mux.Unlock()
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
		mux.Lock()
		err = testGroup.Add(testRange)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testRange, testGroup.Ranges[0]) {
			t.Fatalf("range objects don't match")
		}
	})

	t.Run("Add range higher start", func(t *testing.T) {
		testRange, err := NewRange("testRange", "192.168.2.130", "192.168.2.150", "test range")
		if err != nil {
			t.Fatalf("failed to create test range: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testRange)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testRange, testGroup.Ranges[1]) {
			t.Fatalf("range objects don't match")
		}
	})

	t.Run("Add range lower start", func(t *testing.T) {
		testRange, err := NewRange("testRange", "192.168.1.130", "192.168.2.150", "test range")
		if err != nil {
			t.Fatalf("failed to create test range: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testRange)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testRange, testGroup.Ranges[0]) {
			t.Fatalf("range objects don't match")
		}
	})

	t.Run("Add range higher end", func(t *testing.T) {
		testRange, err := NewRange("testRange", "192.168.2.128", "192.168.2.160", "test range")
		if err != nil {
			t.Fatalf("failed to create test range: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testRange)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testRange, testGroup.Ranges[3]) {
			t.Fatalf("range objects don't match")
		}
	})

	t.Run("Add range lower end", func(t *testing.T) {
		testRange, err := NewRange("testRange", "192.168.2.128", "192.168.2.140", "test range")
		if err != nil {
			t.Fatalf("failed to create test range: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testRange)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testRange, testGroup.Ranges[1]) {
			t.Fatalf("range objects don't match")
		}
	})

	t.Run("Add group", func(t *testing.T) {
		testAddGroup := NewGroup("EFGH", "test group")
		mux.Lock()
		err := testGroup.Add(testAddGroup)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testAddGroup, testGroup.Groups[0]) {
			t.Fatalf("group objects don't match")
		}
	})

	t.Run("Add group higher", func(t *testing.T) {
		testAddGroup := NewGroup("JKLM", "test group")
		mux.Lock()
		err := testGroup.Add(testAddGroup)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testAddGroup, testGroup.Groups[1]) {
			t.Fatalf("group objects don't match")
		}
	})

	t.Run("Add group lower", func(t *testing.T) {
		testAddGroup := NewGroup("ABCD", "test group")
		mux.Lock()
		err := testGroup.Add(testAddGroup)
		mux.Unlock()
		if err != nil {
			t.Fatalf("got error when not expected: %v", err)
		}

		if !reflect.DeepEqual(testAddGroup, testGroup.Groups[0]) {
			t.Fatalf("group objects don't match")
		}
	})

	t.Run("Host already present", func(t *testing.T) {
		testHost, err := NewHost("testHost", "192.168.2.128", "test host")
		if err != nil {
			t.Fatalf("failed to create test host: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testHost)
		mux.Unlock()
		if err != nil {
			// We expected the error so return
			return
		}

		t.Fatalf("expected error due to host already being a member of group")
	})

	t.Run("Network already present", func(t *testing.T) {
		testNetwork, err := NewNetwork("testNetwork", "192.168.2.128", "255.255.255.128", "test network")
		if err != nil {
			t.Fatalf("failed to create test network: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testNetwork)
		mux.Unlock()
		if err != nil {
			// We expected the error so return
			return
		}

		t.Fatalf("expected error due to network already being a member of group")
	})

	t.Run("Range already present", func(t *testing.T) {
		testRange, err := NewRange("testRange", "192.168.2.128", "192.168.2.150", "test range")
		if err != nil {
			t.Fatalf("failed to create test range: %v", err)
		}
		mux.Lock()
		err = testGroup.Add(testRange)
		mux.Unlock()
		if err != nil {
			// We expected the error so return
			return
		}

		t.Fatalf("expected error due to range already being a member of group")
	})

	t.Run("Group already present", func(t *testing.T) {
		testAddGroup := NewGroup("EFGH", "test group")
		mux.Lock()
		err := testGroup.Add(testAddGroup)
		mux.Unlock()
		if err != nil {
			// We expected the error so return
			return
		}

		t.Fatalf("expected error due to range already being a member of group")
	})

	t.Run("Unsupported type", func(t *testing.T) {
		invalid := "lorem ipsum"
		mux.Lock()
		err := testGroup.Add(invalid)
		mux.Unlock()
		if err == nil {
			t.Fatalf("didn't receive error when expected")
		}
	})
}

func TestHasObject(t *testing.T) {
	// Set up the objects we'll need, better to move to own function?
	host1, err := NewHost("host1", "192.168.1.1", "host 1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := NewHost("host2", "192.168.2.1", "host 2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	net1, err := NewNetwork("net1", "192.168.1.0", "255.255.255.0", "net 1")
	if err != nil {
		t.Fatalf("failed to create net1: %v", err)
	}
	net2, err := NewNetwork("net2", "192.168.1.128", "255.255.255.128", "net 1")
	if err != nil {
		t.Fatalf("failed to create net2: %v", err)
	}
	net3, err := NewNetwork("net2", "192.168.2.128", "255.255.255.128", "net 1")
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
	// range3, err := NewRange("range2", "192.168.2.1", "192.168.2.250", "range 2")
	// if err != nil {
	// 	t.Fatalf("failed to create range3: %v", err)
	// }
	testGroup := NewGroup("testGroup", "group for testing")
	testGroup2 := NewGroup("testGroup2", "group 2")

	// testGroup members:
	//   - host1
	//   - range1
	//   - net1
	//   - net2

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
	err = testGroup.Add(testGroup2)

	tests := []struct {
		name   string
		input  interface{}
		strict bool
		want   bool
		err    bool
	}{
		{name: "Host match", input: host1, want: true, err: false},
		{name: "Network match", input: net1, want: true, err: false},
		{name: "Range match", input: range1, want: true, err: false},
		{name: "Group match", input: testGroup2, want: true, err: false},
		{name: "Host no match", input: host2, want: false, err: false},
		{name: "Network no match", input: net3, want: false, err: false},
		{name: "Range no match", input: range2, want: false, err: false},
		{name: "Network match in group member", input: net2, want: true, err: false},
		{name: "Unsupported type", input: "lorem ipsum", want: false, err: true},
		{name: "Unsupported type", input: "lorem ipsum", want: false, err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := testGroup.HasObject(tc.input)
			if err != nil {
				// if we expected error then return, test was successful.
				if tc.err {
					return
				}
				t.Errorf("received error from test when not expected: %v", err)
			}
			if tc.err {
				t.Errorf("expected error from test but did not get one")
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
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
	net1, err := NewNetwork("net1", "192.168.1.0", "255.255.255.0", "net 1")
	if err != nil {
		t.Fatalf("failed to create net1: %v", err)
	}
	net2, err := NewNetwork("net2", "192.168.1.128", "255.255.255.128", "net 1")
	if err != nil {
		t.Fatalf("failed to create net2: %v", err)
	}
	net3, err := NewNetwork("net2", "192.168.2.128", "255.255.255.128", "net 1")
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

	// testGroup members:
	//   - host1
	//   - range1
	//   - net1
	//   - net2

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
	err = testGroup.Add(testGroup2)

	tests := []struct {
		name   string
		input  NetworkObject
		strict bool
		want   bool
		err    bool
	}{
		{name: "Not Strict - Host match", input: host1, want: true},
		{name: "Not Strict - Network match", input: net2, want: true},
		{name: "Not Strict - Range match", input: range2, want: true},
		{name: "Not Strict - Host no match", input: host2, want: false},
		{name: "Not Strict - Network no match", input: net3, want: false},
		{name: "Not Strict - Range no match", input: range3, want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := testGroup.Contains(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

/*

 */
