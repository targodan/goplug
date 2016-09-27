package goplug

// InputSocket represents an input socket on an audio device.
type InputSocket struct {
	conn *OutputSocket
}

// NewInputSocket creates a new InputSocket instance.
func NewInputSocket() *InputSocket {
	return &InputSocket{}
}

// Read reads a sample from a connected OutputSocket or 0
// if none is connected.
func (is *InputSocket) Read() float32 {
	if is.conn == nil {
		return 0
	}
	return is.conn.Read()
}

// Plug connects this InputSocket to an OutputSocket.
func (is *InputSocket) Plug(s *OutputSocket) {
	is.conn = s
}

// OutputSocket represents an output socket on an audio device.
type OutputSocket struct {
	handler *OutputSocketHandler
}

// NewOutputSocket creates a new OutputSocket instance.
func NewOutputSocket() *OutputSocket {
	return &OutputSocket{}
}

// Read reads a sample from the Filter, Source or Drain it belongs to.
func (os *OutputSocket) Read() float32 {
	return os.handler.Read(os)
}

func (os *OutputSocket) registerHanlder(handler *OutputSocketHandler) {
	os.handler = handler
}

// Plug connects this OutputSocket to the given InputSocket.
func (os *OutputSocket) Plug(s *InputSocket) {
	s.Plug(os)
}

// InputSocketHandler is a helper for Filter or Drain implementations.
type InputSocketHandler struct {
	sockets []*InputSocket
}

// NewInputSocketHandler creates a new InputSocketHandler instance
// with the given number of InputSockets.
func NewInputSocketHandler(numSockets int) *InputSocketHandler {
	ret := &InputSocketHandler{}
	for i := 0; i < numSockets; i++ {
		ret.AddSocket(NewInputSocket())
	}
	return ret
}

// AddSocket adds an InputSocket to this Handler.
func (ish *InputSocketHandler) AddSocket(s *InputSocket) {
	ish.sockets = append(ish.sockets, s)
}

// GetSocket returns the i-th socket.
func (ish *InputSocketHandler) GetSocket(i int) *InputSocket {
	return ish.sockets[i]
}

// ReadAll reads one sample for from all InputSockets.
func (ish *InputSocketHandler) ReadAll() []float32 {
	if len(ish.sockets) == 1 {
		return []float32{ish.sockets[0].Read()}
	}

	ret := make([]float32, len(ish.sockets))
	wait := make(chan bool, len(ish.sockets))
	for i := 0; i < len(ish.sockets); i++ {
		go func(i int) {
			ret[i] = ish.sockets[i].Read()
			wait <- true
		}(i)
	}
	for i := 0; i < len(ish.sockets); i++ {
		<-wait
	}
	return ret
}

// OutputSocketHandler is a helper for Filter or Source implementations.
type OutputSocketHandler struct {
	provider      Provider
	sockets       []*OutputSocket
	indexes       map[*OutputSocket]int
	bufferIndexes []int
	buffer        *SampleBuffer
}

// NewOutputSocketHandler creates a new OutputSocketHandler instance
// with the given Provider and number of OutputSockets.
func NewOutputSocketHandler(provider Provider, numSockets int) *OutputSocketHandler {
	ret := &OutputSocketHandler{
		provider: provider,
		indexes:  make(map[*OutputSocket]int),
		buffer:   NewSampleBuffer(16),
	}
	for i := 0; i < numSockets; i++ {
		ret.AddSocket(NewOutputSocket())
	}
	return ret
}

// AddSocket adds an OutputSocket.
func (ish *OutputSocketHandler) AddSocket(s *OutputSocket) {
	ish.sockets = append(ish.sockets, s)
	ish.bufferIndexes = append(ish.bufferIndexes, ish.buffer.ReadIndex())
	ish.indexes[s] = len(ish.sockets) - 1
	s.registerHanlder(ish)
}

// GetSocket returns the i-th OutputSocket.
func (ish *OutputSocketHandler) GetSocket(i int) *OutputSocket {
	return ish.sockets[i]
}

// Read reads a sample from the given OutputSocket and buffers where necessary.
func (ish *OutputSocketHandler) Read(sender *OutputSocket) float32 {
	var ret []float32
	index := ish.indexes[sender]
	if ish.buffer.ReadIndex() == ish.bufferIndexes[index] {
		ish.buffer.Write(ish.provider.Read())
	}

	ret, ish.bufferIndexes[index] = ish.buffer.GetSample(ish.bufferIndexes[index] + 1)
	return ret[index]
}
