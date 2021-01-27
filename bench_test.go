// v0.1.2
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Jan 2021

package test

import (
	"fmt"
	"testing"
)

func Test_Bench(t *testing.T) {
	require, assert := Describe(t)

	br, err := Bench(3, benched1)
	fmt.Println(br.String())
	require.NoError(err)
	br1, err := Bench(5, benched2)
	fmt.Println(br1.String())
	require.NoError(err)
	assert.True(br.Speed < br1.Speed)
}

func benched1() error {
	RandomSlice(1024*1024 + Rng.Intn(500)*1024)
	return nil
}

func benched2() error {
	RandomSlice(10 * 1024 * 1024)
	return nil
}
