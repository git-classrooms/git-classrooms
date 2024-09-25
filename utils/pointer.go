package utils

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// NewPtr returns a pointer to a new value.
func NewPtr[T any](v T) *T {
	p := new(T)
	*p = v
	return p
}
