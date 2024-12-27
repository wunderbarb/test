// v0.2.2
// Author: wunderbarb
// © Sony Pictures Entertainment, Nov 2024

package test

import (
	"bytes"
	"errors"
	"sync"
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

// InRAMWriter implements an io.Writer, io.WriterAt, and io.Closer within RAM memory.
type InRAMWriter struct {
	writeAtBuffer
	closed     bool
	currentPos int64
}

// NewInRAMReader creates a new InRAMReader with the buffer `p`.
func NewInRAMReader(p []byte) *InRAMReader {
	return &InRAMReader{
		Reader: *bytes.NewReader(p),
		closed: false,
	}
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
		writeAtBuffer: *newWriteAtBuffer(buf),
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
	return string(irw.Bytes())
}

// Write implements the io.Writer interface.
func (irw *InRAMWriter) Write(p []byte) (int, error) {
	if irw.closed {
		return 0, ErrClosed
	}
	n, err := irw.WriteAt(p, irw.currentPos)
	irw.currentPos += int64(n)
	return n, err
}

// WriteAt implements the io.WriteAt interface.
func (irw *InRAMWriter) WriteAt(p []byte, pos int64) (int, error) {
	if irw.closed {
		return 0, ErrClosed
	}
	return irw.WriteAt(p, pos)
}

// A writeAtBuffer provides an in memory buffer supporting the io.WriterAt interface
type writeAtBuffer struct {
	buf []byte
	m   sync.Mutex

	// GrowthCoeff defines the growth rate of the internal buffer. By
	// default, the growth rate is 1, where expanding the internal
	// buffer will allocate only enough capacity to fit the new expected
	// length.
	GrowthCoeff float64
}

// newWriteAtBuffer creates a writeAtBuffer with an internal buffer
// provided by buf.
func newWriteAtBuffer(buf []byte) *writeAtBuffer {
	return &writeAtBuffer{buf: buf}
}

// WriteAt writes a slice of bytes to a buffer starting at the position provided
// The number of bytes written will be returned, or error. Can overwrite previous
// written slices if the write ats overlap.
func (b *writeAtBuffer) WriteAt(p []byte, pos int64) (n int, err error) {
	pLen := len(p)
	expLen := pos + int64(pLen)
	b.m.Lock()
	defer b.m.Unlock()
	if int64(len(b.buf)) < expLen {
		if int64(cap(b.buf)) < expLen {
			if b.GrowthCoeff < 1 {
				b.GrowthCoeff = 1
			}
			newBuf := make([]byte, expLen, int64(b.GrowthCoeff*float64(expLen)))
			copy(newBuf, b.buf)
			b.buf = newBuf
		}
		b.buf = b.buf[:expLen]
	}
	copy(b.buf[pos:], p)
	return pLen, nil
}

// Bytes returns a slice of bytes written to the buffer.
func (b *writeAtBuffer) Bytes() []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.buf
}
