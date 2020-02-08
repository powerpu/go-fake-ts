package fake

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
)

// Time generates timestamp values based on predetermined parameters.
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

// TimeStats keeps track of various statistics of Time while it's running.
type TimeStats struct {
	ID        string    `json:"id"`
	CTotal    int64     `json:"cumulativeTotal"`
	CEarliest time.Time `json:"cumulativeEarliestTime"`
	CLatest   time.Time `json:"cumulativeLatestTime"`
	Total     int64     `json:"slotTotal"`
	Earliest  time.Time `json:"slotEarliestTime"`
	Latest    time.Time `json:"slotLatestTime"`
}

// Add adds a value to the running tally.
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

// JSON returns a JSON summary of the current time statistics and resets the
// slot tally.
func (ts *TimeStats) JSON() string {
	out, _ := json.Marshal(ts)
	ts.Total = 0
	ts.Earliest = time.Time{}
	ts.Latest = time.Time{}
	return string(out)

}

// Next generates the next time value.
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

// Val returns the current time value as an interface{}
func (ft *Time) Val() interface{} {
	return ft.v
}

// Vals returns the next count of values as an interface{} array.
func (ft *Time) Vals(count int) []interface{} {
	return makeValues(ft, count)
}

// JSONStats retrieves the current stats as s JSON string.
func (ft *Time) JSONStats() string {
	return ft.Stats.JSON()
}

// Time returns the current time value as time.Time
func (ft *Time) Time() time.Time {
	return ft.v
}

// Times returns the next count of values as a time.Time array.
func (ft *Time) Times(count int) []time.Time {
	out := make([]time.Time, count)

	for i := 0; i < count; i++ {
		out[i] = ft.Time()
		ft.Next()
	}

	return out
}

// NewTime creates a new fake time. A time has a unique id, an initial first
// time, and increment in milliseconds, a variance in milliseconds for every
// sample, a direction of the variance (< 0 for always negative, 0 for 50/50 at
// random, > 0 for always positive) and needs to know wheter to keep internal
// statistics.
func NewTime(id string, initTs time.Time, increment int, variance int, direction int, keepStats bool) (*Time, error) {
	if id == "" {
		return nil, errors.New("ID for a fake time cannot be blank")
	}

	t := &Time{
		id:        id,
		ts:        initTs,
		increment: increment,
		variance:  variance,
		direction: direction,
		firstVal:  true,
		keepStats: keepStats,
		Stats:     &TimeStats{ID: id},
	}

	t.Next()
	return t, nil
}
