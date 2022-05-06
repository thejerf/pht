package pht

import (
	"testing"
)

type wrapTestCase struct {
	in            string
	initialIndent int
	indent        int
	maxLine       int
	out           string
}

func TestWrapper(t *testing.T) {
	for _, test := range []wrapTestCase{
		{"", 0, 0, 80, ""},
		{"a", 0, 0, 80, "a"},
		{"a  a a", 0, 0, 80, "a a a"},
		{"abcdefg", 0, 0, 3, "abcdefg"},
		{"abcdefg h", 0, 0, 3, "abcdefg\nh"},
		{"abcdefg h", 0, 1, 3, "abcdefg\n h"},
		{"abcdefg h", 1, 0, 3, " abcdefg\nh"},
		{"abcdefg h", 1, 1, 3, " abcdefg\n h"},
		{"abcdefg h i", 1, 1, 3, " abcdefg\n h i"},
		{"abcdefg h ij", 1, 1, 3, " abcdefg\n h\n ij"},
	} {
		res := WrapString(test.in, test.initialIndent, test.indent,
			test.maxLine)
		if res != test.out {
			t.Fatalf("Expected %q, got %q",
				test.out, res)
		}
	}
}
