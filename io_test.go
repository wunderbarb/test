// v0.2.0
// (C) Sony Pictures Entertainment, Jan 2021

package test

import "testing"

func Test_FaultyReader_Read(t *testing.T) {
	_, assert := Describe(t)

	var fr FaultyReader
	_, err := fr.Read(RandomSlice(10))
	assert.Error(err)
}

func Test_FaultyReader_Close(t *testing.T) {
	_, assert := Describe(t)

	var fr FaultyReader

	assert.Error(fr.Close())
}
