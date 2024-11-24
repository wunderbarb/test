// v0.3.3
// Author: DIEHL E.
// ©  Nov 2024

package test

import (
	"math/rand/v2"
	"testing"
)

func Test_InRAMReader_Read(t *testing.T) {
	require, assert := Describe(t)

	p := RandomSlice(rand.IntN(10000) + 256)
	rd := NewInRAMReader(p)
	b := make([]byte, 128)

	n, err := rd.Read(b)
	require.NoError(err)
	assert.Equal(128, n)
	assert.Equal(p[:128], b)

	_ = rd.Close()
	_, err = rd.Read(b)
	assert.EqualError(err, ErrClosed.Error())

}

func Test_InRAMWriter_Write(t *testing.T) {
	require, assert := Describe(t)

	wr := NewInRAMWriter()
	defer func() { _ = wr.Close() }()
	buffer := RandomSlice(256)

	n, err := wr.Write(buffer[:128])
	require.NoError(err)
	assert.Equal(128, n)
	n, err = wr.Write(buffer[128:])
	require.NoError(err)
	assert.Equal(128, n)
	assert.Equal(buffer, wr.Bytes())
}

func Test_InRAMWriter_WriteAt(t *testing.T) {
	require, assert := Describe(t)

	wr := NewInRAMWriter()
	buffer := RandomSlice(256)
	_, _ = wr.Write(buffer)

	buf2 := RandomSlice(128)
	off := rand.Int64N(64)
	n, err := wr.WriteAt(buf2, off)
	require.NoError(err)
	assert.Equal(128, n)
	assert.Equal(buf2, wr.Bytes()[off:128+off])

}
