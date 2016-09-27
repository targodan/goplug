package main

import (
	"time"

	"github.com/targodan/goplug/drains"
	"github.com/targodan/goplug/filters"
	"github.com/targodan/goplug/sources"
)

func main() {
	// D5
	d5 := sources.NewSineSource(587.330, 44100)
	// B5
	b5 := sources.NewSineSource(987.767, 44100)
	// E6
	e6 := sources.NewSineSource(1318.51, 44100)
	// D7
	d7 := sources.NewSineSource(2349.32, 44100)
	// Together it forms a Dm7 chord.
	mix := filters.NewMixer(4)
	raw, err := drains.NewRawFileDrain("test.raw")
	if err != nil {
		panic(err)
	}
	defer raw.Close()

	d5.Output(0).Plug(mix.Input(0))
	b5.Output(0).Plug(mix.Input(1))
	e6.Output(0).Plug(mix.Input(2))
	d7.Output(0).Plug(mix.Input(3))

	mix.Output(0).Plug(raw.Input(0))

	go raw.Start()
	time.Sleep(1 * time.Second)
	raw.Stop()
}
