package filters

import (
	"fmt"

	"github.com/targodan/goplug"
)

// Mixer is a goplug.Filter implementation with multiple inputs
// and one output. It mixes the channels together with a 1/n factor
// whereas n is the number of channels.
//
// Additionally you can manipulate levels of each individual channel
// as well as a master level.
type Mixer struct {
	Filter
	numChannels int
	levels      []float32
	masterLevel float32
}

// NewMixer creates a new Mixer instance.
func NewMixer(numChannels int) *Mixer {
	ret := &Mixer{
		numChannels: numChannels,
		levels:      make([]float32, numChannels),
		masterLevel: 1,
	}
	ret.ohs = goplug.NewOutputSocketHandler(ret, 1)
	ret.ihs = goplug.NewInputSocketHandler(numChannels)
	for i := 0; i < numChannels; i++ {
		ret.levels[i] = 1
	}
	return ret
}

// Read reads a sample.
func (m Mixer) Read() []goplug.Sample {
	channels := m.ihs.ReadAll()
	var ret goplug.Sample
	ret.SampleFrequency = 0
	for i, v := range channels {
		if ret.SampleFrequency != 0 && v.SampleFrequency != 0 && ret.SampleFrequency != v.SampleFrequency {
			panic(fmt.Sprintf("Incompatible sample frequencies expected %d but got %d. Please use a resampler first.", ret.SampleFrequency, v.SampleFrequency))
		}
		ret.SampleFrequency = v.SampleFrequency
		ret.Value += v.Value * m.levels[i]
	}
	ret.Value = m.masterLevel * ret.Value / float32(m.numChannels)
	return []goplug.Sample{ret}
}

// SetLevel sets the level for a channel.
func (m *Mixer) SetLevel(channel int, level float32) {
	m.levels[channel] = level
}

// SetMasterLevel sets the master level.
func (m *Mixer) SetMasterLevel(level float32) {
	m.masterLevel = level
}
