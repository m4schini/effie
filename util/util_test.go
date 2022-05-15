package util

import "testing"

type Case struct {
	Number   int
	Expected string
}

func TestGetNumberPostfix(t *testing.T) {
	cases := []Case{
		{
			Number:   0,
			Expected: "th",
		},
		{
			Number:   -1,
			Expected: "st",
		},
		{
			Number:   1,
			Expected: "st",
		},
		{
			Number:   2,
			Expected: "nd",
		},
		{
			Number:   3,
			Expected: "rd",
		},
	}

	for _, c := range cases {
		if GetNumberPostfix(c.Number) != c.Expected {
			t.Fail()
		}
	}
}
