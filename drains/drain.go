package drains

import "github.com/targodan/goplug"

// Drain is a struct that is supposed to make it easier to create goplug.Provider implementations.
type Drain struct {
	ihs *goplug.InputSocketHandler
}

// Input is the implementation of the goplug.Provider interface.
func (d Drain) Input(i int) *goplug.InputSocket {
	return d.ihs.GetSocket(i)
}
