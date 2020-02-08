package fake

import (
	"encoding/json"
	"github.com/jhorwit2/simple-regression"
	"math"
	"math/rand"
)

// Data generates fake values based on desired parameters.
type Data struct {
	id string

	// Common variables
	samples      int64
	stretchStart float64
	stretchEnd   float64
	stretchStep  float64
	slope        float64
	bump         float64
	from         float64
	to           float64
	limitUpper   bool
	limitLower   bool

	// Permanent bump variables
	permaBumpAt       int64
	permaBumpBy       float64
	permaBumpSmoother int64

	// Random variables
	useRandom bool
	rnd       *rand.Rand
	bias      float64

	// Spike variables
	spike             bool
	spikeEvery        int64
	spikeSustain      int64
	spikeTo           int64
	spikeWobble       bool
	spikeWobbleFactor int64
	spikeSmoother     int64

	// Seasonality variables
	seasonality      bool
	seasonalityWave1 int64
	seasonalityWave2 int64
	seasonalityWave3 int64
	seasonalityWave4 int64
	seasonalityWave5 int64

	// Stats variables
	keepStats bool
	Stats     *DataStats

	// Runtime variables
	spikeCount int64
	spikeStart int64
	spikeEnd   int64
	b          float64
	f          float64
	i          int64
	v          float64
}

// DataStats keeps track of various statistics of a Data while it's running.
type DataStats struct {
	// The ID of the Data
	ID string `json:"id"`

	// From of the Data
	From float64 `json:"from"`

	// To of the Data
	To float64 `json:"to"`

	// Random seed of the Data
	Seed int64 `json:"seed"`

	// Cumulative count of how many times Next() was called.
	CTotal int64 `json:"cumulativeTotal"`

	// Cumulative minimum value
	CMin float64 `json:"cumulativeMinimum"`

	// Cumulative spot where "from" is reached (if applicable)
	CHitMinAt int64 `json:"cumulativeHitMinimumAt"`

	// Cumulative number of points below the "from" value (if applicable)
	CPointsLessThan int64 `json:"cumulativePointsBelowLowerLimit"`

	// Cumulative number of points at the "from" value (if applicable)
	CPointsAtLower int64 `json:"cumulativePointsAtLowerLimit"`

	// Cumulative maximum value
	CMax float64 `json:"cumulativeMaximum"`

	// Cumulative spot where "to" is reached (if applicable)
	CHitMaxAt int64 `json:"cumulativeHitMaximumAt"`

	// Cumulative number of points above the "to" value (if applicable)
	CPointsMoreThan int64 `json:"cumulativePointsAboveUpperLimit"`

	// Cumulative number of points at the "to" value (if applicable)
	CPointsAtUpper int64 `json:"cumulativePointsAtUpperLimit"`

	cRegression *linear.Regression

	// Cumulative slope of the data
	CSlope float64 `json:"cumulativeSlope"`

	// Slot count of how many times Next() was called. This gets reset after every JSON() call.
	Total int64 `json:"slotTotal"`

	// Slot minimum value
	Min float64 `json:"slotMinimum"`

	// Slot spot where "from" is reached (if applicable)
	HitMinAt int64 `json:"slotHitMinimumAt"`

	// Slot number of points below the "from" value (if applicable)
	PointsLessThan int64 `json:"slotPointsBelowLowerLimit"`

	// Slot number of points at the "from" value (if applicable)
	PointsAtLower int64 `json:"slotPointsAtLowerLimit"`

	// Slot maximum value
	Max float64 `json:"slotMaximum"`

	// Slot spot where "to" is reached (if applicable)
	HitMaxAt int64 `json:"slotHitMaximumAt"`
	// Slot number of points above the "to" value (if applicable)
	PointsMoreThan int64 `json:"slotPointsAboveUpperLimit"`

	// Slot number of points at the "to" value (if applicable)
	PointsAtUpper int64 `json:"slotPointsAtUpperLimit"`

	regression *linear.Regression

	// Slope of the current slot
	Slope float64 `json:"slotSlope"`
}

