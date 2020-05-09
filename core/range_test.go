package core

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Neffats/ip"
	"github.com/google/uuid"
)

func TestNewRange(t *testing.T) {
	start := ip.Address(3232235777)
	end := ip.Address(3232236030)
	uid := uuid.New()
	testRange := &Range{
		UID:          uid.String(),
		Name:         "testRange",
		StartAddress: &start,
		EndAddress:   &end,
		Comment:      "test range object",
	}
	tests := []struct {
		name  string
		input string
		want  bool
		err   bool
	}{
		{name: "Matching range", input: "192.168.1.1-192.168.1.254", want: true, err: false},
		{name: "Invalid range", input: "192.168.1.254-192.168.1.1", want: false, err: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "-")
			got, err := NewRange("testRange", parts[0], parts[1], "test range object")
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("failed to create test range object: %v", err)
			}
			if !reflect.DeepEqual(testRange.Name, got.Name) {
				t.Fatalf("Name of created range didn't match.")
			}
			if !reflect.DeepEqual(testRange.StartAddress, got.StartAddress) {
				t.Fatalf("start address of created range didn't match.")
			}
			if !reflect.DeepEqual(testRange.EndAddress, got.EndAddress) {
				t.Fatalf("end address of created range didn't match.")
			}
			if !reflect.DeepEqual(testRange.Comment, got.Comment) {
				t.Fatalf("comment of created range didn't match.")
			}
		})

	}
}

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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "-")
			in, err := NewRange("testRange", parts[0], parts[1], "temp range object for test")
			if err != nil {
				t.Fatalf("failed to create test range object: %v", err)
			}
			got := rangeA.Match(in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestRangeContainsHost(t *testing.T) {
	rangeA, err := NewRange("rangeA", "192.168.1.1", "192.168.1.254", "test range")
	if err != nil {
		t.Fatalf("failed to create test range object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "Host contained in range", input: "192.168.1.3", want: true},
		{name: "Host outside range", input: "192.168.2.1", want: false},
		{name: "Host address same as range", input: "192.168.1.1", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in, err := NewHost("testHost", tc.input, "temp host object for test")
			if err != nil {
				t.Fatalf("failed to create test host object: %v", err)
			}
			got := rangeA.Contains(in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestRangeContainsRange(t *testing.T) {
	rangeA, err := NewRange("rangeA", "192.168.1.1", "192.168.1.254", "test range")
	if err != nil {
		t.Fatalf("failed to create test range object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "Range contained in range", input: "192.168.1.3-192.168.1.6", want: true},
		{name: "Range outside range", input: "192.168.2.1-192.168.2.5", want: false},
		{name: "Range start inside finish outside range", input: "192.168.1.5-192.168.2.3", want: false},
		{name: "Range start outside finish inside range", input: "192.168.0.5-192.168.1.33", want: false},
		{name: "Range same size as range", input: "192.168.1.1-192.168.1.254", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "-")
			in, err := NewRange("testRange", parts[0], parts[1], "temp range object for test")
			if err != nil {
				t.Fatalf("failed to create test range object: %v", err)
			}
			got := rangeA.Contains(in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestRangeContainsNetwork(t *testing.T) {
	rangeA, err := NewRange("rangeA", "192.168.1.0", "192.168.1.255", "test range")
	if err != nil {
		t.Fatalf("failed to create test range object: %v", err)
	}
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "Matching network", input: "192.168.1.0/255.255.255.0", want: true},
		{name: "Contains network", input: "192.168.1.0/255.255.255.192", want: true},
		{name: "Outside of range", input: "192.168.5.0/255.255.255.0", want: false},
		{name: "Network that contains test range", input: "192.168.0.0/255.255.240.0", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			in, err := NewNetwork("compareNet", parts[0], parts[1], "temp network for test")
			if err != nil {
				t.Fatalf("failed to create temp test network object: %v", err)
			}
			got := rangeA.Contains(in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}
