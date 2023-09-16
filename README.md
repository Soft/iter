# iter - Generic Iterators for Go 🦄

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`iter` is a generic iterator library for Go 1.18 and greater. It should feel
familiar to those familiar with [Rust's Iterator
trait](https://doc.rust-lang.org/std/iter/trait.Iterator.html).

# Iterators

```go
type Iterator[T any] interface {
        // Next yields a new value from the Iterator.
        Next() Option[T]
}
```

`Iterator[T]` represents an iterator yielding elements of type `T`.

## Creating Iterators

```go
func Slice[T any](slice []T) Iterator[T]
```

`Slice` returns an Iterator that yields elements from a slice.

```go
func Chan[T any](ch <-chan T) Iterator[T]
```

`Chan` returns an Iterator that yields elements from a channel.

```go
func String(input string) Iterator[rune]
```

`String` returns an Iterator yielding runes from the supplied string.

```go
func Range(start, stop, step int) Iterator[int]
```

`Range` returns an Iterator over a range of integers.

```go
func Func[T any](fn func() Option[T]) Iterator[T]
```

`Func` returns an Iterator from a function.

```go
func Once[T any](value T) Iterator[T]
```

`Once` returns an Iterator that returns a value exactly once.

```go
func Empty[T any]() Iterator[T]
```

`Empty` returns an empty Iterator.

```go
func Repeat[T any](value T) Iterator[T]
```

`Repeat` returns an Iterator that repeatedly returns the same value.


## Iterator Adapters

```go
func Chain[T any](first Iterator[T], second Iterator[T]) Iterator[T]
```

`Chain` returns an Iterator that concatenates two iterators.

```go
func Drop[T any](it Iterator[T], n uint) Iterator[T]
```

`Drop` returns an Iterator adapter that drops the first n items from the
underlying Iterator before yielding any values.

```go
func DropWhile[T any](it Iterator[T], pred func(T) bool) Iterator[T]
```

`DropWhile` returns an Iterator adapter that drops items from the underlying
Iterator until pred predicate function returns true.

```go
func Filter[T any](it Iterator[T], pred func(T) bool) Iterator[T]
```

`Filter` returns an Iterator adapter that yields elements from the underlying
Iterator for which pred returns true.

```go
func Flatten[T any](it Iterator[Iterator[T]]) Iterator[T]
```

`Flatten` returns an Iterator adapter that flattens nested iterators.

```go
func Fuse[T any](it Iterator[T]) Iterator[T]
```

`Fuse` returns an Iterator adapter that will keep yielding None after the
underlying Iterator has yielded None once.

```go
func Map[T, R any](it Iterator[T], fn func(T) R) Iterator[R]
```

`Map` is an Iterator adapter that transforms each value yielded by the
underlying iterator using fn.

```go
func Take[T any](it Iterator[T], n uint) Iterator[T]
```

`Take` returns an Iterator adapter that yields the n first elements from the
underlying Iterator.

```go
func TakeWhile[T any](it Iterator[T], pred func(T) bool) Iterator[T]
```

`TakeWhile` returns an Iterator adapter that yields values from the underlying
Iterator as long as pred predicate function returns true.

## Consuming Iterators

```go
func Count[T any](it Iterator[T]) uint
```

`Count` consumes an Iterator and returns the number of items it yielded.

```go
func Fold[T any, B any](it Iterator[T], init B, fn func(B, T) B) B
```

`Fold` reduces Iterator using function fn.

```go
func ForEach[T any](it Iterator[T], fn func(T))
```

`ForEach` consumes the Iterator applying fn to each yielded value.

```go
func ToSlice[T any](it Iterator[T]) []T
```

`ToSlice` consumes an Iterator creating a slice from the yielded values.

```go
func ToChan[T any](it Iterator[T]) <-chan T
```

`ToChan` consumes an Iterator writing the yielded values to a channel.

```go
func ToString(it Iterator[rune]) string
```

`ToString` consumes a rune Iterator creating a string.

```go
func Find[T any](it Iterator[T], pred func(T) bool) Option[T]
```

`Find` the first element from Iterator that satisfies pred predicate function.

```go
func All[T any](it Iterator[T], pred func(T) bool) bool
```

All tests if every element of the Iterator matches a predicate. An empty
Iterator returns true.

```go
func Any[T any](it Iterator[T], pred func(T) bool) bool
```

Any tests if any element of the Iterator matches a predicate. An empty Iterator
returns false.

```go
func Equal[T comparable](first Iterator[T], second Iterator[T]) bool
```

Determines if the elements of two Iterators are equal.

```go
func EqualBy[T any](first Iterator[T], second Iterator[T], cmp func(T, T) bool) bool
```

Determines if the elements of two Iterators are equal using function cmp to
compare elements.

```go
func Nth[T any](it Iterator[T], n uint) Option[T]
```

Nth returns nth element of the Iterator.



# Optional Values

```go
type Option[T any] struct {
        // Has unexported fields.
}
```

`Options[T]` represents an optional value of type `T`.

```go
func Some[T any](v T) Option[T]
```

`Some` returns an Option containing a value.

```go
func None[T any]() Option[T]
```

`None` returns an empty Option.

```go
func (opt Option[T]) IsSome() bool
```

`IsSome` returns true if Option contains a value.

```go
func (opt Option[T]) IsNone() bool
```

`IsNone` returns true if Option is empty.

```go
func (opt Option[T]) Unwrap() T
```

`Unwrap` extracts a value from Option. Panics if Option does not contain a
value.

```go
func (opt Option[T]) UnwrapOr(def T) T
```

`UnwrapOr` extracts a value from Option or returns a default value def if the
Option is empty.

```go
func (opt Option[T]) UnwrapOrElse(fn func() T) T
```

`UnwrapOrElse` extracts a value from Option or computes a value by calling fn if
the Option is empty.

```go
func MapOption[T any, R any](opt Option[T], fn func(T) R) Option[R]
```

`MapOption` applies a function fn to the contained value if it exists.

