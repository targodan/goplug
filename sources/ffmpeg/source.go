package ffmpeg

import (
	"fmt"

	"github.com/targodan/goad"
	"github.com/targodan/goplug"

	"gopkg.in/targodan/ffgopeg.v1/avformat"
	"gopkg.in/targodan/ffgopeg.v1/avutil"
)

type InputFile struct {
	ohs         *goplug.OutputSocketHandler
	decoder     *goad.Decoder
	sampleRate  uint
	bufferSize  int
	streamIndex int
	numChannels int
	buffer      <-chan []float32
}

func NewInputFile(filename string, bufferSize int, sampleRate uint) (*InputFile, error) {
	file := &InputFile{}
	var err error
	file.decoder, err = goad.NewDecoder(filename)
	if err != nil {
		return nil, err
	}
	file.sampleRate = sampleRate
	file.bufferSize = bufferSize
	return file, nil
}

func (file *InputFile) Streams() []*avformat.Stream {
	return file.decoder.Streams()
}

func (file *InputFile) ChoseFirstStream() error {
	if file.ohs != nil {
		panic("Can only chose one stream.")
	}

	for i, s := range file.Streams() {
		if s.CodecPar().CodecType() == avutil.AVMEDIA_TYPE_AUDIO {
			return file.ChoseStream(i)
		}
	}

	return fmt.Errorf("No audio stream found.")
}

func (file *InputFile) ChoseStream(streamIndex int) error {
	if file.ohs != nil {
		panic("Can only chose one stream.")
	}

	ch, err, sr := file.decoder.EnableStream(streamIndex, file.bufferSize, int(file.sampleRate))
	if err != nil {
		return err
	}

	file.sampleRate = uint(sr)
	file.streamIndex = streamIndex

	file.buffer = ch

	file.ohs = goplug.NewOutputSocketHandler(file, file.Streams()[streamIndex].CodecPar().Channels())

	return nil
}

func (file InputFile) Start() <-chan error {
	return file.decoder.Start()
}

func (file InputFile) ChosenStream() *avformat.Stream {
	return file.Streams()[file.streamIndex]
}

func (file InputFile) Channels() int {
	return file.ChosenStream().CodecPar().Channels()
}

func (file InputFile) SampleRate() uint {
	return file.sampleRate
}

func (file *InputFile) Close() {
	file.decoder.Close()
}

// Output is the implementation of the goplug.Provider interface.
func (file InputFile) Output(i int) *goplug.OutputSocket {
	return file.ohs.GetSocket(i)
}

// Read reads a sample.
func (file *InputFile) Read() []goplug.Sample {
	samples, ok := <-file.buffer
	if !ok {
		// done reading. React and/or notify user.
	}
	ret := make([]goplug.Sample, len(samples))
	for i, s := range samples {
		ret[i] = goplug.Sample{s, file.SampleRate()}
	}
	return ret
}
