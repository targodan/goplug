package sources

import (
	"math"

	"github.com/targodan/goplug"
)

// SineSource is a goplug.Source that is capable of creating sine waves.
// The frequency of the sine wave can still be manipulated after a connected
// golang.Drain has started.
type SineSource struct {
	Source
	sampleFreq  uint
	deltaT      float64
	freq        float64
	t           float64
	sampleIndex uint64
}

// NewSineSource creates a new SineSource instance.
// Make sure that the sample frequency is the same as the drains
// or resample the signal later.
func NewSineSource(freq float64, sampleFreq uint) *SineSource {
	ret := &SineSource{
		sampleFreq:  sampleFreq,
		freq:        freq,
		t:           0,
		sampleIndex: 0,
	}
	ret.calculateDeltaT()
	ret.ohs = goplug.NewOutputSocketHandler(ret, 1)
	return ret
}

// SetFrequency sets the frequency of the sine wave.
func (s *SineSource) SetFrequency(freq float64) {
	s.freq = freq
	s.calculateDeltaT()
}

func (s *SineSource) calculateDeltaT() {
	s.deltaT = 2 * math.Pi * s.freq / float64(s.sampleFreq)
}

// Read reads a sample.
func (s *SineSource) Read() []goplug.Sample {
	s.t += s.deltaT
	return []goplug.Sample{goplug.Sample{float32(math.Sin(s.t)), s.sampleFreq}}
}
