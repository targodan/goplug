// Package goplug is an easy to use library for real-time manipulation of
// audio streams.
package goplug

// Sample holds the value and frequency information of a sample.
type Sample struct {
	Value           float32
	SampleFrequency uint
}