// Add adds a value to the running tally.
func (ds *DataStats) Add(v float64) {

	// Initial values
	if ds.CTotal == 0 {
		ds.CMin = v
		ds.Min = v
		ds.CMax = v
		ds.Max = v
	}

	ds.cRegression.Push(float64(ds.CTotal), v)
	ds.regression.Push(float64(ds.Total), v)

	ds.CTotal++
	ds.Total++

	if v < ds.CMin {
		ds.CMin = v
	}

	if v < ds.Min {
		ds.Min = v
	}

	if ds.CHitMinAt == 0 && v <= ds.From {
		ds.CHitMinAt = ds.CTotal
	}

	if ds.HitMinAt == 0 && v <= ds.From {
		ds.HitMinAt = ds.Total
	}

	if v > ds.CMax {
		ds.CMax = v
	}

	if v > ds.Max {
		ds.Max = v
	}

	if ds.CHitMaxAt == 0 && v >= ds.To {
		ds.CHitMaxAt = ds.CTotal
	}

	if ds.HitMaxAt == 0 && v >= ds.To {
		ds.HitMaxAt = ds.Total
	}

	if v < ds.From {
		ds.CPointsLessThan++
		ds.PointsLessThan++
	}

	if v == ds.From {
		ds.CPointsAtLower++
		ds.PointsAtLower++
	}

	if v == ds.To {
		ds.CPointsAtUpper++
		ds.PointsAtUpper++
	}

	if v > ds.To {
		ds.CPointsMoreThan++
		ds.PointsMoreThan++
	}
}

// JSON returns a summary of the current pattern statistics and resets the slot
// tally.
func (ds *DataStats) JSON() string {
	ds.CSlope = ds.cRegression.Slope()
	ds.Slope = ds.regression.Slope()
	out, _ := json.MarshalIndent(ds, "", " ")
	ds.Total = 0
	ds.Min = 0
	ds.HitMinAt = 0
	ds.PointsLessThan = 0
	ds.PointsAtLower = 0
	ds.Max = 0
	ds.HitMaxAt = 0
	ds.PointsMoreThan = 0
	ds.PointsAtUpper = 0
	ds.regression = linear.NewRegression()
	ds.Slope = 0
	return string(out)

}

func (fd *Data) calculateNextSpikeStartStop() {
	fd.spikeStart = fd.spikeCount * (fd.spikeEvery - fd.spikeSmoother)
	fd.spikeEnd = fd.spikeStart + (2 * fd.spikeSmoother) + fd.spikeSustain

	for i := (fd.spikeCount + 1); fd.spikeStart < fd.i; i++ {
		fd.spikeStart = i * (fd.spikeEvery - fd.spikeSmoother)
		fd.spikeEnd = fd.spikeStart + (2 * fd.spikeSmoother) + fd.spikeSustain
		if fd.spikeStart <= 0 {
			break
		}
	}
}

