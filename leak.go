// v0.1.0
// Author: DIEHL E.
// © Jul 2024

package test

import (
	"strings"
	"testing"

	"github.com/ysmood/gotrace"
)

// NoLeak verifies whether there was no new running go routines after the current test.
func NoLeak(t *testing.T) {
	ign := gotrace.IgnoreCurrent()
	gotrace.CheckLeak(t, 0, ign)
}

// NoLeakButPersistentHTTP verifies whether there was no new running go routines after the current test excepted
// the ones that HTTP may create for a persistent connection
func NoLeakButPersistentHTTP(t *testing.T) {
	ign := gotrace.CombineIgnores(
		gotrace.IgnoreCurrent(),
		func(t *gotrace.Trace) bool {
			return strings.Contains(t.Raw, "net/http.(*persistConn).writeLoop")
		},
		gotrace.IgnoreFuncs("internal/poll.runtime_pollWait"),
	)
	gotrace.CheckLeak(t, 0, ign)
}
