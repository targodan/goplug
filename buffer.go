package goplug

// SampleBuffer is a simple buffer for samples.
// It is not thread safe and on overflow earlier samples are
// overwritten without any notice.
type SampleBuffer struct {
	buffer     [][]Sample
	writeIndex int
}

// NewSampleBuffer creates a SampleBuffer instance of the specified size in samples.
func NewSampleBuffer(size int) *SampleBuffer {
	return &SampleBuffer{
		buffer:     make([][]Sample, size),
		writeIndex: 0,
	}
}

func (b *SampleBuffer) incWriteIndex() {
	b.writeIndex = (b.writeIndex + 1) % len(b.buffer)
}

// Write writes a sample to the buffer.
// If the size is exceeded the oldest sample is overwritten.
func (b *SampleBuffer) Write(sample []Sample) {
	b.buffer[b.writeIndex] = sample
	b.incWriteIndex()
}

// WriteIndex returns the next write index.
func (b *SampleBuffer) WriteIndex() int {
	return b.writeIndex
}

// ReadIndex returns the index of the youngest written sample.
func (b *SampleBuffer) ReadIndex() int {
	ret := b.writeIndex - 1
	if ret < 0 {
		ret += len(b.buffer)
	}
	return ret % len(b.buffer)
}

// GetSample returns the i-th sample.
// i is moduloed by the buffer size and then returned as well.
func (b *SampleBuffer) GetSample(i int) ([]Sample, int) {
	i = i % len(b.buffer)
	return b.buffer[i], i
}
