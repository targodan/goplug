package sources

import (
	"math"

	"github.com/targodan/goplug"
)

type SineSource struct {
	Source
	sampleFreq  uint
	deltaT      float64
	freq        float64
	t           float64
	sampleIndex uint64
}

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

func (s *SineSource) SetFrequency(freq float64) {
	s.freq = freq
	s.calculateDeltaT()
}

func (s *SineSource) calculateDeltaT() {
	s.deltaT = 2 * math.Pi * s.freq / float64(s.sampleFreq)
}

func (s *SineSource) Read() []float32 {
	s.t += s.deltaT
	return []float32{float32(math.Sin(s.t))}
}
