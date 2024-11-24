// V 0.9.1
// Author: DIEHL E.
// (C)  Jul 2021

package test

import (
	"bufio"
	"bytes"
	"errors"
)

// RAMWriter is an emulation of bufio writer stored in RAM.
// It is mainly for test purpose.
type RAMWriter struct {
	b   *bytes.Buffer
	bio *bufio.Writer
}

// NewRAMWriter initializes a new RAMWriter.
func NewRAMWriter() *RAMWriter {
	var ramw RAMWriter
	var buf []byte
	ramw.b = bytes.NewBuffer(buf)
	ramw.bio = bufio.NewWriter(ramw.b)
	return &ramw
}

// AsBytes returns the value of the RAMWriter as a slice.
func (ramw *RAMWriter) AsBytes() []byte {
	ramw.Flush()
	return ramw.b.Bytes()
}

// AsString returns the value of the RAMWriter as a string.
func (ramw *RAMWriter) AsString() string {
	ramw.Flush()
	return ramw.b.String()
}

// Flush flushes the writer
func (ramw *RAMWriter) Flush() error {
	return ramw.Writer().Flush()
}

// Write the slice p to the writer.  It retruns the
// number of written bytes.  It implements the io.Writer
// interface.
func (ramw *RAMWriter) Write(p []byte) (int, error) {
	return ramw.b.Write(p)
}

// Writer returns the pointer to the bufio.Writer.
func (ramw *RAMWriter) Writer() *bufio.Writer {

	return ramw.bio
}

// RAMReader is an emulation of bufio.Reader  and io.ReadCloser
// stored in RAM. It is mainly for test purpose.
type RAMReader struct {
	b   *bytes.Reader
	bio *bufio.Reader
}

// NewRAMReader initializes a new RAMReader with the information in `data`.
func NewRAMReader(data []byte) *RAMReader {
	var ramr RAMReader
	ramr.b = bytes.NewReader(data)
	ramr.bio = bufio.NewReader(ramr.b)
	return &ramr
}

// NewRAMReaderFromString initializes a new RAMReader with
// the string `m`
func NewRAMReaderFromString(m string) *RAMReader {
	var ramr RAMReader
	ramr.b = bytes.NewReader([]byte(m))
	ramr.bio = bufio.NewReader(ramr.b)
	return &ramr
}

// Reader returns the pointer to the bufio.Reader.
func (ramr *RAMReader) Reader() *bufio.Reader {
	return ramr.bio
}

// Close implements io.Closer interface
func (ramr *RAMReader) Close() error {
	return nil
}

// Read is the io.Reader interface implementation.
func (ramr *RAMReader) Read(b []byte) (n int, err error) {
	return ramr.b.Read(b)
}

// Seek implements the io.Seeker interface.
func (ramr *RAMReader) Seek(offset int64, whence int) (int64, error) {
	return ramr.b.Seek(offset, whence)

}

// ErrReader is an io.Reader that returns an error
// when called.
// Deprecated: replaced by FaultyReader
type ErrReader struct{}

func (er ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("this is an error")
}
