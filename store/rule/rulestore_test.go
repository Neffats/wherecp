package rulestore

import (
	"testing"

	"github.com/Neffats/wherecp/core"
)

type testPuller struct {}

func (tp *testPuller) PullRules() ([]*core.Rule, error) {
	
}


