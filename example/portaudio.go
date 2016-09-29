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

	iDev, err := sources.NewDefaultInputDevice(1, 44100)
	if err != nil {
		panic(err)
	}
	defer iDev.Close()
	oDev, err := drains.NewDefaultOutputDevice(1, 44100)
	if err != nil {
		panic(err)
	}
	defer oDev.Close()

	dist := filters.NewDistortion()
	dist.SetGain(35)
	dist.SetVolume(0.1)
	split := filters.NewSplitter(2)
	merge := filters.NewMixer(2)

	goplug.Chain(iDev, 0, 0, split, 0, 0, merge, 0, 0, oDev)
	goplug.Chain(split, 1, 0, dist, 0, 1, merge, 0, 0, oDev)

	go oDev.Start()
	time.Sleep(60 * time.Second)
	oDev.Stop()
}
