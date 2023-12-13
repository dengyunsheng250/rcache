package policy

type entry struct {
	key   any
	value Value
}

type Value interface {
	Len() int
}