// Next generates the next fake number value.
func (fd *Data) Next() {
	// We're at the end of a spike, let's calculate when the next spike starts
	if fd.i == 0 || fd.i == fd.spikeEnd {
		fd.spikeCount++
		fd.calculateNextSpikeStartStop()
	}

	spread := math.Abs(fd.to) + math.Abs(fd.from)

	// Column A
	a := fd.rnd.Float64()

	// Column E
	e := (fd.from + fd.to) / 2

	if fd.useRandom {
		// Column B
		if fd.i == 0 {
			// fd.b = a * fd.to
			fd.b = a
		}

		bPrev := fd.b
		if fd.i > 0 {
			fd.b = bPrev
		}

		if fd.bias <= 0 {
			fd.b = bPrev + fd.b
		} else if fd.bias >= 1 {
			fd.b = bPrev - fd.b
		} else if a > fd.bias {
			fd.b = bPrev + (((a - 0.5) * (a - 0.5)) * -1)
		} else if a < fd.bias {
			fd.b = bPrev + ((a - 0.5) * (a - 0.5))
		}

		// Column C
		c := fd.b + (math.Log(float64(fd.samples)) / math.Log(2.5))

		// Column D
		d := ((c / (math.Log(float64(fd.samples)) / math.Log(2.5))) / 2)

		e = (d * spread) - spread
	}

	// Let's do seasonality!
	sv := float64(0)
	if fd.seasonality {
		divisor := float64(0)
		rad1 := (math.Sin(float64(fd.i)*degToRad(float64(1)/(float64(fd.seasonalityWave1)/float64(360)))) / 2) + 0.5
		if fd.seasonalityWave1 == 1 {
			rad1 = 0
		} else {
			divisor++
		}

		rad2 := (math.Sin(float64(fd.i)*degToRad(float64(1)/(float64(fd.seasonalityWave2)/float64(360)))) / 2) + 0.5
		if fd.seasonalityWave2 == 1 {
			rad2 = 0
		} else {
			divisor++
		}

		rad3 := (math.Sin(float64(fd.i)*degToRad(float64(1)/(float64(fd.seasonalityWave3)/float64(360)))) / 2) + 0.5
		if fd.seasonalityWave3 == 1 {
			rad3 = 0
		} else {
			divisor++
		}

		rad4 := (math.Sin(float64(fd.i)*degToRad(float64(1)/(float64(fd.seasonalityWave4)/float64(360)))) / 2) + 0.5
		if fd.seasonalityWave4 == 1 {
			rad4 = 0
		} else {
			divisor++
		}

		rad5 := (math.Sin(float64(fd.i)*degToRad(float64(1)/(float64(fd.seasonalityWave5)/float64(360)))) / 2) + 0.5
		if fd.seasonalityWave5 == 1 {
			rad5 = 0
		} else {
			divisor++
		}

		if divisor > 0 {
			sv = (spread * ((rad1 + rad2 + rad3 + rad4 + rad5) / divisor)) - (spread / 2)
		}
	}

	f := e + sv // Column F

	// Let's do the permanent bump which includes a smoother as we're wither
	// going up or down from a baseline
	if fd.i >= fd.permaBumpAt && fd.permaBumpAt > 0 && fd.permaBumpSmoother > 0 {
		bumpBaseline := (float64(fd.permaBumpBy) / 100) * fd.to
		if fd.i-fd.permaBumpAt > fd.permaBumpSmoother {
			f = f + bumpBaseline
		} else {
			tmp := (float64(fd.i-fd.permaBumpAt) / float64(fd.permaBumpSmoother))
			adjusted := bumpBaseline * tmp * tmp
			f = f + adjusted
		}
	}

	f = f + (float64(fd.i) * fd.slope) + fd.bump
	v := float64(0)

	// Let's stretch or squish. We need to do it both up and down.
	stv := fd.stretchStart + (float64(fd.i) * fd.stretchStep)
	if stv > 1 { // Stretch above 1 means "stretch"
		if f > fd.f { // Stretching on the way up
			v = f + ((f - fd.f) * stv)
		} else if f < fd.f { // Stretching on the way down
			v = f - ((fd.f - f) * stv)
		} else { // Flat value so just use whatever the previous value was
			v = fd.v
		}
	} else if stv < 1 { // Stretch below 1 means "squish"
		v = f * stv
	}

	// Let's do spikes!
	if fd.spike && fd.i >= fd.spikeStart && !(fd.i > fd.spikeEnd) {
		multiplier := int64(0)
		spikeValue := (float64(fd.spikeTo) / 100) * fd.to

		// Quick div/0 safety check
		if fd.spikeWobbleFactor == 0 {
			fd.spikeWobbleFactor = 1
		}

		// Quick div/0 safety check
		if fd.spikeSmoother == 0 {
			fd.spikeSmoother = 1
		}

		if fd.i >= fd.spikeStart && fd.i < (fd.spikeStart+fd.spikeSmoother) { // Going up?
			multiplier = fd.spikeSmoother - ((fd.spikeStart + fd.spikeSmoother) - fd.i) + 1
		} else if fd.i > (fd.spikeEnd - fd.spikeSmoother) { // Going down?
			multiplier = fd.spikeEnd - fd.i + 1
		}

		if multiplier == 0 {
			if fd.spikeWobble {
				if fd.spikeWobbleFactor > 0 {
					v = spikeValue - ((a * spikeValue) / float64(fd.spikeWobbleFactor))
				} else {
					v = spikeValue + ((a * spikeValue) / float64(fd.spikeWobbleFactor))
				}
			} else {
				v = spikeValue
			}
		} else {
			// Let's apply a smoother as we're wither going up or down from the peak
			tmp := (float64(1) - math.Abs(1/(float64(multiplier)))) * (1 + (1 / float64(fd.spikeSmoother)))
			v = v + (tmp * tmp * (spikeValue - v))
		}
	}

	// Let's limit
	if fd.limitLower && v < fd.from {
		v = fd.from
	} else if fd.limitUpper && v > fd.to {
		v = fd.to
	}

	// Setup next iteration
	fd.f = f
	fd.i = fd.i + 1
	fd.v = v

	if fd.keepStats {
		fd.Stats.Add(fd.v)
	}
}

