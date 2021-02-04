// V0.7.2
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Feb 2021

package test

import (
	"testing"
)

func Test_CLI(t *testing.T) {
	_, assert := Describe(t)

	_, err := CLI(false, "ls")
	assert.NoError(err)
	_, err = CLI(false, "ls", "~/Dev1")
	assert.Error(err)
}

func Test_CompareFiles(t *testing.T) {
	_, assert := Describe(t)

	assert.True(CompareFiles("4test.go", "4test.go"))
	assert.False(CompareFiles("4test.go", "random.go"))
}
