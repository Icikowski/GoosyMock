package utils

// PointerTo returns a pointer for given value
func PointerTo[T any](val T) *T {
	return &val
}

// ValueOf returns a value of given pointer
func ValueOf[T any](ptr *T) T {
	var val T
	if ptr != nil {
		val = *ptr
	}
	return val
}

// ValueOrFallback returns a value of given pointer or fallback
// value if pointer is nil
func ValueOrFallback[T any](ptr *T, fallback T) T {
	if ptr != nil {
		return *ptr
	}
	return fallback
}
