package goplug

type SampleBuffer struct {
	buffer     [][]float32
	writeIndex int
}

func NewSampleBuffer(size int) *SampleBuffer {
	return &SampleBuffer{
		buffer:     make([][]float32, size),
		writeIndex: 0,
	}
}

func (b *SampleBuffer) incWriteIndex() {
	b.writeIndex = (b.writeIndex + 1) % len(b.buffer)
}

func (b *SampleBuffer) Write(sample []float32) {
	b.buffer[b.writeIndex] = sample
	b.incWriteIndex()
}

func (b *SampleBuffer) WriteIndex() int {
	return b.writeIndex
}

func (b *SampleBuffer) ReadIndex() int {
	ret := b.writeIndex - 1
	if ret < 0 {
		ret += len(b.buffer)
	}
	return ret % len(b.buffer)
}

func (b *SampleBuffer) GetSample(i int) ([]float32, int) {
	i = i % len(b.buffer)
	return b.buffer[i], i
}
