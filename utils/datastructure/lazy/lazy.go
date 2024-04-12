package lazy

import _interface "github.com/hopeio/cherry/utils/definition/interface"

type Lazy[T _interface.Init] struct {
	init bool
	Prop T
}

func (l *Lazy[T]) GetProp() T {
	if l.init {
		return l.Prop
	}
	l.Prop.Init()
	l.init = true
	return l.Prop
}
