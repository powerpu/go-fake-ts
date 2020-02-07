package fake

import (
	"encoding/json"
	"errors"
	"math"
)

// Pattern generates true/false values based on a predetermined pattern.
type Pattern struct {
	id           string
	i            int64
	patternGood  int32
	patternBad   int32
	patternRatio float64
	keepStats    bool
	Stats        *PatternStats
	v            bool
}

// PatternStats keeps track of various statistics of a Pattern while it's running.
type PatternStats struct {

    // The ID of the Pattern.
	Id         string  `json:"id"`

    // Cumulative count of how many times Next() was called.
	CTotal     int64   `json:"cumulativeTotal"`

    // Cumulative count of how many times the value was "good".
	CGoodCount int64   `json:"cumulativeGoodCount"`

    // Cumulative count of how many times the value was "bad".
	CBadCount  int64   `json:"cumulativeBadCount"`

    // Cumulative ratio of good/bad
	CRatio     float64 `json:"cumulativeRatio"`

    // Slot count of how many times Next() was called. This gets reset after every JSON() call.
	Total      int64   `json:"slotTotal"`

    // Slot count of how many times the value was "good". This gets reset after every JSON() call.
	GoodCount  int64   `json:"slotGoodCount"`

    // Slot count of how many times the value was "bad". This gets reset after every JSON() call.
	BadCount   int64   `json:"slotBadCount"`

    // Slot ratio of good/bad. This gets reset after every JSON() call.
	Ratio      float64 `json:"slotGoodRatio"`
}

// Add adds a value to the running tally.
func (ps *PatternStats) Add(v bool) {
	ps.CTotal++
	ps.Total++

	if v {
		ps.CGoodCount++
		ps.GoodCount++
	} else {
		ps.CBadCount++
		ps.BadCount++
	}

	ps.CRatio = float64(ps.CGoodCount) / float64(ps.CTotal)
	ps.Ratio = float64(ps.GoodCount) / float64(ps.Total)
}

// Returns a JSON summary of the current pattern statistics and resets the slot
// tally. 
func (ps *PatternStats) Json() string {
	out, _ := json.Marshal(ps)
	ps.Total = 0
	ps.GoodCount = 0
	ps.BadCount = 0
	ps.Ratio = 0
	return string(out)

}

// Generates the next pattern value.
func (fp *Pattern) Next() {
	fp.i = fp.i + 1
	f := math.Floor((float64(fp.i)/float64(fp.patternGood+fp.patternBad))*100000000) / 100000000
	fp.v = (fp.patternBad == 0 || ((f-math.Trunc(f)) <= fp.patternRatio && (f-math.Trunc(f)) != 0))
	if fp.keepStats {
		fp.Stats.Add(fp.v)
	}
}

// Returns the current pattern value.
func (fp *Pattern) Val() interface{} {
	return fp.v
}

// Returns the next count of values as an interface{} array.
func (fp *Pattern) Vals(count int) []interface{} {
	return makeValues(fp, count)
}

// Retrieves the current stats as s JSON string.
func (fp *Pattern) JsonStats() string {
	return fp.Stats.Json()
}

// Returns whether the current value is "good".
func (fp *Pattern) Good() bool {
	return fp.v
}

// Returns whether the current value is "bad".
func (fp *Pattern) Bad() bool {
	return !fp.v
}

// Returns the next count of values as a bool array.
func (fp *Pattern) Values(count int) []bool {
	out := make([]bool, count)

	for i := 0; i < count; i++ {
		fp.Next()
		out[i] = fp.Good()
	}

	return out
}

// Creates a new Pattern. A pattern has a unique id, number of required "good"
// samples followed by a number of required "bad" samples and needs to know
// wheter to keep internal statistics.
func NewPattern(id string, patternGood int32, patternBad int32, keepStats bool) (*Pattern, error) {
	if id == "" {
		return nil, errors.New("ID for a fake pattern cannot be blank")
	}

	if patternGood < 0 || patternBad < 0 {
		return nil, errors.New("patternGood or patternBad in a fake pattern with id '" + id + "' cannot be less than 0")
	}

	if patternGood == 0 && patternBad == 0 {
		return nil, errors.New("patternGood and patternBad in a fake pattern with id '" + id + "' cannot both be 0")
	}

	return &Pattern{
		id:           id,
		patternGood:  patternGood,
		patternBad:   patternBad,
		patternRatio: float64(patternGood) / float64(patternGood+patternBad),
		keepStats:    keepStats,
		Stats:        &PatternStats{Id: id},
	}, nil
}
