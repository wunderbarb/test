// V0.7.3
// Author: DIEHL E.
// © Nov 2024

package test

import (
	"testing"
)

func Test_CompareFiles(t *testing.T) {
	_, assert := Describe(t)

	assert.True(CompareFiles("4test.go", "4test.go"))
	assert.False(CompareFiles("4test.go", "random.go"))
}
