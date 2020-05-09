package rulestore

import (
	"errors"
	"sync"
	
	"github.com/Neffats/wherecp/core"
)

const (
	ErrRuleNotFound = errors.New("rule not found") 
)

type RulePuller interface {
	PullRules() ([]*core.Rule, error)
}

type RuleStore struct {
	Rules []*core.Rule
	Puller RulePuller

	mux sync.RWMutex
}

func New(puller RulePuller) *RuleStore {
	return &RuleStore{
		Rules: make([]*core.Rule, 0),
		Puller: puller,
		mux: &sync.RWMutex,
	}
}

func (rs *RuleStore) Init() error {
	rules, err := rs.Puller.PullRules()
	if err != nil {
		return fmt.Errorf("failed to pull rules from source: %v", err)
	}
	rs.Rules = rules
	return nil
}
	
func (rs *RuleStore) All() []*core.Rule {
	rs.mux.RLock()
	defer rs.mux.RUnlock()
	return rs.Rules
}

func (rs *RuleStore) Create(rule *core.Rule) error {
	rs.mux.Lock()
	defer rs.mux.Unlock()
	rs.Rules = append(rs.Rules, rule)
}

func (rs *RuleStore) Get(uid string) (*core.Rule, error) {
	// TODO: Make this more efficient.
	rs.mux.RLock()
	defer rs.mux.RUnlock()
	for _, r := range rs.Rules {
		if r.UID == uid {
			return r, nil
		}
	}
	return nil, ErrRuleNotFound 
}

func (rs *RuleStore) Update(uid string, updated *core.Rule) error {
	rs.mux.RLock()
	defer rs.mux.RUnlock()
	for i, rule := range rs.Rules {
		if rule.UID == uid {
			rs.mux.Lock()
			rs.Rules[i] = updated
			rs.mux.Unlock()
			return nil
		}
	}
	return ErrRuleNotFound
}

func (rs *RuleStore) Delete(uid string) error {
	rs.mux.RLock()
	defer rs.mux.RUnlock()
	for i, rule := range rs.Rules {
		if rule.UID == uid {
			rs.mux.Lock()
			newRules := make([]*core.Rule, len(rs.Rules)-1)
			copy(newRules[:i], rs.Rules[:i])
			copy(newRules[i:], rs.Rules[i+1:])
			rs.mux.Unlock()
			return
		}
	}
	return ErrRuleNotFound
}
