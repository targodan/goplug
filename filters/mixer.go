package filters

import "github.com/targodan/goplug"

type Mixer struct {
	Filter
	numChannels int
	levels      []float32
	masterLevel float32
}

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

func (m Mixer) Read() []float32 {
	channels := m.ihs.ReadAll()
	var ret float32
	for i, v := range channels {
		ret += v * m.levels[i]
	}
	return []float32{m.masterLevel * ret / float32(m.numChannels)}
}

func (m *Mixer) SetLevel(channel int, level float32) {
	m.levels[channel] = level
}

func (m *Mixer) SetMasterLevel(level float32) {
	m.masterLevel = level
}
