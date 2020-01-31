package fake

import (
	"encoding/json"
	"errors"
	"math"
)

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

type PatternStats struct {
	Id         string  `json:"id"`
	CTotal     int64   `json:"cumulativeTotal"`
	CGoodCount int64   `json:"cumulativeGoodCount"`
	CBadCount  int64   `json:"cumulativeBadCount"`
	CRatio     float64 `json:"cumulativeRatio"`
	Total      int64   `json:"slotTotal"`
	GoodCount  int64   `json:"slotGoodCount"`
	BadCount   int64   `json:"slotBadCount"`
	Ratio      float64 `json:"slotGoodRatio"`
}

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

func (ps *PatternStats) Json() string {
	out, _ := json.Marshal(ps)
	ps.Total = 0
	ps.GoodCount = 0
	ps.BadCount = 0
	ps.Ratio = 0
	return string(out)

}

func (fp *Pattern) Next() {
	fp.i = fp.i + 1
	f := math.Floor((float64(fp.i)/float64(fp.patternGood+fp.patternBad))*100000000) / 100000000
	fp.v = (fp.patternBad == 0 || ((f-math.Trunc(f)) <= fp.patternRatio && (f-math.Trunc(f)) != 0))
	if fp.keepStats {
		fp.Stats.Add(fp.v)
	}
}

func (fp *Pattern) Val() interface{} {
	return fp.v
}

func (fp *Pattern) Vals(count int) []interface{} {
	return makeValues(fp, count)
}

func (fp *Pattern) JsonStats() string {
	return fp.Stats.Json()
}

func (fp *Pattern) Good() bool {
	return fp.v
}

func (fp *Pattern) Bad() bool {
	return !fp.v
}

func (fp *Pattern) Values(count int) []bool {
	out := make([]bool, count)

	for i := 0; i < count; i++ {
		fp.Next()
		out[i] = fp.Good()
	}

	return out
}

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
