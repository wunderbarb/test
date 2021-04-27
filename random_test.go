// V 0.7.1
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Apr 2021

package test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func Test_test_RandomString(t *testing.T) {
	_, assert := Describe(t)
	for i := 1; i < 50; i++ {
		s := RandomString(i)
		assert.Equal(i, len(s), "failed to generate a string of %d characters", i)
	}

	s := RandomString(0)
	assert.NotEqual(0, len(s), "failed to generate with 0 size: %d", len(s))
}

func TestGenerateAlphaRandomString(t *testing.T) {
	_, assert := Describe(t)
	s := RandomAlphaString(254, All)
	assert.Equal(254, len(s))
	s = RandomAlphaString(254, AlphaNum)
	assert.NotContains(s, "@.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, AlphaNumNoSpace)
	assert.NotContains(s, " @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, Alpha)
	assert.NotContains(s, "01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, Caps)
	assert.NotContains(s, "abcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, Small)
	assert.NotContains(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, 1000)
	assert.Equal("", s)
}

func Test_test_GenerateName(t *testing.T) {
	_, assert := Describe(t)
	for i := 1; i < 50; i++ {
		s := RandomName(i)
		assert.Equal(i, len(s), "failed to generate a string of %d characters", i)
	}

	s := RandomName(0)
	assert.NotEqual(0, len(s), "failed to generate with 0 size: %d", len(s))
}

func Test_test_RandomSlice(t *testing.T) {
	_, assert := Describe(t)
	n := rand.Intn(255) + 1
	d := RandomSlice(n)
	assert.Equal(n, len(d))

	d1 := RandomSlice(0)
	assert.NotZero(len(d1))
}

func Test_test_RandomID(t *testing.T) {
	_, assert := Describe(t)
	id := RandomID()
	assert.Equal(16, len(id))
	assert.NotContains(id, " @.!$&_+-:;*?#/\\,()[]{}<>%")
}

func Test_RandomCSVFile(t *testing.T) {
	require, assert := Describe(t)

	name := "testdata/" + RandomID() + ".csv"

	err := RandomCSVFile(name, Rng.Intn(10)+1, Rng.Intn(10)+1, ';')
	require.NoError(err)
	assert.FileExists(name)
	os.Remove(name)

	require.Error(RandomCSVFile("bad/"+name, 10, 10, ';'))

}

func Test_RandomFile(t *testing.T) {
	require, assert := Describe(t)

	n, err := RandomFile(3, "essai", true)
	require.NoError(err)
	assert.FileExists(n)
	os.Remove(n)
}

func Test_RandomFileWithDir(t *testing.T) {
	require, assert := Describe(t)

	n, err := RandomFileWithDir(3, "essai", "")
	require.NoError(err)
	assert.FileExists(n)
	os.Remove(n)

	os.MkdirAll("testdata/essai", 0777)
	n, err = RandomFileWithDir(3, "essai", "testdata/essai")
	require.NoError(err)
	assert.FileExists(filepath.Join("testdata/essai", n))
	os.Remove(filepath.Join("testdata/essai", n))

	_, err = RandomFileWithDir(3, "essai", "testdata/bad")
	assert.Error(err)
}

func Benchmark_RandomSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomSlice(1024)
	}

}

func Benchmark_RandomFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n, _ := RandomFileWithDir(10, "tst", "testdata/essai")
		os.Remove(filepath.Join("testdata/essai", n))
	}

}
