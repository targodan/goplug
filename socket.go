package goplug

type InputSocket struct {
	conn *OutputSocket
}

func NewInputSocket() *InputSocket {
	return &InputSocket{}
}

func (is *InputSocket) Read() float32 {
	if is.conn == nil {
		return 0
	}
	return is.conn.Read()
}

func (is *InputSocket) Plug(s *OutputSocket) {
	is.conn = s
}

type OutputSocket struct {
	handler *OutputSocketHandler
}

func NewOutputSocket() *OutputSocket {
	return &OutputSocket{}
}

func (os *OutputSocket) Read() float32 {
	return os.handler.Read(os)
}

func (os *OutputSocket) registerHanlder(handler *OutputSocketHandler) {
	os.handler = handler
}

func (os *OutputSocket) Plug(s *InputSocket) {
	s.Plug(os)
}

type InputSocketHandler struct {
	sockets []*InputSocket
}

func NewInputSocketHandler(numSockets int) *InputSocketHandler {
	ret := &InputSocketHandler{}
	for i := 0; i < numSockets; i++ {
		ret.AddSocket(NewInputSocket())
	}
	return ret
}

func (ish *InputSocketHandler) AddSocket(s *InputSocket) {
	ish.sockets = append(ish.sockets, s)
}

func (ish *InputSocketHandler) GetSocket(i int) *InputSocket {
	return ish.sockets[i]
}

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

type OutputSocketHandler struct {
	provider      Provider
	sockets       []*OutputSocket
	indexes       map[*OutputSocket]int
	bufferIndexes []int
	buffer        *SampleBuffer
}

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

func (ish *OutputSocketHandler) AddSocket(s *OutputSocket) {
	ish.sockets = append(ish.sockets, s)
	ish.bufferIndexes = append(ish.bufferIndexes, ish.buffer.ReadIndex())
	ish.indexes[s] = len(ish.sockets) - 1
	s.registerHanlder(ish)
}

func (ish *OutputSocketHandler) GetSocket(i int) *OutputSocket {
	return ish.sockets[i]
}

func (ish *OutputSocketHandler) Read(sender *OutputSocket) float32 {
	var ret []float32
	index := ish.indexes[sender]
	if ish.buffer.ReadIndex() == ish.bufferIndexes[index] {
		ish.buffer.Write(ish.provider.Read())
	}

	ret, ish.bufferIndexes[index] = ish.buffer.GetSample(ish.bufferIndexes[index] + 1)
	return ret[index]
}
