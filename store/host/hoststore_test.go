package hoststore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Neffats/wherecp/core"
)

type testPuller struct {}

func (tp *testPuller) PullHosts() ([]*core.Host, error) {
	return make([]*core.Host, 0), nil
}

func TestAll(t *testing.T) {
	// Setup the test data.
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	host3, err := core.NewHost("host3", "192.168.51.1", "host3")
	if err != nil {
		t.Fatalf("failed to create host3: %v", err)
	}

	hosts := make([]*core.Host, 0)
	hosts = append(hosts, host1)
	hosts = append(hosts, host2)
	hosts = append(hosts, host3)
	
	testStore := &HostStore{
		Hosts: hosts,
		Puller: &testPuller{},
	}

	got := testStore.All()
	if diff := reflect.DeepEqual(got, hosts); !diff {
		t.Errorf("want: %+v\ngot: %+v", hosts, got)
	}
}

func TestGet(t *testing.T) {
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	host3, err := core.NewHost("host3", "192.168.51.1", "host3")
	if err != nil {
		t.Fatalf("failed to create host3: %v", err)
	}
	host4, err := core.NewHost("host4", "10.10.51.1", "host3")
	if err != nil {
		t.Fatalf("failed to create host4: %v", err)
	}

	hosts := make([]*core.Host, 0)
	hosts = append(hosts, host1)
	hosts = append(hosts, host2)
	hosts = append(hosts, host3)
	
	testStore := &HostStore{
		Hosts: hosts,
		Puller: &testPuller{},
	}

	tests := []struct {
		name  string
		input string
		want  *core.Host
		err   bool
	}{
		{name: "Get host 1",
			input: host1.UID(),
			want: host1,
			err:   false},
		{name: "Get host 2",
			input: host2.UID(),
			want: host2,
			err:   false},
		{name: "Non-exitent host",
			input: host4.UID(),
			want: nil,
			err:   true},
		{name: "Bad UID format",
			input: "lorem ipsum",
			want: nil,
			err:   true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := testStore.Get(tc.input)
			if err != nil {
				// If we expected the error, then test successful.
				if tc.err {
					return
				}
				t.Errorf("get error when not expected: %v", err)
			}
			if diff := reflect.DeepEqual(tc.want, got); !diff {
				t.Errorf("want: %+v\ngot: %+v", tc.want, got)
			}
			
		})
	}
}

func TestCreate(t *testing.T) {
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	host3, err := core.NewHost("host3", "192.168.51.1", "host3")
	if err != nil {
		t.Fatalf("failed to create host3: %v", err)
	}

	hosts := make([]*core.Host, 0)
	hosts = append(hosts, host3)
	
	testStore := &HostStore{
		Hosts: hosts,
		Puller: &testPuller{},
	}
	
	tests := []struct {
		name  string
		input *core.Host
		want  *core.Host
		err   bool
	}{
		{name: "Create host 1",
			input: host1,
			want: host1,
			err:   false},
		{name: "Create host 2",
			input: host2,
			want: host2,
			err:   false},
		{name: "Create host 3 - but already present",
			input: host3,
			want: host3,
			err:   true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testStore.Create(tc.input)
			if err != nil {
				// If we expected the error, then test successful.
				if tc.err {
					return
				}
				t.Errorf("get error when not expected: %v", err)
			}
			got, err := testStore.Get(tc.input.UID())
			if err != nil {
				t.Errorf("failed to get rule: %v", err)
			}
			if diff := reflect.DeepEqual(tc.want, got); !diff {
				t.Errorf("want: %+v\ngot: %+v", tc.want, got)
			}
			
		})
	}
}

func TestDelete(t *testing.T) {
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	host3, err := core.NewHost("host3", "192.168.51.1", "host3")
	if err != nil {
		t.Fatalf("failed to create host3: %v", err)
	}
	host4, err := core.NewHost("host4", "10.10.51.1", "host3")
	if err != nil {
		t.Fatalf("failed to create host4: %v", err)
	}

	hosts := make([]*core.Host, 0)
	hosts = append(hosts, host1)
	hosts = append(hosts, host2)
	hosts = append(hosts, host3)
	hosts = append(hosts, host4)
	
	testStore := &HostStore{
		Hosts: hosts,
		Puller: &testPuller{},
	}

	tests := []struct {
		name  string
		input string
		want error
		err   bool
	}{
		{name: "Delete rule 1",
			input: host1.UID(),
			want: ErrHostNotFound,
			err:   false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testStore.Delete(tc.input)
			if err != nil {
				// If we expected the error, then test successful.
				if tc.err {
					return
				}
				t.Errorf("get error when not expected: %v", err)
			}
			_, err = testStore.Get(tc.input)
			if !errors.Is(err, tc.want) {
				t.Errorf("Expected error: %v\nError received: %v", tc.want, err)
			}
		})
	}
	
}
