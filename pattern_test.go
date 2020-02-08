package fake

import (
	"fmt"
)

func ExampleNewPattern() {
	fp, _ := NewPattern("fakePattern1", 2, 1, true)
	fmt.Printf("%v\n", fp.Vals(10))
	// Output: [true true false true true false true true false true]
}

func ExampleNewPattern_b() {
	fp, _ := NewPattern("fakePattern2", 3, 1, true)

	for i := 1; i <= 10; i++ {
		fmt.Printf("%v ", fp.Good())
		// Whoops! We forgot to call fp.Next() here so repeated calls get the current value.
	}
	// Output: true true true true true true true true true true
}
