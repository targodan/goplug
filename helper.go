package goplug

import "fmt"

func chkType(ok bool, exp string, i int) {
	if !ok {
		panic(fmt.Sprintf("Invalid type. Expected element at position %d to be of type %s.", i, exp))
	}
}

func Chain(elements ...interface{}) {
	if len(elements) < 4 {
		return
	}

	for i := 0; i < len(elements)-1; i += 3 {
		from, ok := elements[i].(Provider)
		chkType(ok, "Provider", i)
		fromChan, ok := elements[i+1].(int)
		chkType(ok, "int", i+1)
		toChan, ok := elements[i+2].(int)
		chkType(ok, "int", i+2)
		to, ok := elements[i+3].(Consumer)
		chkType(ok, "Consumer", i+3)
		from.Output(fromChan).Plug(to.Input(toChan))
	}
}
