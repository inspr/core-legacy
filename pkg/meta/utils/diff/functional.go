package diff

import (
	"github.com/inspr/inspr/pkg/ierrors"
)

// FilterDiffsByKind filters the changelog by changes of a spefific kind or combination of kinds.
// To use a combination, just bitwise or the different kinds, so, if you want the changes on channels and
// apps, use
//
// c.FilterDiffsByKind(AppKind|ChannelKind)
func (c Changelog) FilterDiffsByKind(kind Kind) Changelog {
	return c.FilterDiffs(func(_ string, d Difference) bool {
		return d.Kind&kind > 0
	})
}

// ForEach applies a function on each change of the changelog
func (c Changelog) ForEach(f func(c Change) error) error {
	errors := ierrors.MultiError{Errors: []error{}}
	for _, change := range c {
		if change.Scope == "*" {
			change.Scope = ""
		}
		errors.Add(f(change))
	}
	return &errors
}

// FilterDiffs filters the entries in the changelog by an arbitrary function
func (c Changelog) FilterDiffs(comp DifferenceFilter) Changelog {
	newChangelog := Changelog{}
	for _, change := range c {
		if change.Scope == "*" {
			change.Scope = ""
		}
		hasAdded := false
		for _, d := range change.Diff {
			if comp(change.Scope, d) {
				if !hasAdded {
					newChangelog = append(newChangelog, Change{
						Scope: change.Scope,
						Diff:  []Difference{},
					})
					hasAdded = true
				}
				newChangelog[len(newChangelog)-1].Diff = append(newChangelog[len(newChangelog)-1].Diff, d)
			}
		}
	}
	return newChangelog
}

// DifferenceFilter filters differnces
type DifferenceFilter func(scope string, d Difference) bool

// DifferenceOperation applies a function in a difference and returns its error
type DifferenceOperation func(scope string, d Difference) error

// DifferenceReaction is a struct that holds a filter and an operation to be applied on
// a Changelog. The filter determines which differences will be applied to and the operation defines
// what will be applied to each difference.
type DifferenceReaction struct {
	filter    DifferenceFilter
	operation DifferenceOperation
}

// NewDifferenceKindReaction creates a DifferenceOperation that filters per type.
//
// See DifferenceOperation
func NewDifferenceKindReaction(kind Kind, apply DifferenceOperation) DifferenceReaction {
	return DifferenceReaction{
		filter: func(scope string, d Difference) bool {
			return d.Kind&kind > 0
		},
		operation: apply,
	}
}

// NewDifferenceReaction creates a DifferenceOperation for the given filter and apply function
//
// See DifferenceOperation
func NewDifferenceReaction(filter DifferenceFilter, apply DifferenceOperation) DifferenceReaction {
	return DifferenceReaction{
		filter, apply,
	}
}

// ForEachDiffFiltered applies each operation on the diffs contained in the changelog.
//
// The operations are applied only if the filters defined on them return true, and every filter
// is applied on each diff in the changelog.
//
// Errors are concatenated
func (c Changelog) ForEachDiffFiltered(operations ...DifferenceReaction) error {
	errors := ierrors.MultiError{
		Errors: []error{},
	}
	for _, change := range c {
		if change.Scope == "*" {
			change.Scope = ""
		}
		for _, d := range change.Diff {
			for _, filter := range operations {
				if filter.filter(change.Scope, d) {
					errors.Add(filter.operation(change.Scope, d))
				}
			}
		}
	}
	if errors.Empty() {
		return nil
	}
	return &errors
}

// ChangeFilter filters changes
type ChangeFilter func(c Change) bool

// ChangeOperation applies an operation in a change and returns an error
type ChangeOperation func(c Change) error

// ChangeReaction is a struct that holds a filter and an operation to be applied on
// a Changelog. The filter determines which changess will be applied to and the operation defines
// what will be applied to each change.
type ChangeReaction struct {
	filter ChangeFilter
	apply  ChangeOperation
}

// NewChangeReaction creates a ChangeOperation for the given filter and apply function
//
// See ChangeOperation
func NewChangeReaction(filter ChangeFilter, apply ChangeOperation) ChangeReaction {
	return ChangeReaction{
		filter, apply,
	}
}

// NewChangeKindReaction creates a ChangeOperation that filters per type.
//
// See ChangeOperation
func NewChangeKindReaction(kind Kind, apply func(c Change) error) ChangeReaction {
	return ChangeReaction{
		filter: func(c Change) bool {
			return c.Kind&kind > 0
		},
		apply: apply,
	}
}

// Filter filters the differences of the Change with the return value of the given function
func (c Change) Filter(f func(d Difference) bool) (ret Change) {
	ret.Scope = c.Scope
	ret.Kind = c.Kind
	ret.Diff = []Difference{}

	for _, d := range c.Diff {
		if f(d) {
			ret.Diff = append(ret.Diff, d)
		}
	}
	return
}

// FilterKind filters the diffs of the change by its kind.
func (c Change) FilterKind(kind Kind) Change {
	return c.Filter(func(d Difference) bool {
		return d.Kind&kind > 0
	})
}

// ForEachFiltered applies each operation on the changes contained in the changelog.
//
// The operations are applied only if the filters defined on them return true, and every filter
// is applied on each change in the changelog.
//
// Errors are concatenated.
func (c Changelog) ForEachFiltered(operations ...ChangeReaction) error {
	errors := ierrors.MultiError{
		Errors: []error{},
	}
	for _, change := range c {
		if change.Scope == "*" {
			change.Scope = ""
		}
		for _, op := range operations {
			if op.filter(change) {
				errors.Add(op.apply(change))
			}
		}
	}
	if errors.Empty() {
		return nil
	}
	return &errors
}

// ForEach applies the function for each change in the changelog
//
// Errors are concatenated
func (c Change) ForEach(f DifferenceOperation) error {
	errors := ierrors.MultiError{
		Errors: []error{},
	}
	for _, d := range c.Diff {
		errors.Add(f(c.Scope, d))
	}
	if errors.Empty() {
		return nil
	}
	return &errors
}
