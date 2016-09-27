package filters

import "github.com/targodan/goplug"

// Splitter is a goplug.Filter with one input and multiple outputs.
// Each output channel contains the same samples as the input receives.
type Splitter struct {
	Filter
	numOutputs int
}

// NewSplitter creates a new Splitter instance.
func NewSplitter(numOutputs int) *Splitter {
	ret := &Splitter{
		numOutputs: numOutputs,
	}
	ret.ohs = goplug.NewOutputSocketHandler(ret, numOutputs)
	ret.ihs = goplug.NewInputSocketHandler(1)
	return ret
}

// Read reads a Sample.
func (s Splitter) Read() []float32 {
	val := s.ihs.GetSocket(0).Read()
	ret := make([]float32, s.numOutputs)
	for i := 0; i < s.numOutputs; i++ {
		ret[i] = val
	}
	return ret
}
