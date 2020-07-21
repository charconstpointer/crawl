package set

type set interface {
	Add(value string) bool
	Contains(value string) bool
	Remove(value string) bool
}

type Set struct {
	items map[string]struct{}
}

func NewSet() Set {
	return Set{items: make(map[string]struct{})}
}

func (s *Set) Add(value string) bool {
	if s.Contains(value) {
		return false
	}
	s.items[value] = struct{}{}
	return true
}

func (s *Set) Remove(value string) {
	if s.Contains(value) {
		delete(s.items, value)
	}
}

func (s *Set) Contains(value string) bool {
	_, ok := s.items[value]
	return ok
}
