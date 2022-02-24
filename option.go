package iter

// Options[T] represents an optional value of type T.
type Option[T any] struct{ value *T }

// Some returns an Option containing a value.
func Some[T any](v T) Option[T] {
	return Option[T]{value: &v}
}

// None returns an empty Option.
func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

// IsNone returns true if Option is empty.
func (opt Option[T]) IsNone() bool {
	return opt.value == nil
}

// IsSome returns true if Option contains a value.
func (opt Option[T]) IsSome() bool {
	return !opt.IsNone()
}

// Unwrap extracts a value from Option. Panics if Option does not contain a
// value.
func (opt Option[T]) Unwrap() T {
	if opt.IsNone() {
		panic("Attempted to unwrap an empty Option.")
	}
	return *opt.value
}

// UnwrapOr extracts a value from Option or returns a default value def if the
// Option is empty.
func (opt Option[T]) UnwrapOr(def T) T {
	if opt.IsSome() {
		return opt.Unwrap()
	}
	return def
}

// UnwrapOrElse extracts a value from Option or computes a value by calling fn
// if the Option is empty.
func (opt Option[T]) UnwrapOrElse(fn func() T) T {
	if opt.IsSome() {
		return opt.Unwrap()
	}
	return fn()
}

// MapOption applies a function fn to the contained value if it exists.
func MapOption[T any, R any](opt Option[T], fn func(T) R) Option[R] {
	if !opt.IsSome() {
		return None[R]()
	}
	return Some(fn(opt.Unwrap()))
}
