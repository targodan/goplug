package portaudio

import (
	"github.com/gordonklaus/portaudio"
	"github.com/targodan/goplug"
)

type OutputDevice struct {
	ihs         *goplug.InputSocketHandler
	stream      *portaudio.Stream
	buffer      chan []float32
	sampleRate  uint
	hasStopped  chan bool
	running     bool
	numChannels int
}

func NewDefaultOutputDevice(channels int, sampleRate uint) (*OutputDevice, error) {
	device, err := portaudio.DefaultOutputDevice()
	if err != nil {
		return nil, err
	}
	params := portaudio.LowLatencyParameters(nil, device)
	params.Output.Channels = channels
	params.SampleRate = float64(sampleRate)
	return NewOutputDevice(params)
}

func NewOutputDevice(params portaudio.StreamParameters) (*OutputDevice, error) {
	dev := &OutputDevice{
		stream: nil,
		buffer: make(chan []float32, 64),
	}
	err := portaudio.IsFormatSupported(params, dev.send)
	if err != nil {
		panic(err)
	}
	stream, err := portaudio.OpenStream(params, dev.send)
	dev.sampleRate = uint(stream.Info().SampleRate)
	dev.numChannels = params.Output.Channels
	if err != nil {
		return nil, err
	}
	dev.stream = stream
	dev.ihs = goplug.NewInputSocketHandler(dev.numChannels)
	return dev, nil
}

func (dev *OutputDevice) send(data [][]float32) {
	for i := range data[0] {
		ch := <-dev.buffer
		for c := 0; c < dev.numChannels; c++ {
			data[c][i] = ch[c]
		}
	}
}

func (dev *OutputDevice) Close() {
	dev.stream.Close()
}

// Output is the implementation of the goplug.Provider interface.
func (dev OutputDevice) Input(i int) *goplug.InputSocket {
	return dev.ihs.GetSocket(i)
}

func (dev *OutputDevice) Start() {
	dev.hasStopped = make(chan bool)
	dev.running = true
	dev.stream.Start()
	for dev.running {
		cs := dev.ihs.ReadAll()
		tmp := make([]float32, dev.numChannels)
		for i, s := range cs {
			if s.SampleFrequency != 0 && s.SampleFrequency != dev.sampleRate {
				dev.running = false
				break
			}
			tmp[i] = s.Value
		}
		dev.buffer <- tmp
	}
	dev.hasStopped <- true
}

func (dev *OutputDevice) Stop() {
	dev.running = false
	<-dev.hasStopped
}
