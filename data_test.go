package fake

import (
	"fmt"
)

func ExampleNewData() {
	fd, _ := NewData(
		"d1",
		int64(10),

		float64(1),
		float64(1),
		float64(0),
		float64(0),
		float64(50),
		float64(100),
		false,
		false,

		int64(0),
		float64(0),
		int64(0),

		true,
		int64(1),
		float64(0.5),

		false,
		int64(5),
		int64(100),
		int64(100),
		false,
		int64(200),
		int64(20),

		false,
		int64(300),
		int64(1),
		int64(1),
		int64(1),
		int64(1),

		false)
	fmt.Printf("%v\n", fd.Vals(10))
	// Output: [-57.280531906981736 -63.071999277046245 -63.88021576640003 -63.764429475766235 -63.59492175044973 -64.63661509786012 -59.00562870317643 -55.48448571431909 -50.63657453451286 -49.45361566654604]
}

func ExampleNewData_b() {
	ft, _ := NewData(
		"d1",
		int64(10),

		float64(1),
		float64(1),
		float64(0),
		float64(0),
		float64(50),
		float64(100),
		false,
		false,

		int64(0),
		float64(0),
		int64(0),

		true,
		int64(1),
		float64(0.5),

		false,
		int64(5),
		int64(100),
		int64(100),
		false,
		int64(200),
		int64(20),

		false,
		int64(300),
		int64(1),
		int64(1),
		int64(1),
		int64(1),

		false)

	for i := 1; i <= 10; i++ {
		fmt.Printf("%v ", ft.Val())
		// Whoops! We forgot to call ft.Next() here so repeated calls get the current value.
	}
	// Output: -57.280531906981736 -57.280531906981736 -57.280531906981736 -57.280531906981736 -57.280531906981736 -57.280531906981736 -57.280531906981736 -57.280531906981736 -57.280531906981736 -57.280531906981736
}
