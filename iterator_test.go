package iter

import (
	"reflect"
	"testing"
)

func equals[T any](t *testing.T, a, b T) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%+v != %+v", a, b)
	}
}

func TestSlice(t *testing.T) {
	it := Slice([]int{1, 2, 3})
	equals(t, it.Next().Unwrap(), 1)
	equals(t, it.Next().Unwrap(), 2)
	equals(t, it.Next().Unwrap(), 3)
	equals(t, it.Next().IsNone(), true)
}

func TestChan(t *testing.T) {
	ch := make(chan int)
	arr := []int{1, 2, 3}
	go func() {
		for _, v := range arr {
			ch <- v
		}
		close(ch)
	}()
	it := Chan(ch)
	equals(t, it.Next().Unwrap(), 1)
	equals(t, it.Next().Unwrap(), 2)
	equals(t, it.Next().Unwrap(), 3)
	equals(t, it.Next().IsNone(), true)
}

func TestRepeat(t *testing.T) {
	it := Repeat(5)
	equals(t, it.Next().Unwrap(), 5)
	equals(t, it.Next().Unwrap(), 5)
	equals(t, it.Next().Unwrap(), 5)
}

func TestTake(t *testing.T) {
	it := Take(Repeat(5), 2)
	equals(t, it.Next().Unwrap(), 5)
	equals(t, it.Next().Unwrap(), 5)
	equals(t, it.Next().IsNone(), true)
}

func TestMap(t *testing.T) {
	it := Repeat(5)
	equals(t, it.Next().Unwrap(), 5)
	equals(t, it.Next().Unwrap(), 5)
	equals(t, it.Next().Unwrap(), 5)
}

func TestFunc(t *testing.T) {
	v := 0
	it := Func(func() Option[int] {
		if v < 3 {
			ret := Some(v)
			v++
			return ret
		} else {
			return None[int]()
		}
	})
	equals(t, it.Next().Unwrap(), 0)
	equals(t, it.Next().Unwrap(), 1)
	equals(t, it.Next().Unwrap(), 2)
	equals(t, it.Next().IsNone(), true)
}

func TestEmpty(t *testing.T) {
	it := Empty[int]()
	equals(t, it.Next().IsNone(), true)
}

func TestOnce(t *testing.T) {
	it := Once[int](10)
	equals(t, it.Next().Unwrap(), 10)
	equals(t, it.Next().IsNone(), true)
}

func TestToSlice(t *testing.T) {
	slice1 := ToSlice(Take(Repeat(5), 3))
	equals(t, slice1, []int{5, 5, 5})
	slice2 := ToSlice(Empty[int]())
	equals(t, slice2, []int{})
}

func TestToChan(t *testing.T) {
	slice1 := ToSlice(Chan(ToChan(Take(Repeat(5), 3))))
	equals(t, slice1, []int{5, 5, 5})
	slice2 := ToSlice(Chan(ToChan(Empty[int]())))
	equals(t, slice2, []int{})
}

func TestDrop(t *testing.T) {
	slice1 := ToSlice(Drop(Slice([]int{1, 2, 3, 4, 5}), 3))
	equals(t, slice1, []int{4, 5})
	slice2 := ToSlice(Drop(Slice([]int{1, 2, 3, 4, 5}), 0))
	equals(t, slice2, []int{1, 2, 3, 4, 5})
	slice3 := ToSlice(Drop(Slice([]int{1, 2, 3, 4, 5}), 5))
	equals(t, slice3, []int{})
}

func TestDropWhile(t *testing.T) {
	slice := ToSlice(
		DropWhile(
			Slice([]int{1, 2, 3, 4, 5}),
			func(i int) bool {
				return i < 4
			},
		),
	)
	equals(t, slice, []int{4, 5})
}

func TestTakeWhile(t *testing.T) {
	slice := ToSlice(
		TakeWhile(
			Slice([]int{1, 2, 3, 4, 5}),
			func(i int) bool {
				return i < 4
			},
		),
	)
	equals(t, slice, []int{1, 2, 3})
}

func TestForEach(t *testing.T) {
	it := Take(Repeat(1), 5)
	var ret int
	ForEach(it, func(i int) {
		ret += i
	})
	equals(t, ret, 5)
}

func TestFold(t *testing.T) {
	it := Slice([]int{1, 2, 3, 4, 5})
	ret := Fold(it, 0, func(acc, i int) int {
		return acc + i
	})
	equals(t, ret, 15)
}

