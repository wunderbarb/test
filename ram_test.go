// V0.7.1
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Apr 2020

package test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RamWriter(t *testing.T) {
	Describe(t, "RamWriter")
	text := RandomString(0)

	ramw := NewRAMWriter()
	bio := ramw.Writer()
	_, err := bio.WriteString(text)
	require.NoError(t, err)
	assert.Equal(t, text, ramw.AsString())
	assert.Equal(t, []byte(text), ramw.AsBytes())
}

func Test_RamReader(t *testing.T) {
	Describe(t, "RAMReader")
	text := RandomString(0) + "\n"

	ramr := NewRAMReader([]byte(text))
	bio := ramr.Reader()

	data, err := bio.ReadString('\n')
	require.NoError(t, err)
	assert.Equal(t, text, data)

}

func Test_RamReaderFromString(t *testing.T) {
	Describe(t, "RAMReaderFromString")
	text := RandomString(0) + "\n"

	ramr := NewRAMReaderFromString(text)
	bio := ramr.Reader()

	data, err := bio.ReadString('\n')
	require.NoError(t, err)
	assert.Equal(t, text, data)

}

func Test_ErrReader(t *testing.T) {
	Describe(t)

	var er ErrReader
	assert.Error(t, fake(er))
}

func fake(rd io.Reader) error {
	var p []byte
	_, err := rd.Read(p)
	return err
}
