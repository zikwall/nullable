package nullable

import "encoding/json"

// Nullable is a generic type that can hold a value or represent a null state.
type Nullable[T any] struct {
	value   T
	notNull bool
}

// New creates a new Nullable instance with the given value, marked as not null.
func New[T any](v T) Nullable[T] {
	return Nullable[T]{
		value:   v,
		notNull: true,
	}
}

// ConvertRef converts a reference of type S to of Nullable of type T using the provided conversion function.
// If the reference is nil, it returns a null Nullable.
func ConvertRef[S, T any](s *S, convert func(*S) T) Nullable[T] {
	if s == nil {
		return Nullable[T]{}
	}

	return Nullable[T]{
		value:   convert(s),
		notNull: true,
	}
}

// FromRef creates a Nullable from a pointer to a value of type T.
// If the pointer is nil, it returns a null Nullable.
func FromRef[T any](v *T) Nullable[T] {
	if v == nil {
		return Nullable[T]{}
	}

	return Nullable[T]{
		value:   *v,
		notNull: true,
	}
}

// Value returns the value stored in the Nullable. If the Nullable is null, it returns the zero value of T.
func (n Nullable[T]) Value() T {
	if n.IsNull() {
		var v T
		return v
	}

	return n.value
}

// IsNull returns true if the Nullable is in a null state.
func (n Nullable[T]) IsNull() bool {
	return !n.notNull
}

// NotNull returns true if the Nullable contains a non-null value.
func (n Nullable[T]) NotNull() bool {
	return n.notNull
}

// Ref returns a pointer to the value contained in the Nullable.
// If the Nullable is null, it returns nil.
func (n Nullable[T]) Ref() *T {
	if n.IsNull() {
		return nil
	}

	return &n.value
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the contained value if not null, otherwise it serializes as null.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.notNull {
		return json.Marshal(nil)
	}
	return json.Marshal(n.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes JSON data into the Nullable struct, setting it to null if the data is JSON null.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.notNull = false
		return nil
	}

	if len(data) == 0 {
		n.notNull = false
		return nil
	}

	err := json.Unmarshal(data, &n.value)
	if err != nil {
		return err
	}

	n.notNull = true
	return nil
}
