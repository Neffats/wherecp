package rulestore

import (
	"errors"
	"fmt"
	"sync"
	
	"github.com/Neffats/wherecp/core"
)

var (
	ErrRuleNotFound = errors.New("rule not found") 
)

// RulePuller is the interface that any 
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

// All return all of the rules in the store. 
func (rs *RuleStore) All() []*core.Rule {
	rs.mux.RLock()
	defer rs.mux.RUnlock()
	// Create a new copy of rules to stop accidental modification.
	r := make([]*core.Rule, len(rs.Rules))
	copy(r, rs.Rules)
	return r
}

func (rs *RuleStore) Create(rule *core.Rule) error {
	// Check whether the rule is already in the store.
	existing, err := rs.Get(rule.UID)
	if !errors.Is(err, ErrRuleNotFound) {
		return fmt.Errorf("failed to determine whether rule is already present: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("rule already in store")
	}
	rs.mux.Lock()
	defer rs.mux.Unlock()
	rs.Rules = append(rs.Rules, rule)
	return nil
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
	//rs.mux.RLock()
	//defer rs.mux.RUnlock()
	for i, rule := range rs.Rules {
		if rule.UID == uid {
			rs.mux.Lock()
			newRules := make([]*core.Rule, len(rs.Rules)-1)
			copy(newRules[:i], rs.Rules[:i])
			copy(newRules[i:], rs.Rules[i+1:])
			rs.Rules = newRules
			rs.mux.Unlock()
			return nil
		}
	}
	return ErrRuleNotFound
}
