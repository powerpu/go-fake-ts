package fake

import (
	"fmt"
	"time"
)

func ExampleNewTime() {
	t := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	ft, _ := NewTime("fakeTime1", t, 5000, 1000, -1, true)

	for i := 1; i <= 10; i++ {
		fmt.Printf("%v ", ft.Time().Unix())
		ft.Next()
	}
	// Output: 1258490098 1258490102 1258490107 1258490113 1258490118 1258490122 1258490128 1258490133 1258490138 1258490143
}

func ExampleNewTime_b() {
	t := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	ft, _ := NewTime("fakeTime1", t, 5000, 1000, -1, true)

	for i := 1; i <= 10; i++ {
		fmt.Printf("%v ", ft.Time().Unix())
		// Whoops! We forgot to call t.Next() here so repeated calls get the current value.
	}
	// Output: 1258490098 1258490098 1258490098 1258490098 1258490098 1258490098 1258490098 1258490098 1258490098 1258490098
}
