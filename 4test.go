// V 0.9.0
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Jan 2021

package test

import (
	"fmt"
	"os/exec"
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

// CLI executes the command `app` with the parameters `params`.
// `expErr` indicates whether an error is expected.  If the expectation
// is not met, then it outputs the full CLI command that failed.
//
// DEPRECATED CLI should be repalced by githyb.com/wunderbarb/syst.Run using
// option WithVerboseTest.
func CLI(expErr bool, app string, params ...string) ([]byte, error) {
	cmd := exec.Command(app, params...)
	answer, err := cmd.Output()
	if (err == nil) == expErr {
		s := app
		for _, p := range params {
			s += " " + p
		}
		fmt.Printf("failed on: %s\n", s)
	}
	return answer, err
}

// CompareFiles returns true if files f1 and f2 are identical
func CompareFiles(f1 string, f2 string) bool {
	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
	equal, _ := cmp.CompareFile(f1, f2)
	return equal
}

// Describe displays the order of the test, the name of the function
//  and its optional description provided by 'msg'.  It initializes an assert
// and a require and returns them.
func Describe(t *testing.T, msg ...string) (*require.Assertions,
	*assert.Assertions) {

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
//  and its optional description provided by 'msg'.
func Describeb(b *testing.B, msg ...string) {

	dispMsg := ""
	if len(msg) != 0 {
		dispMsg = msg[0]
	}
	name := strings.TrimPrefix(b.Name(), "Bench_")
	fmt.Printf("Bench %d: %s %s\n", testCounter, name, dispMsg)
	testCounter++
}
