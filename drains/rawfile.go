package drains

import (
	"bufio"
	"encoding/binary"
	"os"

	"github.com/targodan/goplug"
)

type RawFileDrain struct {
	Drain
	file       *os.File
	buf        *bufio.Writer
	run        bool
	hasStopped chan bool
}

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

func (wf *RawFileDrain) Close() {
	wf.buf.Flush()
	wf.file.Close()
}

func (wf *RawFileDrain) Start() {
	wf.run = true
	for wf.run {
		binary.Write(wf.buf, binary.LittleEndian, wf.ihs.GetSocket(0).Read())
	}
	wf.hasStopped <- true
}

func (wf *RawFileDrain) Stop() {
	wf.hasStopped = make(chan bool)
	wf.run = false
	<-wf.hasStopped
}
