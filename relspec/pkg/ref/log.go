package ref

import "golang.org/x/exp/slices"

type Log[T ID[T]] struct {
	m map[T]Reference[T]
}

func (l *Log[T]) Set(r Reference[T]) {
	pr := l.m[r.id]
	if !r.used {
		r.used = pr.used
	}
	if !r.resolved {
		r.resolved = pr.resolved
	}
	if r.err == nil {
		r.err = pr.err
	}
	l.m[r.id] = r
}

func (l *Log[T]) ForEach(fn func(Reference[T])) {
	if l == nil {
		return
	}

	for _, ref := range l.m {
		fn(ref)
	}
}

func (l *Log[T]) Merge(others ...*Log[T]) *Log[T] {
	if len(others) == 0 {
		return l
	} else if l == nil {
		return NewLog[T]().Merge(others...)
	}

	for _, log := range others {
		log.ForEach(func(ref Reference[T]) { l.Set(ref) })
	}

	return l
}

func (l *Log[T]) SetUsed(flag bool) {
	l.ForEach(func(ref Reference[T]) {
		ref.used = flag
		l.m[ref.id] = ref
	})
}

func (l *Log[T]) Filter(fn func(Reference[T]) bool) []Reference[T] {
	var s []Reference[T]
	l.ForEach(func(ref Reference[T]) {
		if fn(ref) {
			s = append(s, ref)
		}
	})
	slices.SortFunc(s, func(a, b Reference[T]) bool { return a.ID().Less(b.ID()) })
	return s
}

func (l *Log[T]) AllReferences() []Reference[T] {
	return l.Filter(func(_ Reference[T]) bool { return true })
}

func (l *Log[T]) UsedReferences() []Reference[T] {
	return l.Filter(func(ref Reference[T]) bool { return ref.Used() })
}

func (l *Log[T]) Used() bool {
	return len(l.UnresolvedReferences()) > 0
}

func (l *Log[T]) ResolvedReferences() []Reference[T] {
	return l.Filter(func(ref Reference[T]) bool {
		return ref.Used() && ref.Resolved()
	})
}

func (l *Log[T]) UnresolvedReferences() []Reference[T] {
	return l.Filter(func(ref Reference[T]) bool {
		return ref.Used() && !ref.Resolved()
	})
}

func (l *Log[T]) Resolved() bool {
	return len(l.UnresolvedReferences()) == 0
}

func (l *Log[T]) OKReferences() []Reference[T] {
	return l.Filter(func(ref Reference[T]) bool {
		return ref.Used() && ref.Resolved() && ref.Error() == nil
	})
}

func (l *Log[T]) ErroredReferences() []Reference[T] {
	return l.Filter(func(ref Reference[T]) bool {
		return ref.Used() && ref.Resolved() && ref.Error() != nil
	})
}

func (l *Log[T]) OK() bool {
	problems := l.Filter(func(ref Reference[T]) bool {
		return ref.Used() && (!ref.Resolved() || ref.Error() != nil)
	})
	return len(problems) == 0
}

func NewLog[T ID[T]]() *Log[T] {
	return &Log[T]{m: make(map[T]Reference[T])}
}

func CopyLog[T ID[T]](from *Log[T]) *Log[T] {
	return NewLog[T]().Merge(from)
}

func InitialLog[T ID[T]](refs ...Reference[T]) *Log[T] {
	log := NewLog[T]()
	for _, ref := range refs {
		log.Set(ref)
	}
	return log
}
