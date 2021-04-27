// v0.3.2
// Author: DIEHL E.
// © Sony Pictures Entertainment, Apr 2021

package test

import (
	"testing"
)

func Test_InRAMReader_Read(t *testing.T) {
	require, assert := Describe(t)

	p := RandomSlice(Rng.Intn(10000) + 256)
	rd := NewInRAMReader(p)
	b := make([]byte, 128)

	n, err := rd.Read(b)
	require.NoError(err)
	assert.Equal(128, n)
	assert.Equal(p[:128], b)

	rd.Close()
	_, err = rd.Read(b)
	assert.EqualError(err, ErrClosed.Error())

}

func Test_InRAMWriter_Write(t *testing.T) {
	require, assert := Describe(t)

	wr := NewInRAMWriter()
	defer wr.Close()
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
	wr.Write(buffer)

	buf2 := RandomSlice(128)
	off := Rng.Int63n(64)
	n, err := wr.WriteAt(buf2, off)
	require.NoError(err)
	assert.Equal(128, n)
	assert.Equal(buf2, wr.Bytes()[off:128+off])

}
