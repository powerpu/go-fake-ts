package fake

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

type Random struct {
	id        string
	rnd       *rand.Rand
	pctGood   float64
	keepStats bool
	Stats     *RandomStats
	v         bool
}

type RandomStats struct {
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

func (rs *RandomStats) Json() string {
	out, _ := json.Marshal(rs)
	rs.Total = 0
	rs.GoodCount = 0
	rs.BadCount = 0
	rs.Ratio = 0
	return string(out)

}

func (fr *Random) Next() {
	fr.v = fr.rnd.Float64() < fr.pctGood
	if fr.keepStats {
		fr.Stats.Add(fr.v)
	}
}

func (fr *Random) Val() interface{} {
	return fr.v
}

func (fr *Random) Vals(count int) []interface{} {
	return makeValues(fr, count)
}

func (fr *Random) JsonStats() string {
	return fr.Stats.Json()
}

func (fr *Random) Good() bool {
	return fr.v
}

func (fr *Random) Bad() bool {
	return !fr.v
}

func (fr *Random) Values(count int) []bool {
	out := make([]bool, count)

	for i := 0; i < count; i++ {
		fr.Next()
		out[i] = fr.Good()
	}

	return out
}

func NewRandom(id string, seed int64, pctGood float64, keepStats bool) (*Random, error) {
	if id == "" {
		return nil, errors.New("ID for a fake random cannot be blank")
	}

	if pctGood < 0 || pctGood > 1 {
		return nil, errors.New("Percentage good for a FakeRandom with id '" + id + "' must be between 0 and 1 but was '" + fmt.Sprintf("%v", pctGood) + "'")
	}

	return &Random{
		id:        id,
		rnd:       GenerateRandom(seed),
		pctGood:   pctGood,
		keepStats: keepStats,
		Stats:     &RandomStats{Id: id},
	}, nil
}
