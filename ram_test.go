// V0.7.4
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Feb 2021

package test

import (
	"io"
	"testing"
)

func Test_RamWriter(t *testing.T) {
	require, assert := Describe(t)
	text := RandomString(0)

	ramw := NewRAMWriter()
	bio := ramw.Writer()
	_, err := bio.WriteString(text)
	require.NoError(err)
	assert.Equal(text, ramw.AsString())
	assert.Equal([]byte(text), ramw.AsBytes())
}

func Test_RamReader(t *testing.T) {
	require, assert := Describe(t)
	text := RandomString(0) + "\n"

	ramr := NewRAMReader([]byte(text))
	bio := ramr.Reader()

	data, err := bio.ReadString('\n')
	require.NoError(err)
	assert.Equal(text, data)

}

func Test_RamReaderFromString(t *testing.T) {
	require, assert := Describe(t)
	text := RandomString(0) + "\n"

	ramr := NewRAMReaderFromString(text)
	bio := ramr.Reader()

	data, err := bio.ReadString('\n')
	require.NoError(err)
	assert.Equal(text, data)

}

func Test_ErrReader(t *testing.T) {
	_, assert := Describe(t)

	var er ErrReader
	assert.Error(fake(er))
}

func fake(rd io.Reader) error {
	var p []byte
	_, err := rd.Read(p)
	return err
}
