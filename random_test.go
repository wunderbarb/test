// V 0.6.2
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Sep 2020

package test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_test_RandomString(t *testing.T) {
	Describe(t)
	for i := 1; i < 50; i++ {
		s := RandomString(i)
		assert.Equal(t, i, len(s), "failed to generate a string of %d characters", i)
	}

	s := RandomString(0)
	assert.NotEqual(t, 0, len(s), "failed to generate with 0 size: %d", len(s))
}

func TestGenerateAlphaRandomString(t *testing.T) {
	Describe(t)
	s := RandomAlphaString(254, All)
	assert.Equal(t, 254, len(s))
	s = RandomAlphaString(254, AlphaNum)
	assert.NotContains(t, s, "@.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, AlphaNumNoSpace)
	assert.NotContains(t, s, " @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, Alpha)
	assert.NotContains(t, s, "01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, Caps)
	assert.NotContains(t, s, "abcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, Small)
	assert.NotContains(t, s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%")
	s = RandomAlphaString(254, 1000)
	assert.Equal(t, "", s)
}

func Test_test_GenerateName(t *testing.T) {
	Describe(t)
	for i := 1; i < 50; i++ {
		s := RandomName(i)
		assert.Equal(t, i, len(s), "failed to generate a string of %d characters", i)
	}

	s := RandomName(0)
	assert.NotEqual(t, 0, len(s), "failed to generate with 0 size: %d", len(s))
}

func Test_test_RandomSlice(t *testing.T) {
	Describe(t)
	n := rand.Intn(255) + 1
	d := RandomSlice(n)
	assert.Equal(t, n, len(d))

	d1 := RandomSlice(0)
	assert.NotZero(t, len(d1))
}

func Test_test_RandomID(t *testing.T) {
	Describe(t)
	id := RandomID()
	assert.Equal(t, 16, len(id))
	assert.NotContains(t, id, " @.!$&_+-:;*?#/\\,()[]{}<>%")
}

func Test_RandomCSVFile(t *testing.T) {
	Describe(t)

	name := "testdata/" + RandomID() + ".csv"

	err := RandomCSVFile(name, Rng.Intn(10)+1, Rng.Intn(10)+1, ';')
	require.NoError(t, err)
	assert.FileExists(t, name)
	os.Remove(name)

	require.Error(t, RandomCSVFile("bad/"+name, 10, 10, ';'))

}

func Test_RandomFile(t *testing.T) {
	Describe(t)

	n, err := RandomFile(3, "essai", true)
	require.NoError(t, err)
	assert.FileExists(t, n)
	os.Remove(n)
}

func Test_RandomFileWithDir(t *testing.T) {
	Describe(t)

	n, err := RandomFileWithDir(3, "essai", "")
	require.NoError(t, err)
	assert.FileExists(t, n)
	os.Remove(n)

	os.MkdirAll("testdata/essai", 0777)
	n, err = RandomFileWithDir(3, "essai", "testdata/essai")
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join("testdata/essai", n))
	os.Remove(filepath.Join("testdata/essai", n))

	_, err = RandomFileWithDir(3, "essai", "testdata/bad")
	assert.Error(t, err)
}
