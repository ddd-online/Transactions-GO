// Package set provides a generic Set implementation.
package set

// Set is a generic set based on map[T]struct{}.
type Set[T comparable] map[T]struct{}

// New creates and returns a new empty set.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// NewWithElements creates a set and adds the given elements.
func NewWithElements[T comparable](elements ...T) Set[T] {
	s := New[T]()
	for _, e := range elements {
		s.Add(e)
	}
	return s
}

// Add inserts an element into the set.
func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

// Remove deletes an element from the set.
func (s Set[T]) Remove(element T) {
	delete(s, element)
}

// Has checks if the element exists in the set.
func (s Set[T]) Has(element T) bool {
	_, exists := s[element]
	return exists
}

// Size returns the number of elements in the set.
func (s Set[T]) Size() int {
	return len(s)
}

// Clear removes all elements from the set.
func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Values returns a slice of all elements in the set (order is undefined).
func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))
	for k := range s {
		values = append(values, k)
	}
	return values
}

// Union returns a new set with elements from both this set and another.
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := New[T]()
	for k := range s {
		result.Add(k)
	}
	for k := range other {
		result.Add(k)
	}
	return result
}

// Intersection returns a new set with elements common to both sets.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	result := New[T]()
	for k := range s {
		if other.Has(k) {
			result.Add(k)
		}
	}
	return result
}

// Difference returns a new set with elements in this set but not in the other (A - B).
func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := New[T]()
	for k := range s {
		if !other.Has(k) {
			result.Add(k)
		}
	}
	return result
}

// IsSubsetOf checks if this set is a subset of the other set.
func (s Set[T]) IsSubsetOf(other Set[T]) bool {
	if s.Size() > other.Size() {
		return false
	}
	for k := range s {
		if !other.Has(k) {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this set is a superset of the other set.
func (s Set[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(s)
}

// Equals checks if two sets contain exactly the same elements.
func (s Set[T]) Equals(other Set[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	return s.IsSubsetOf(other)
}
