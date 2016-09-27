package goplug

// Provider is an interface for anything that provides audio data.
type Provider interface {
	Output(i int) *OutputSocket
	Read() []float32
}

// Consumer is an interface for anything that consumes audio data.
type Consumer interface {
	Input(i int) *InputSocket
}

// Drain is a Consumer that drains any connected providers.
type Drain interface {
	Consumer
	Run()
}

// Filter is an interface for anything that both consumes and provides audio data.
type Filter interface {
	Provider
	Consumer
}