// Val returns the current fake numeric value.
func (fd *Data) Val() interface{} {
	return fd.v
}

// Vals returns the next count of fake values as an interface{} array.
func (fd *Data) Vals(count int) []interface{} {
	return makeValues(fd, count)
}

// JSONStats retrieves the current stats as s JSON string.
func (fd *Data) JSONStats() string {
	return fd.Stats.JSON()
}

// Float returns the current value as a float64.
func (fd *Data) Float() float64 {
	return fd.v
}

// Floats returns the next count of fake values as a float64 array.
func (fd *Data) Floats(count int) []float64 {
	out := make([]float64, count)

	for i := 0; i < count; i++ {
		out[i] = fd.Float()
		fd.Next()
	}

	return out
}

// NewData create a new fake data. Data has a unique id, number of required
// samples to generate. Other parameters are:
//
// stretchStart and stretchEnd
//
// When a default graph is generated this value "stretches" or "squishes" the
// data up and down. For example if the generated data has a minimum of 35 and
// a maximum of 60, "stretching" it will increase the difference between the
// minumum and maximum while "squishing" will decrease it.
//
// A value above 1 will stretch, a value below 1 but more than 0 will "squish"
// while a value of 1 will be "normal".
//
// stretchStart and stretchEnd values indicates the start and stop values.
// This is useful if you want to "amplify over time" or "dampen over time".
//
// For example if stretchStart is 100, stretchEnd is 0 and we want 100
// samples then we will gradually reduce the stretch from 100 to 0 over 100
// linear steps (i.e. the stretch value at the 50th sample will be 50).
//
// slope
//
// This will set the slope of the generated data. 0 means whatever the
// underlying data is. Positive values means data will trend upwards, negative
// values mean data will trend downward.
//
// bump
//
// Tweaking this number will "bump" the random value up (or down if negative).
// For example if a maximum is 60, a positive "bump" will increase this value,
// a negative "bump" will decrease it. This is useful if you want to ensure
// values breach min/max in certain cases.
//
// from
//
// Generate data "from" this number. E.g. CPU values may be from 0 to 100.
//
// to
//
// Generate data "to" this number. E.g. CPU values may be from 0 to 100.
//
// limitUpper
//
// If numbers go above "to", when this is "TRUE" values will be set to the
// "to" value. E.g. a CPU cannot go above 100%.
//
// limitLower
//
// If numbers go below "from", when this is "TRUE" values will be set to the
// "from" value. E.g. a CPU cannot go below 0%.
//
// permaBumpAt
//
// Generate data "from" this number. E.g. CPU values may be from 0 to 100. Use
// 0 to disable permanent bump.
//
// permaBumpBy
//
// What value should we bump to, expressed as a percentage of "to". Can be
// negative too.
//
// permaBumpSmoother
//
// When smoother is 1 we will go straight from current value to bump value.
// When smoother is above 1 when we will take this number of samples to reach
// the bump value (i.e. gradually but quickly rise to a spike). Use 0 to
// disable permanent bump.
//
// useRandom
//
// Whether to generate random numbers. When "TRUE" numbers will be generated
// based on the below parameters.
//
// seed
//
// A seed to use for any random numbers used internally
//
// bias
//
// A parameter between 0 and 1. Use it to control the range and slope of your
// data. Generally lower biases make the slope negative but which values
// entirely depend on the underlying dataset. As always, have a play to
// see the effects.
//
// spike
//
// Whether to generate spikes. When "TRUE" spikes will be generated based on
// the below parameters.
//
// spikeEvery
//
// Starting at 0, every n samples will reach the desired spike value.
//
// spikeSustain
//
// When a spike is reached, for how many samples should we sustain it.
//
// spikeTo
//
// What value should we spike to, expressed as a percentage of "to".
//
// spikeWobble
//
// When we're spiking and sustaining it, do we use a flat value at the top or
// do we "bounce off the top" a little for a little variation? When "TRUE" it
// means that when sustaining we'll add a little variance. We respect and will
// never break the SPIKE_TO value though!
//
// spikeWobbleFactor
//
// A magic number to tweak the "wobbliness". Have a play with this value.
// Generally a higher value will mean smoother values.
//
// spikeSmoother
//
// When smoother is 0 we will go straight from current value to spike value.
// When smoother is above 0 when we will take this number of samples to reach
// the spike value (i.e. gradually but quickly rise to a spike)
//
// seasonality
//
// Whether to generate seasonality using SIN. When "TRUE" waves will be
// generated using the WAVE parameters below.
//
// seasonalityWave[1-5]
//
// Indicates number of points where one SIN cycle will be complete.  Each wave
// is summed to generate interference.
func NewData(
	id string,
	samples int64,

	stretchStart float64,
	stretchEnd float64,
	slope float64,
	bump float64,
	from float64,
	to float64,
	limitUpper bool,
	limitLower bool,

	permaBumpAt int64,
	permaBumpBy float64,
	permaBumpSmoother int64,

	useRandom bool,
	seed int64,
	bias float64,

	spike bool,
	spikeEvery int64,
	spikeSustain int64,
	spikeTo int64,
	spikeWobble bool,
	spikeWobbleFactor int64,
	spikeSmoother int64,

	seasonality bool,
	seasonalityWave1 int64,
	seasonalityWave2 int64,
	seasonalityWave3 int64,
	seasonalityWave4 int64,
	seasonalityWave5 int64,

	keepStats bool) (*Data, error) {

	stretchStep := math.Abs(math.Abs(stretchEnd)-math.Abs(stretchStart)) / float64(samples)
	if stretchEnd < stretchStart {
		stretchStep = stretchStep * -1
	}

	d := &Data{
		id:      id,
		samples: samples,

		stretchStart: stretchStart,
		stretchEnd:   stretchEnd,
		stretchStep:  stretchStep,
		slope:        slope,
		bump:         bump,
		from:         from,
		to:           to,
		limitUpper:   limitUpper,
		limitLower:   limitLower,

		useRandom: useRandom,
		rnd:       generateRandom(seed),
		bias:      bias,

		permaBumpAt:       permaBumpAt,
		permaBumpBy:       permaBumpBy,
		permaBumpSmoother: permaBumpSmoother,

		spike:             spike,
		spikeSustain:      spikeSustain,
		spikeEvery:        spikeEvery,
		spikeTo:           spikeTo,
		spikeWobble:       spikeWobble,
		spikeWobbleFactor: spikeWobbleFactor,
		spikeSmoother:     spikeSmoother,

		seasonality:      seasonality,
		seasonalityWave1: seasonalityWave1,
		seasonalityWave2: seasonalityWave2,
		seasonalityWave3: seasonalityWave3,
		seasonalityWave4: seasonalityWave4,
		seasonalityWave5: seasonalityWave5,

		keepStats: keepStats,
		Stats: &DataStats{
			ID:          id,
			From:        from,
			To:          to,
			Seed:        seed,
			regression:  linear.NewRegression(),
			cRegression: linear.NewRegression(),
		},
	}

	d.Next()
	return d, nil
}
