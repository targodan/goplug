package drains

import "github.com/targodan/goplug"

type Drain struct {
	ihs *goplug.InputSocketHandler
}

func (d Drain) Input(i int) *goplug.InputSocket {
	return d.ihs.GetSocket(i)
}
