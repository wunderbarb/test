// v0.9.1
// Author: DIEHL E.
// © Sony Pictures Entertainment, Nov 2024

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

var (
	// testCounter is a counter used with function Describe
	testCounter = 1
)

// CompareFiles returns true if files f1 and f2 are identical
func CompareFiles(f1 string, f2 string) bool {
	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
	equal, _ := cmp.CompareFile(f1, f2)
	return equal
}

// Describe displays the order of the test, the name of the function and its optional description provided by 'msg'.
// It initializes an assert and require and returns them.
func Describe(t *testing.T, msg ...string) (*require.Assertions, *assert.Assertions) {

	dispMsg := ""
	if len(msg) != 0 {
		dispMsg = msg[0]
	}
	name := strings.TrimPrefix(t.Name(), "Test_")
	fmt.Printf("Test %d: %s %s\n", testCounter, name, dispMsg)
	testCounter++
	return require.New(t), assert.New(t)
}

// Describeb displays the order of the test, the name of the function
//
//	and its optional description provided by 'msg'.
func Describeb(b *testing.B, msg ...string) {

	dspMsg := ""
	if len(msg) != 0 {
		dspMsg = msg[0]
	}
	name := strings.TrimPrefix(b.Name(), "Bench_")
	fmt.Printf("Bench %d: %s %s\n", testCounter, name, dspMsg)
	testCounter++
}
