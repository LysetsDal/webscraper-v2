package main

import (
	"strconv"
	"testing"
)

// If the number is divisible by 3, write "Foo" otherwise, the number
func Fooer(input int) string {
	isfoo := (input % 3) == 0
	if isfoo {
		return "Foo"
	}
	return strconv.Itoa(input)
}

func TestFooerTableDriven(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name  string
		input int
		want  string
	}{
		// the table itself
		{"9 should be Foo", 9, "Foo"},
		{"3 should be Foo", 3, "Foo"},
		{"1 is not Foo", 1, "1"},
		{"0 should be Foo", 0, "Foo"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := Fooer(tt.input)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
	/* NOTE:
	 * The execution loop calls t.Run(), which defines a subtest. 
	 * As a result, each row of the table defines a subtest named 
	 * [NameOfTheFuction]/[NameOfTheSubTest]
	 */
}
