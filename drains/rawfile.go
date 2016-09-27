package drains

import (
	"bufio"
	"encoding/binary"
	"os"

	"github.com/targodan/goplug"
)

// RawFileDrain is a goplug.Drain, that writes the samples it receives to a file.
// No conversion is done, the float32 are just written byte per byte in little endian order.
type RawFileDrain struct {
	Drain
	file       *os.File
	buf        *bufio.Writer
	run        bool
	hasStopped chan bool
}

// NewRawFileDrain creates a new RawFileDrain opening the file with the given filename.
func NewRawFileDrain(filename string) (*RawFileDrain, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewWriter(f)
	ret := &RawFileDrain{
		file: f,
		buf:  buf,
		run:  true,
	}
	ret.ihs = goplug.NewInputSocketHandler(1)
	return ret, nil
}

// Close flushes the buffers and closes any file descriptors.
func (wf *RawFileDrain) Close() {
	wf.buf.Flush()
	wf.file.Close()
}

// Start starts the drain process. This should be called as a goroutine.
func (wf *RawFileDrain) Start() {
	wf.run = true
	for wf.run {
		binary.Write(wf.buf, binary.LittleEndian, wf.ihs.GetSocket(0).Read())
	}
	wf.hasStopped <- true
}

// Stop tells the goroutine running the Start Method to stop and blocks until
// it actually stopped.
func (wf *RawFileDrain) Stop() {
	wf.hasStopped = make(chan bool)
	wf.run = false
	<-wf.hasStopped
}
