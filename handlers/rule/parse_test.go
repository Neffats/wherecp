package rulehandler

import (
	"testing"
	
	"github.com/Neffats/wherecp/core"
)

func TestParse(t *testing.T) {
	host1, err := core.NewHost("host1", "192.168.1.1", "host1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := core.NewHost("host2", "192.168.2.1", "host2")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	host3, err := core.NewHost("host3", "192.168.3.1", "host3")
	if err != nil {
		t.Fatalf("failed to create host3: %v", err)
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

	rule := core.NewRule(src, dst, svc, true, "")
	
	tests := []struct {
		name string
		input string
		want filterFn		
		err bool
	}{
		{name: "Just_Has",
			input: "(has \"192.168.1.1\" in dst)",
			want: Has(host1, core.HasInDestination()),
			err: false},
		{name: "Object not in rule",
			input: "(has \"192.168.1.1\" in dst)",
			want: Has(host3, core.HasInDestination()),
			err: false},
		{name: "Unknown keyword",
			input: "(hasn't \"192.168.1.1\" in dst)",
			want: Has(host1, core.HasInDestination()),
			err: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotFun, err := Parse(tc.input)
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("got parse error when not expected: %v", err)
			}
			if tc.err {
				t.Fatalf("expected error, but didn't get one")
			}
			got, err := gotFun(rule)
			if err != nil {
				t.Fatalf("got error from returned filterFn: %v", err)
			}
			wanted, err := tc.want(rule)
			if err != nil {
				t.Fatalf("got error from want filterFn: %v", err)
			}
			if got != wanted {
				t.Fatalf("got: %t\nwant: %t", got, wanted)
			}
			
		})
	}
}
