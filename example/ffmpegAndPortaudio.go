package main

import (
	"fmt"
	"os"

	"github.com/gordonklaus/portaudio"
	"github.com/targodan/goplug"
	drains "github.com/targodan/goplug/drains/portaudio"
	sources "github.com/targodan/goplug/sources/ffmpeg"
)

const numOutChannels = 2

func main() {
	// Handle missuse.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run decode.go <audiofile>")
		os.Exit(-1)
	}

	inFilename := os.Args[1]

	portaudio.Initialize()
	defer portaudio.Terminate()

	file, err := sources.NewInputFile(inFilename, 10000, 44100)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = file.ChoseFirstStream()
	if err != nil {
		panic(err)
	}

	fmt.Printf("File: \"%s\" Sample rate: %d Stream: %d Channels: %d\n", inFilename, file.SampleRate(), file.ChosenStream().StreamIdentifier(), file.Channels())

	errChan := file.Start()
	go func() {
		for err := range errChan {
			if err != nil {
				panic(err)
			}
		}
	}()

	oDev, err := drains.NewDefaultOutputDevice(numOutChannels, file.SampleRate())
	if err != nil {
		panic(err)
	}
	defer oDev.Close()

	maxChannels := 0
	if numOutChannels <= file.Channels() {
		maxChannels = numOutChannels
	} else {
		maxChannels = file.Channels()
	}
	for i := 0; i < maxChannels; i++ {
		goplug.Chain(file, i, i, oDev)
	}

	go oDev.Start()
	defer oDev.Stop()
}
