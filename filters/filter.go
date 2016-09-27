package filters

import "github.com/targodan/goplug"

// Filter implements the Output and Input functions of the
// goplug.Provider and goplug.Consumer interfaces.
// This can be used as an anonymous field in a goplug.Filter
// implementation.
type Filter struct {
	ihs *goplug.InputSocketHandler
	ohs *goplug.OutputSocketHandler
}

// Output is the implementaiton of the goplug.Provider interface.
func (s Filter) Output(i int) *goplug.OutputSocket {
	return s.ohs.GetSocket(i)
}

// Input is the implementaiton of the goplug.Consumer interface.
func (s Filter) Input(i int) *goplug.InputSocket {
	return s.ihs.GetSocket(i)
}
