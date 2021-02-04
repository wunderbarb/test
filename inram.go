// v0.2.1
// Author: DIEHL E.
// © Sony Pictures Entertainment, Feb 2021

package test

import (
	"bytes"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
)

var (
	// ErrClosed occurs when attempting to access a closed InRAMReader or InRAMWriter.
	ErrClosed = errors.New("reader or writer is closed")
)

// InRAMReader implements everything for io.Reader, io.Closer, and io.Seeker with an initial value but
// within RAM memory.
type InRAMReader struct {
	bytes.Reader
	closed bool
}

// InRAMWriter implemets an io.Writer, io.WriterAt, and io.Closer within RAM memory.
type InRAMWriter struct {
	aws.WriteAtBuffer
	closed     bool
	currentPos int64
}

// NewInRAMReader creates a new InRAMReader with the buffer `p`.
func NewInRAMReader(p []byte) *InRAMReader {
	irrd := &InRAMReader{
		Reader: *bytes.NewReader(p),
		closed: false,
	}
	return irrd
}

// NewInRAMReaderAsString creates a new InRAMReader with the string `s`.
func NewInRAMReaderAsString(s string) *InRAMReader {
	return NewInRAMReader([]byte(s))
}

// Close implements the io.Closer interface.
func (irr *InRAMReader) Close() error {
	if irr.closed {
		return ErrClosed
	}
	irr.closed = true
	return nil
}

// Read implements the io.Reader interface.
func (irr *InRAMReader) Read(p []byte) (int, error) {
	if irr.closed {
		return 0, ErrClosed
	}
	return irr.Reader.Read(p)
}

// Seek implements the io.Seeker interface.
func (irr *InRAMReader) Seek(offset int64, whence int) (int64, error) {
	if irr.closed {
		return 0, ErrClosed
	}
	return irr.Reader.Seek(offset, whence)
}

// NewInRAMWriter creates a new InRAMWriter.
func NewInRAMWriter() *InRAMWriter {
	buf := make([]byte, 10)

	return &InRAMWriter{
		WriteAtBuffer: *aws.NewWriteAtBuffer(buf),
		closed:        false,
		currentPos:    0,
	}
}

// Close implements the io.Closer interface.
func (irw *InRAMWriter) Close() error {
	if irw.closed {
		return ErrClosed
	}
	irw.closed = true
	return nil
}

func (irw *InRAMWriter) String() string {
	return string(irw.WriteAtBuffer.Bytes())
}

// Write implements the io.Writer interface.
func (irw *InRAMWriter) Write(p []byte) (int, error) {
	if irw.closed {
		return 0, ErrClosed
	}
	n, err := irw.WriteAtBuffer.WriteAt(p, irw.currentPos)
	irw.currentPos += int64(n)
	return n, err
}

// WriteAt implements the io.WriteAt interface.
func (irw *InRAMWriter) WriteAt(p []byte, pos int64) (int, error) {
	if irw.closed {
		return 0, ErrClosed
	}
	return irw.WriteAtBuffer.WriteAt(p, pos)
}
