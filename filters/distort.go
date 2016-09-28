package filters

import (
	"errors"
	"math"

	"github.com/targodan/goplug"
)

type Distortion struct {
	Filter
	gain   float32
	volume float32
}

func NewDistortion() *Distortion {
	ret := &Distortion{gain: 1, volume: 1}

	ret.ohs = goplug.NewOutputSocketHandler(ret, 1)
	ret.ihs = goplug.NewInputSocketHandler(1)

	return ret
}

func (d *Distortion) SetGain(g float32) {
	d.gain = g
}

func (d *Distortion) SetVolume(v float32) error {
	if v < 0 || v > 1 {
		return errors.New("Volume must be between 0 and 1.")
	}
	d.volume = v
	return nil
}

// Read reads a Sample.
func (d Distortion) Read() []goplug.Sample {
	val := d.ihs.GetSocket(0).Read()
	val.Value = d.volume * float32(math.Tanh(float64(d.gain*val.Value)))
	return []goplug.Sample{val}
}
