package ads

type Set[T comparable] map[T]int

func (s Set[T]) Add(k T) {
	if _, ok := s[k]; !ok {
		s[k] = len(s)
	}
}

func (s Set[T]) Range(fn func(k T, idx int) bool) {
	for _, k := range s.Slice() {
		if !fn(k, s[k]) {
			break
		}
	}
}

// Slice
func (s Set[T]) Slice() []T {
	var keys = make([]T, len(s))
	for k, idx := range s {
		keys[idx] = k
	}
	return keys
}
