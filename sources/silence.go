package sources

import "github.com/targodan/goplug"

type Silence struct {
	Source
}

func NewSilence() *Silence {
	ret := &Silence{}
	ret.ohs = goplug.NewOutputSocketHandler(ret, 1)
	return ret
}

// Read reads a sample.
func (s *Silence) Read() []goplug.Sample {
	return []goplug.Sample{goplug.Sample{float32(0), 0}}
}
