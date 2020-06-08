package rulestore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Neffats/wherecp/core"
)

type testPuller struct {}

func (tp *testPuller) PullRules() ([]*core.Rule, error) {
	return make([]*core.Rule, 0), nil
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
	http, err := core.NewPort("http", 80, "tcp", "http port")
	if err != nil {
		t.Fatalf("failed to create http service: %v", err)
	}
	src := core.NewGroup("src", "src")
	err = src.Add(host1)
	if err != nil {
		t.Fatalf("failed to add host1 to src: %v", err)
	}

	dst := core.NewGroup("dst", "dst")
	err = dst.Add(host2)
	if err != nil {
		t.Fatalf("failed to add host2 to dst: %v", err)
	}

	svc := core.NewPortGroup("svc", "svc")
	err = svc.Add(http)
	if err != nil {
		t.Fatalf("failed to add http to svc: %v", err)
	}

	rule1 := core.NewRule(1, src, dst, svc, true, "")
	rule2 := core.NewRule(2, dst, src, svc, true, "")
	rule3 := core.NewRule(3, src, src, svc, false, "")

	rules := make([]*core.Rule, 0)
	rules = append(rules, rule1)
	rules = append(rules, rule2)
	rules = append(rules, rule3)
	
	testStore := &RuleStore{
		Rules: rules,
		Puller: &testPuller{},
	}

	got := testStore.All()
	if len(got) != len(rules) {
		t.Errorf("length of expected rules doesn't match rules we got, wanted length: %d, length gotten: %d\n", len(rules), len(got))
		return
	}
	for i := 0; i < len(rules); i++ {
		if diff := reflect.DeepEqual(got[i], rules[i]); !diff {
			t.Errorf("expected: %+v\ngot:%+v", got[i], rules[i])
		}
	}
}

func TestGet(t *testing.T) {
	// Setup the test data.
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	http, err := core.NewPort("http", 80, "tcp", "http port")
	if err != nil {
		t.Fatalf("failed to create http service: %v", err)
	}
	src := core.NewGroup("src", "src")
	err = src.Add(host1)
	if err != nil {
		t.Fatalf("failed to add host1 to src: %v", err)
	}

	dst := core.NewGroup("dst", "dst")
	err = dst.Add(host2)
	if err != nil {
		t.Fatalf("failed to add host2 to dst: %v", err)
	}

	svc := core.NewPortGroup("svc", "svc")
	err = svc.Add(http)
	if err != nil {
		t.Fatalf("failed to add http to svc: %v", err)
	}

	rule1 := core.NewRule(1, src, dst, svc, true, "")
	rule2 := core.NewRule(2, dst, src, svc, true, "")
	rule3 := core.NewRule(3, src, src, svc, false, "")

	rules := make([]*core.Rule, 0)
	rules = append(rules, rule1)
	rules = append(rules, rule2)
	//rules = append(rules, rule3)
	
	testStore := &RuleStore{
		Rules: rules,
		Puller: &testPuller{},
	}

	tests := []struct {
		name  string
		input string
		want  *core.Rule
		err   bool
	}{
		{name: "Get rule 1",
			input: rule1.UID(),
			want: rule1,
			err:   false},
		{name: "Get rule 2",
			input: rule2.UID(),
			want: rule2,
			err:   false},
		{name: "Non-exitent rule",
			input: rule3.UID(),
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

func TestInsert(t *testing.T) {
	// Setup the test data.
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	http, err := core.NewPort("http", 80, "tcp", "http port")
	if err != nil {
		t.Fatalf("failed to create http service: %v", err)
	}
	src := core.NewGroup("src", "src")
	err = src.Add(host1)
	if err != nil {
		t.Fatalf("failed to add host1 to src: %v", err)
	}

	dst := core.NewGroup("dst", "dst")
	err = dst.Add(host2)
	if err != nil {
		t.Fatalf("failed to add host2 to dst: %v", err)
	}

	svc := core.NewPortGroup("svc", "svc")
	err = svc.Add(http)
	if err != nil {
		t.Fatalf("failed to add http to svc: %v", err)
	}

	rule1 := core.NewRule(1, src, dst, svc, true, "")
	rule2 := core.NewRule(2, dst, src, svc, true, "")
	rule3 := core.NewRule(3, src, src, svc, false, "")

	rules := make([]*core.Rule, 0)
	rules = append(rules, rule3)
	
	testStore := &RuleStore{
		Rules: rules,
		Puller: &testPuller{},
	}
	
	tests := []struct {
		name  string
		input *core.Rule
		want  *core.Rule
		err   bool
	}{
		{name: "Create rule 1",
			input: rule1,
			want: rule1,
			err:   false},
		{name: "Create rule 2",
			input: rule2,
			want: rule2,
			err:   false},
		{name: "Create rule 3 - but already present",
			input: rule3,
			want: rule3,
			err:   true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testStore.Insert(tc.input)
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
				t.Errorf("want: %+v\ngot:%+v", tc.want, got)
			}
			
		})
	}
}

func TestDelete(t *testing.T) {
	// Setup the test data.
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	http, err := core.NewPort("http", 80, "tcp", "http port")
	if err != nil {
		t.Fatalf("failed to create http service: %v", err)
	}
	src := core.NewGroup("src", "src")
	err = src.Add(host1)
	if err != nil {
		t.Fatalf("failed to add host1 to src: %v", err)
	}

	dst := core.NewGroup("dst", "dst")
	err = dst.Add(host2)
	if err != nil {
		t.Fatalf("failed to add host2 to dst: %v", err)
	}

	svc := core.NewPortGroup("svc", "svc")
	err = svc.Add(http)
	if err != nil {
		t.Fatalf("failed to add http to svc: %v", err)
	}

	rule1 := core.NewRule(1, src, dst, svc, true, "")
	rule2 := core.NewRule(2, dst, src, svc, true, "")
	rule3 := core.NewRule(3, src, src, svc, false, "")

	rules := make([]*core.Rule, 0)
	rules = append(rules, rule1)
	rules = append(rules, rule2)
	rules = append(rules, rule3)
	
	testStore := &RuleStore{
		Rules: rules,
		Puller: &testPuller{},
	}

	tests := []struct {
		name  string
		input string
		want error
		err   bool
	}{
		{name: "Delete rule 1",
			input: rule1.UID(),
			want: ErrRuleNotFound,
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
