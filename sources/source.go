package sources

import "github.com/targodan/goplug"

// Source implements the Output function of the goplug.Source interface.
// It is meant to be used as an anonymous field in a goplug.Source implementation.
type Source struct {
	ohs *goplug.OutputSocketHandler
}

// Output is the implementation of the goplug.Provider interface.
func (s Source) Output(i int) *goplug.OutputSocket {
	return s.ohs.GetSocket(i)
}
