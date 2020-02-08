package fake

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

// Random generates true/false values based on a predetermined percentage.
type Random struct {
	id        string
	rnd       *rand.Rand
	pctGood   float64
	keepStats bool
	Stats     *RandomStats
	v         bool
}

// RandomStats keeps track of various statistics of a Random while it's running.
type RandomStats struct {

	// The ID of the Random.
	ID string `json:"id"`

	// Cumulative count of how many times Next() was called.
	CTotal int64 `json:"cumulativeTotal"`

	// Cumulative count of how many times the value was "good".
	CGoodCount int64 `json:"cumulativeGoodCount"`

	// Cumulative count of how many times the value was "bad".
	CBadCount int64 `json:"cumulativeBadCount"`

	// Cumulative ratio of good/bad
	CRatio float64 `json:"cumulativeRatio"`

	// Slot count of how many times Next() was called. This gets reset after every JSON() call.
	Total int64 `json:"slotTotal"`

	// Slot count of how many times the value was "good". This gets reset after every JSON() call.
	GoodCount int64 `json:"slotGoodCount"`

	// Slot count of how many times the value was "bad". This gets reset after every JSON() call.
	BadCount int64 `json:"slotBadCount"`

	// Slot ratio of good/bad. This gets reset after every JSON() call.
	Ratio float64 `json:"slotGoodRatio"`
}

// Add adds a value to the running tally.
func (rs *RandomStats) Add(v interface{}) {
	rs.CTotal++
	rs.Total++

	if v.(bool) {
		rs.CGoodCount++
		rs.GoodCount++
	} else {
		rs.CBadCount++
		rs.BadCount++
	}

	rs.CRatio = float64(rs.CGoodCount) / float64(rs.CTotal)
	rs.Ratio = float64(rs.GoodCount) / float64(rs.Total)
}

// JSON returns a JSON summary of the current random statistics and resets the
// slot tally.
func (rs *RandomStats) JSON() string {
	out, _ := json.Marshal(rs)
	rs.Total = 0
	rs.GoodCount = 0
	rs.BadCount = 0
	rs.Ratio = 0
	return string(out)

}

// Next generates the next random value.
func (fr *Random) Next() {
	fr.v = fr.rnd.Float64() < fr.pctGood
	if fr.keepStats {
		fr.Stats.Add(fr.v)
	}
}

// Val returns the current random value.
func (fr *Random) Val() interface{} {
	return fr.v
}

// Vals returns the next count of values as an interface{} array.
func (fr *Random) Vals(count int) []interface{} {
	return makeValues(fr, count)
}

// JSONStats retrieves the current stats as s JSON string.
func (fr *Random) JSONStats() string {
	return fr.Stats.JSON()
}

// Good returns whether the current value is "good".
func (fr *Random) Good() bool {
	return fr.v
}

// Bad returns whether the current value is "bad".
func (fr *Random) Bad() bool {
	return !fr.v
}

// Values returns the next count of values as a bool array.
func (fr *Random) Values(count int) []bool {
	out := make([]bool, count)

	for i := 0; i < count; i++ {
		out[i] = fr.Good()
		fr.Next()
	}

	return out
}

// NewRandom creates a new Random. A random has a unique id, a random seed to
// ensure consistency when generating random numbers for the same seed, a
// percentage of required "good" samples and needs to know wheter to keep
// internal statistics.
func NewRandom(id string, seed int64, pctGood float64, keepStats bool) (*Random, error) {
	if id == "" {
		return nil, errors.New("ID for a fake random cannot be blank")
	}

	if pctGood < 0 || pctGood > 1 {
		return nil, errors.New("Percentage good for a FakeRandom with id '" + id + "' must be between 0 and 1 but was '" + fmt.Sprintf("%v", pctGood) + "'")
	}

	r := &Random{
		id:        id,
		rnd:       generateRandom(seed),
		pctGood:   pctGood,
		keepStats: keepStats,
		Stats:     &RandomStats{ID: id},
	}

	r.Next()
	return r, nil
}
