package portaudio

import (
	"github.com/gordonklaus/portaudio"
	"github.com/targodan/goplug"
)

type InputDevice struct {
	ohs         *goplug.OutputSocketHandler
	stream      *portaudio.Stream
	buffer      chan []float32
	sampleRate  uint
	numChannels int
}

func NewDefaultInputDevice(channels int, sampleRate uint) (*InputDevice, error) {
	device, err := portaudio.DefaultInputDevice()
	if err != nil {
		return nil, err
	}
	params := portaudio.LowLatencyParameters(device, nil)
	params.Input.Channels = channels
	params.SampleRate = float64(sampleRate)
	return NewInputDevice(params)
}

func NewInputDevice(params portaudio.StreamParameters) (*InputDevice, error) {
	dev := &InputDevice{
		stream: nil,
		buffer: make(chan []float32, 64),
	}
	stream, err := portaudio.OpenStream(params, dev.recv)
	dev.sampleRate = uint(stream.Info().SampleRate)
	dev.numChannels = params.Input.Channels
	if err != nil {
		return nil, err
	}
	dev.stream = stream
	dev.ohs = goplug.NewOutputSocketHandler(dev, dev.numChannels)
	dev.stream.Start()
	return dev, nil
}

func (dev *InputDevice) recv(data [][]float32) {
	for i := range data[0] {
		tmp := make([]float32, dev.numChannels)
		for c := 0; c < dev.numChannels; c++ {
			tmp[c] = data[c][i]
		}
		dev.buffer <- tmp
	}
}

func (dev *InputDevice) Close() {
	dev.stream.Close()
}

// Output is the implementation of the goplug.Provider interface.
func (dev InputDevice) Output(i int) *goplug.OutputSocket {
	return dev.ohs.GetSocket(i)
}

// Read reads a sample.
func (dev *InputDevice) Read() []goplug.Sample {
	ch := <-dev.buffer
	ret := make([]goplug.Sample, dev.numChannels)
	for i, s := range ch {
		ret[i] = goplug.Sample{s, dev.sampleRate}
	}
	return ret
}
