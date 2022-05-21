package iter

import "testing"

func TestOption(t *testing.T) {
	str := "hello world"
	var opt Option[string]
	opt = NewSome(str)

	if opt.IsNone() || !opt.IsSome() {
		t.Fatal("option should be Some not None")
	}

	val := opt.Unwrap()

	if val != str {
		t.Fatalf("unwrap gave '%s', but '%s' was expected", val, str)
	}

	switch opt.Get().(type) {
	case None[string]:
		t.Fatal("option is None, should be Some")
	}

	opt = NewNone[string]()

	switch opt.Get().(type) {
	case *Some[string]:
		t.Fatal("option is Some, should be None")
	}

	if !opt.IsNone() || opt.IsSome() {
		t.Fatal("option should be None not Some")
	}
}
