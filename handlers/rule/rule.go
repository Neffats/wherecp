package rulehandler

import (
	"fmt"

	"github.com/Neffats/wherecp/core"
)

// filterFn
type filterFn func(*core.Rule) (bool, error)

// And takes a number of filterFn functions, args, and returns a
// single filterFn. The returned function returns true only if all
// filterFn functions in args return true. With this you can chain
// filters together to make a more complex filter.
//
// Example:
//   filter := And(Has(hostA, InSource()), Has(HostB, InDestintation()))
//   result, err := filter(ruleA)
func And(args ...filterFn) filterFn {
	return func(r *core.Rule) (bool, error) {
		for _, arg := range args {
			res, err := arg(r)
			if err != nil {
				return false, err
			}
			if !res {
				return false, nil
			}
		}
		return true, nil
	}
}

func Or(args ...filterFn) filterFn {
	return func(r *core.Rule) (bool, error) {
		for _, arg := range args {
			res, err := arg(r)
			if err != nil {
				return false, err
			}
			if res {
				return true, nil
			}
		}
		return false, nil
	}
}

func Not(arg filterFn) filterFn {
	return func(r *core.Rule) (bool, error) {
		res, err := arg(r)
		if err != nil {
			return false, err
		}
		return !res, nil
	}
}

// Has takes an object and a comp function. The comp function
// determines which of the rule's components is to be searched
// i.e. Source, Destination or Service. The returned filterFn returns
// true if the specified component has the specified object.
//
// Example:
//   filter := Has(hostA, InDestination())
//   result, err := filter(ruleA)
func Has(obj interface{}, comp func(*core.Rule) core.Haser) filterFn {
	return func(r *core.Rule) (bool, error) {
		component := comp(r)
		has, err := component.HasObject(obj)
		if err != nil {
			return false, fmt.Errorf("failed to determine if object is in rule: %v", err)
		}

		return has, nil
	}
}

func ContainsNet(obj core.NetworkObject, comp func(*core.Rule) core.NetContainser) filterFn {
	return func(r *core.Rule) (bool, error) {
		component := comp(r)
		contains := component.Contains(obj)
		return contains, nil
	}
}

func ContainsPort(obj core.PortObject) filterFn {
	return func(r *core.Rule) (bool, error) {
		contains := r.Port().Contains(obj)
		return contains, nil
	}
}

// (and (has "192.168.1.1" (in dest)) (has "8.8.8.8" (in src)) (contains "tcp/80" (in svc)))
/*
func CreateFilter(source string) filterFn {
	tokens := LexFilter(source)
	filter := ParseTokens(tokens)
}
*/
