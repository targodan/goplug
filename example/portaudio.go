package main

import (
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/targodan/goplug"
	drains "github.com/targodan/goplug/drains/portaudio"
	"github.com/targodan/goplug/filters"
	sources "github.com/targodan/goplug/sources/portaudio"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	iDev, err := sources.NewDefaultInputDevice()
	if err != nil {
		panic(err)
	}
	defer iDev.Close()
	oDev, err := drains.NewDefaultOutputDevice()
	if err != nil {
		panic(err)
	}
	defer oDev.Close()

	dist := filters.NewDistortion()
	dist.SetGain(35)
	dist.SetVolume(0.2)
	split := filters.NewSplitter(2)
	merge := filters.NewMixer(2)

	goplug.Chain(iDev, 0, 0, split, 0, 0, merge, 0, 0, oDev)
	goplug.Chain(split, 1, 0, dist, 0, 1, merge, 0, 1, oDev)

	oDev.Start()
	time.Sleep(5 * time.Second)
	oDev.Stop()
}
