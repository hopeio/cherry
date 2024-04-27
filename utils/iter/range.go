package iter

import "golang.org/x/exp/constraints"

type Range[T constraints.Integer] struct {
	Begin T
	End   T
}

func (a *Range[T]) Get() (T, T) {
	return a.Begin, a.End
}

func (a *Range[T]) Has(index T) bool {
	return index > a.Begin && index < a.End
}

func (a *Range[T]) Iterator() Iterator[T] {
	return &rangeIter[T]{begin: a.Begin, end: a.End, step: 1, idx: 0}
}

type rangeIter[T constraints.Integer] struct {
	begin, end, step, idx T
}

// Range returns an Iterator over a range of integers.
func RangeOf[T constraints.Integer](begin, end, step T) Iterator[T] {
	return &rangeIter[T]{
		begin: begin,
		end:   end,
		step:  step,
		idx:   0,
	}
}

func (it *rangeIter[T]) Next() (T, bool) {
	v := it.begin + it.step*it.idx
	if it.step > 0 {
		if v >= it.end {
			return *new(T), false
		}
	} else {
		if v <= it.end {
			return *new(T), false
		}
	}
	it.idx++
	return v, true
}
