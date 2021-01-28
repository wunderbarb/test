// v0.1.0

package test

import "testing"

func Test_FaultyReader_Read(t *testing.T) {
	_, assert := Describe(t)

	var fr FaultyReader
	_, err := fr.Read(RandomSlice(10))
	assert.Error(err)
}
