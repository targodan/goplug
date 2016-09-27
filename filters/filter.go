package filters

import "github.com/targodan/goplug"

type Filter struct {
	ihs *goplug.InputSocketHandler
	ohs *goplug.OutputSocketHandler
}

func (s Filter) Output(i int) *goplug.OutputSocket {
	return s.ohs.GetSocket(i)
}

func (s Filter) Input(i int) *goplug.InputSocket {
	return s.ihs.GetSocket(i)
}
