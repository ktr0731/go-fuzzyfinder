package scoring

import "testing"

func TestCalculate(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		s1, s2    string
		willPanic bool
	}{
		"must not panic":  {s1: "foo", s2: "foo"},
		"must not panic2": {s1: "", s2: ""},
		"must panic":      {s1: "foo", s2: "foobar", willPanic: true},
	}

	for _, c := range cases {
		if c.willPanic {
			defer func() {
				if err := recover(); err == nil {
					t.Error("Calculate must panic")
				}
			}()
		}
		Calculate(c.s1, c.s2)
	}
}

func Test_max(t *testing.T) {
	t.Parallel()

	if n := max(); n != 0 {
		t.Errorf("max must return 0 if no args, but got %d", n)
	}

	if n := max(0, -1, 10, 3); n != 10 {
		t.Errorf("max must return the maximun number 10, but got %d", n)
	}
}
