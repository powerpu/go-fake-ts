package fake

type Value interface {
	Next()
	Val() interface{}
	Vals(count int) []interface{}
	JsonStats() string
}

func makeValues(fv Value, count int) []interface{} {
	out := make([]interface{}, count)

	for i := 0; i < count; i++ {
		fv.Next()
		out[i] = fv.Val()
	}

	return out
}
