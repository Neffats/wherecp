package core

import (
	"testing"
)

func TestRuleHas(t *testing.T) {
	host1, err := NewHost("host1", "192.168.1.1", "Host 1")
	if err != nil {
		t.Fatalf("failed to create host1: %v", err)
	}
	host2, err := NewHost("host2", "8.8.8.8", "Google")
	if err != nil {
		t.Fatalf("failed to create host2: %v", err)
	}
	http, err := NewPort("http", 80, "tcp", "")
	if err != nil {
		t.Fatalf("failed to create http port: %v", err)
	}
	rule := NewRule(
		NewGroup("src", ""),
		NewGroup("dst", ""), 
		NewPortGroup("svc", ""),
		true,
		"Rule 1")
	rule.Source.Add(host1)
	rule.Destination.Add(host2)
	rule.Port.Add(http)

	got, err := rule.Has(host1, InSource())
	if err != nil {
		t.Fatalf("got error when not expected: %v", err)
	}
	if got != true {
		t.Fatalf("wanted: %v, got: %v", true, got)
	}
}
