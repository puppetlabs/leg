package ref

type ID[T any] interface {
	comparable
	Less(other T) bool
}

type View interface {
	Used() bool
	Resolved() bool
	OK() bool
}

type Collection interface {
	View
	SetUsed(flag bool)
}

type Collections []Collection

var _ Collection = Collections(nil)

func (cs Collections) Used() bool {
	for _, c := range cs {
		if c != nil && c.Used() {
			return true
		}
	}
	return false
}

func (cs Collections) Resolved() bool {
	for _, c := range cs {
		if c != nil && !c.Resolved() {
			return false
		}
	}
	return true
}

func (cs Collections) OK() bool {
	for _, c := range cs {
		if c != nil && !c.OK() {
			return false
		}
	}
	return true
}

func (cs Collections) SetUsed(flag bool) {
	for _, c := range cs {
		if c != nil {
			c.SetUsed(flag)
		}
	}
}

type Reference[T ID[T]] struct {
	id       T
	used     bool
	resolved bool
	err      error
}

func (r Reference[T]) ID() T {
	return r.id
}

func (r Reference[_]) Used() bool {
	return r.used
}

func (r Reference[_]) Resolved() bool {
	return r.resolved
}

func (r Reference[_]) Error() error {
	return r.err
}

func Unused[T ID[T]](ref Reference[T]) Reference[T] {
	ref.used = false
	return ref
}

func Observed[T ID[T]](id T) Reference[T] {
	return Reference[T]{
		id:   id,
		used: true,
	}
}

func OK[T ID[T]](id T) Reference[T] {
	r := Observed(id)
	r.resolved = true
	return r
}

func Errored[T ID[T]](id T, err error) Reference[T] {
	r := Observed(id)
	r.resolved = true
	r.err = err
	return r
}

type References[T any] interface {
	Collection
	Merge(others ...T) T
}

type EmptyReferences struct{}

func (EmptyReferences) Used() bool                                    { return false }
func (EmptyReferences) Resolved() bool                                { return true }
func (EmptyReferences) OK() bool                                      { return true }
func (er EmptyReferences) Merge(_ ...EmptyReferences) EmptyReferences { return er }
func (EmptyReferences) SetUsed(flag bool)                             {}
