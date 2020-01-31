package fake

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
)

type Time struct {
	id           string
	increment    int
	variance     int
	direction    int
	ts           time.Time
	varianceTime time.Time
	firstVal     bool
	keepStats    bool
	Stats        *TimeStats
	v            time.Time
}

type TimeStats struct {
	Id        string    `json:"id"`
	CTotal    int64     `json:"cumulativeTotal"`
	CEarliest time.Time `json:"cumulativeEarliestTime"`
	CLatest   time.Time `json:"cumulativeLatestTime"`
	Total     int64     `json:"slotTotal"`
	Earliest  time.Time `json:"slotEarliestTime"`
	Latest    time.Time `json:"slotLatestTime"`
}

func (ts *TimeStats) Add(v interface{}) {
	ts.CTotal++
	ts.Total++

	ts.CLatest = v.(time.Time)
	ts.Latest = v.(time.Time)

	if ts.CEarliest.IsZero() {
		ts.CEarliest = v.(time.Time)
	}

	if ts.Earliest.IsZero() {
		ts.Earliest = v.(time.Time)
	}
}

func (ts *TimeStats) Json() string {
	out, _ := json.Marshal(ts)
	ts.Total = 0
	ts.Earliest = time.Time{}
	ts.Latest = time.Time{}
	return string(out)

}

func (ft *Time) Next() {
	a := rand.Float64()

	// Ensure first time doesn't have any variance to respect the start time parameter
	if ft.firstVal {
		ft.firstVal = false
		ft.v = ft.ts

		if ft.keepStats {
			ft.Stats.Add(ft.v)
		}

		return
	}

	ft.ts = ft.ts.Add(time.Duration(ft.increment) * time.Millisecond)
	tmp := (float64(ft.variance) * a) - float64(int64(float64(ft.variance)*a))
	tmp2 := float64(-1)

	if ft.direction < 0 {
		tmp2 = float64(-1)
	} else if ft.direction > 0 {
		tmp2 = float64(1)
	} else if tmp > 0.5 {
		tmp2 = float64(1)
	}

	c := int64(Round(float64(ft.variance)*a, 0.0000000005) * tmp2)
	ft.v = ft.ts.Add(time.Duration(c) * time.Millisecond)

	if ft.keepStats {
		ft.Stats.Add(ft.v)
	}
}

func (ft *Time) Val() interface{} {
	return ft.v
}

func (ft *Time) Vals(count int) []interface{} {
	return makeValues(ft, count)
}

func (ft *Time) JsonStats() string {
	return ft.Stats.Json()
}

func (ft *Time) Time() time.Time {
	return ft.v
}

func (ft *Time) Times(count int) []time.Time {
	out := make([]time.Time, count)

	for i := 0; i < count; i++ {
		ft.Next()
		out[i] = ft.Time()
	}

	return out
}

func NewTime(id string, initTs time.Time, increment int, variance int, direction int, keepStats bool) (*Time, error) {
	if id == "" {
		return nil, errors.New("ID for a fake time cannot be blank")
	}

	return &Time{
		id:        id,
		ts:        initTs,
		increment: increment,
		variance:  variance,
		direction: direction,
		firstVal:  true,
		keepStats: keepStats,
		Stats:     &TimeStats{Id: id},
	}, nil
}
