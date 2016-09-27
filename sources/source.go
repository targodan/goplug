package sources

import "github.com/targodan/goplug"

type Source struct {
	ohs *goplug.OutputSocketHandler
}

func (s Source) Output(i int) *goplug.OutputSocket {
	return s.ohs.GetSocket(i)
}
