package goplug

// Provider is an interface for anything that provides audio data.
type Provider interface {
	// Output should return the i-th OutputSocket.
	Output(i int) *OutputSocket
	// Read should return one sample for each OutputSocket.
	Read() []float32
}

// Consumer is an interface for anything that consumes audio data.
type Consumer interface {
	// Input should return the i-th InputSocket.
	Input(i int) *InputSocket
}

// Drain is a Consumer that drains any connected providers.
type Drain interface {
	Consumer
	// Start should start the process of draining any connected Sources.
	Start()
}

// Filter is an interface for anything that both consumes and provides audio data.
type Filter interface {
	Provider
	Consumer
}
