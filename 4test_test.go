// V0.7.1
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Apr 2020

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func Test_CLI(t *testing.T) {
	Describe(t)

	_, err := CLI(false, "ls")
	assert.NoError(t, err)
	_, err = CLI(false, "ls", "~/Dev1")
	assert.Error(t, err)
}

func Test_CompareFiles(t *testing.T) {
	Describe(t)

	assert.True(t, CompareFiles("4test.go", "4test.go"))
	assert.False(t, CompareFiles("4test.go", "random.go"))
}