func TestFuse(t *testing.T) {
	state := true
	it := Fuse(
		Func(func() Option[int] {
			ret := None[int]()
			if state {
				ret = Some(1)
			}
			state = !state
			return ret
		}),
	)
	equals(t, it.Next().Unwrap(), 1)
	equals(t, it.Next().IsNone(), true)
	equals(t, it.Next().IsNone(), true)
}

func TestChain(t *testing.T) {
	it := Chain(
		Slice([]int{1, 2}),
		Slice([]int{3, 4}),
	)
	equals(t, ToSlice(it), []int{1, 2, 3, 4})
}

func TestFind(t *testing.T) {
	it := Slice([]int{1, 2, 3, 4, 5})
	v := Find(it, func(i int) bool {
		return i > 3
	})
	equals(t, v, Some(4))
}

func TestFlatten(t *testing.T) {
	it := Slice(
		[]Iterator[int]{
			Slice([]int{1, 2}),
			Slice([]int{3, 4}),
		},
	)
	equals(t, ToSlice(Flatten(it)), []int{1, 2, 3, 4})
}

func TestRange(t *testing.T) {
	equals(t, ToSlice(Range(0, 5, 1)), []int{0, 1, 2, 3, 4})
	equals(t, ToSlice(Range(0, 0, 1)), []int{})
	equals(t, ToSlice(Range(0, -5, -1)), []int{0, -1, -2, -3, -4})
	equals(t, ToSlice(Range(5, 10, -1)), []int{})
	equals(t, ToSlice(Range(5, 10, 1)), []int{5, 6, 7, 8, 9})
}

func TestAll(t *testing.T) {
	equals(t,
		All(
			Empty[int](),
			func(n int) bool {
				return n > 0
			},
		),
		true,
	)
	equals(t,
		All(
			Slice([]int{5, 6, 7, 8, 9}),
			func(n int) bool {
				return n > 4
			},
		),
		true,
	)
	equals(t,
		All(
			Slice([]int{1, 2, 3, 4, 5}),
			func(n int) bool {
				return n <= 3
			},
		),
		false,
	)
}

func TestAny(t *testing.T) {
	equals(t,
		Any(
			Empty[int](),
			func(n int) bool {
				return n > 0
			},
		),
		false,
	)
	equals(t,
		Any(
			Slice([]int{5, 6, 7, 8, 9}),
			func(n int) bool {
				return n > 7
			},
		),
		true,
	)
	equals(t,
		Any(
			Slice([]int{1, 2, 3, 4, 5}),
			func(n int) bool {
				return n > 5
			},
		),
		false,
	)
}

func TestNth(t *testing.T) {
	equals(t, Nth(Slice([]int{1, 2, 3, 4, 5}), 0).Unwrap(), 1)
	equals(t, Nth(Slice([]int{1, 2, 3, 4, 5}), 4).Unwrap(), 5)
	equals(t, Nth(Slice([]int{1, 2, 3, 4, 5}), 10).IsNone(), true)
}

func TestEqual(t *testing.T) {
	equals(t, Equal(Slice([]int{}), Slice([]int{})), true)
	equals(t, Equal(Slice([]int{1, 2, 3}), Slice([]int{1, 2, 3})), true)
	equals(t, Equal(Slice([]int{1, 2, 3}), Slice([]int{1, 3, 2})), false)
	equals(t, Equal(Slice([]int{1, 2, 3}), Slice([]int{1, 2})), false)
}

func TestEqualBy(t *testing.T) {
	type Pair struct {
		First, Second string
	}
	equals(
		t,
		EqualBy(
			Slice([]Pair{Pair{"a", "b"}, Pair{"c", "d"}}),
			Slice([]Pair{Pair{"a", "e"}, Pair{"c", "f"}}),
			func(a, b Pair) bool {
				return a.First == b.First
			},
		),
		true,
	)
	equals(
		t,
		EqualBy(
			Slice([]Pair{Pair{"a", "b"}, Pair{"c", "d"}}),
			Slice([]Pair{Pair{"a", "e"}, Pair{"c", "f"}}),
			func(a, b Pair) bool {
				return a.Second == b.Second
			},
		),
		false,
	)
}

func TestString(t *testing.T) {
	equals(
		t,
		Equal(
			String("Hello"),
			Slice([]rune{'H', 'e', 'l', 'l', 'o'}),
		),
		true,
	)
}

func TestToString(t *testing.T) {
	equals(
		t,
		ToString(String("Hello")),
		"Hello",
	)
}
