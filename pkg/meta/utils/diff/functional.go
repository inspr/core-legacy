package diff

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
func (c Changelog) ForEach(f func(c Change)) {
	for _, change := range c {
		f(change)
	}
}

// FilterDiffs filters the entries in the changelog by an arbitrary function
func (c Changelog) FilterDiffs(comp func(c string, d Difference) bool) Changelog {
	newChangelog := Changelog{}
	for _, change := range c {
		hasAdded := false
		for _, d := range change.Diff {
			if comp(change.Context, d) {
				if !hasAdded {
					newChangelog = append(newChangelog, Change{
						Context: change.Context,
						Diff:    []Difference{},
					})
					hasAdded = true
				}
				newChangelog[len(newChangelog)-1].Diff = append(newChangelog[len(newChangelog)-1].Diff, d)
			}
		}
	}
	return newChangelog
}

// DifferenceOperation is a struct that holds a filter and an operation to be applied on
// a Changelog. The filter determines which differences will be applied to and the operation defines
// what will be applied to each difference.
type DifferenceOperation struct {
	filter    func(scope string, d Difference) bool
	operation func(scope string, d Difference)
}

// NewDifferenceKindOperation creates a DifferenceOperation that filters per type.
//
// See DifferenceOperation
func NewDifferenceKindOperation(kind Kind, apply func(scope string, d Difference)) DifferenceOperation {
	return DifferenceOperation{
		filter: func(scope string, d Difference) bool {
			return d.Kind&kind > 0
		},
		operation: apply,
	}
}

// NewDifferenceOperation creates a DifferenceOperation for the given filter and apply function
//
// See DifferenceOperation
func NewDifferenceOperation(filter func(scope string, d Difference) bool, apply func(scope string, d Difference)) DifferenceOperation {
	return DifferenceOperation{
		filter, apply,
	}
}

// ForEachDiffFiltered applies each operation on the diffs contained in the changelog.
//
// The operations are applied only if the filters defined on them return true, and every filter
// is applied on each diff in the changelog.
func (c Changelog) ForEachDiffFiltered(operations ...DifferenceOperation) {
	for _, change := range c {
		for _, d := range change.Diff {
			for _, filter := range operations {
				if filter.filter(change.Context, d) {
					filter.operation(change.Context, d)
				}
			}
		}
	}
}

// ChangeOperation is a struct that holds a filter and an operation to be applied on
// a Changelog. The filter determines which changess will be applied to and the operation defines
// what will be applied to each change.
type ChangeOperation struct {
	filter func(c Change) bool
	apply  func(c Change)
}

// NewChangeOperation creates a ChangeOperation for the given filter and apply function
//
// See ChangeOperation
func NewChangeOperation(filter func(c Change) bool, apply func(c Change)) ChangeOperation {
	return ChangeOperation{
		filter, apply,
	}
}

// NewChangeKindOperation creates a ChangeOperation that filters per type.
//
// See ChangeOperation
func NewChangeKindOperation(kind Kind, apply func(c Change)) ChangeOperation {
	return ChangeOperation{
		filter: func(c Change) bool {
			return c.Kind&kind > 0
		},
		apply: apply,
	}
}

// Filter filters the differences of the Change with the return value of the given function
func (c Change) Filter(f func(d Difference) bool) (ret Change) {
	ret.Context = c.Context
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
func (c Changelog) ForEachFiltered(operations ...ChangeOperation) {
	for _, change := range c {
		for _, op := range operations {
			if op.filter(change) {
				op.apply(change)
			}
		}
	}
}

// ForEach applies the function for each change in the changelog
func (c Change) ForEach(f func(d Difference)) {
	for _, d := range c.Diff {
		f(d)
	}
}
