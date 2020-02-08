package fake

import (
	"fmt"
)

func ExampleNewRandom() {
	fr, _ := NewRandom("fakeRandom1", 4, 0.5, true)
	fmt.Printf("%v\n", fr.Vals(10))
	// Output: [true true true false true true false false false false]
}

func ExampleNewRandom_b() {
	fr, _ := NewRandom("fakeRandom2", 4, 0.5, true)

	for i := 1; i <= 10; i++ {
		fmt.Printf("%v ", fr.Good())
		// Whoops! We forgot to call fr.Next() here so repeated calls get the current value.
	}
	// Output: true true true true true true true true true true
}
