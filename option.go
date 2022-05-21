package iter

import "errors"

// SomeOrNone represents the inner value of an object
type SomeOrNone[T any] interface {
	Unwrappable() bool
}

// Some represents a value that exists. Note: we always pass Some by value,
//because we don't want to assume how it will be used. Some also never changes,
// so we never need a reference. This means we never force an allocation where
// one isn't needed, and `val` can still be a pointer to avoid copying if desired.
type Some[T any] struct{ val T }

// Unwrap unwraps the value in Some
func (s Some[T]) Unwrappable() bool {
	return true
}

// Get returns the inner value of Some
func (s Some[T]) Get() T {
	return s.val
}

// None represents a value that does not exist. Note: this is better than `nil` because
// `nil` is still a value assigned to a variable of a given type. A empty struct uses no
// memory, so really this is just a container for type data to satisfy the compiler. Whereas
// say a slice with a `nil` value will set the zero values for it's internal fields. Using
// nil to represent something that's un-initialized will take the minimum memory for that given
// type, where as `None[[]int]` will take no space because `[]int` doesn't even exist in
// memory. Even a nil pointer takes the size of a pointer still, because `nil` is still a value
// and is just assumed to mean the zero value. This is probably a meaningless optimization, but
// it's cool nonetheless. I mean why use memory to represent something that doesn't exist?
type None[T any] struct{}

// Unwrap returns false specifying that there is no value to unwrap
func (n None[T]) Unwrappable() bool {
	return false
}

var ErrUnwrapNone = errors.New("unwrapped an empty `Option`")

// Option represents an actual option
type Option[T any] interface {
	Unwrap() T
	UnwrapOr(v T) T
	UnwrapOrElse(fn func() T) T
	IsSome() bool
	IsNone() bool
	Get() SomeOrNone[T]
}

// the actual Option implementation
type option[T any] struct {
	inner SomeOrNone[T]
}

func NewSome[T any](v T) Option[T] {
	return option[T]{inner: Some[T]{v}}
}

func NewNone[T any]() Option[T] {
	return option[T]{inner: None[T]{}}
}

func (o option[T]) Unwrap() (t T) {
	exists := o.inner.Unwrappable()
	if !exists {
		panic(ErrUnwrapNone)
	}

	// we should know that this is should assert as Some at this point
	t = o.inner.(Some[T]).Get()
	return
}

func (o option[T]) UnwrapOr(v T) T {
	exists := o.inner.Unwrappable()
	if exists {
		// we should know that this is should assert as Some at this point
		return o.inner.(Some[T]).Get()
	}

	return v
}

func (o option[T]) UnwrapOrElse(fn func() T) T {
	exists := o.inner.Unwrappable()
	if exists {
		// we should know that this is should assert as Some at this point
		return o.inner.(Some[T]).Get()
	}

	return fn()
}

func (o option[T]) IsSome() bool {
	switch o.inner.(type) {
	case Some[T]:
		return true
	default:
		return false
	}
}
func (o option[T]) IsNone() bool {
	switch o.inner.(type) {
	case None[T]:
		return true
	default:
		return false
	}
}

func (o option[T]) Get() SomeOrNone[T] {
	return o.inner
}

// MapOption applies a function fn to the contained value if it exists.
func MapOption[T any, R any](opt Option[T], fn func(T) R) Option[R] {
	if !opt.IsSome() {
		return NewNone[R]()
	}
	return NewSome(fn(opt.Unwrap()))
}
