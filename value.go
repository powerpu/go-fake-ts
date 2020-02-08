package fake

// Value is an interface that represents fake values
type Value interface {
	// Next generates the next value internally
	Next()

	// Val retrieves the current fake value
	Val() interface{}

	// Vals retrieves the next count of fake values
	Vals(count int) []interface{}

	// JSONStats retrieves the current internal statistict for this fake value (if kept)
	JSONStats() string
}

func makeValues(fv Value, count int) []interface{} {
	out := make([]interface{}, count)

	for i := 0; i < count; i++ {
		out[i] = fv.Val()
		fv.Next()
	}

	return out
}
