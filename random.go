// V0.9.1a
// Author: Diehl E.
// (C) Sony Pictures Entertainment, Sep 2020

package test

import (
	"encoding/csv"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/wunderbarb/mypkg/toolbox"
)

// Rng is a randomly seeded random number generator that can be used for tests.
// The random number generator is not cryptographically safe.
var Rng *rand.Rand

// init initializes the random number generator.
func init() {
	Rng = rand.New(rand.NewSource(time.Now().UnixNano()))

}

// AlphaNumType represents the kind of characters
// that will be generated.
type AlphaNumType int

const (
	// All requests all the characters from the character set
	// ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%"
	All AlphaNumType = iota
	// AllCVS requests the same characters than 'All' at the exception of ; and ,.  It is to be
	// used when appluing to CVS files.
	AllCVS
	// AlphaNum requests only characters that are alpha-numerical with space included.
	AlphaNum
	// AlphaNumNoSpace requests only characters thaat are alpha-numerical without space.
	AlphaNumNoSpace
	// Alpha requests only characters that are alphabetical with space included.
	Alpha
	// AlphaNoSpace requests only characters thaat are alphabetical without space.
	AlphaNoSpace
	// Caps requests only upper characters without space.
	Caps
	// Small requests only minor characters without space.
	Small
)

// RandomID returns a random 16 character, alphanumeric, ID.
func RandomID() string {
	return RandomAlphaString(16, AlphaNumNoSpace)
}

// RandomName returns a random string with size characters.
// If size is null, then the length of the string is random in the range
// 1 to 256 characters.
//
// The character set is ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// mypckg/cryptobox package with GenerateNewKey
func RandomName(size int) string {

	return RandomAlphaString(size, AlphaNoSpace)
}

// RandomSlice returns a random slice with size bytes.
// If size is null, then the number of bytes in the slice is random in the range
// 1 to 256 characters.
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// mypckg/cryptobox package with GenerateNewKey
func RandomSlice(size int) []byte {
	if size == 0 {
		size = Rng.Intn(256) + 1
	}
	var buffer []byte
	for i := 0; i < size; i++ {
		buffer = append(buffer, byte(Rng.Intn(255)))
	}
	return buffer
}

// RandomString returns a random string with size characters.
// If size is null, then the length of the string is random in the range
// 1 to 256 characters.
//
// The character set is ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// mypckg/cryptobox package with GenerateNewKey
func RandomString(size int) string {

	return RandomAlphaString(size, All)
}

// RandomAlphaString generates a size-character random string which character
// set depends on the value of t.  if t is not a prpoer value, the returned value
// is the empty string.
// If size is null, then the length of the string is random in the range
// 1 to 256 characters.
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// mypkg/cryptobox package with GenerateNewKey
func RandomAlphaString(size int, t AlphaNumType) string {
	conv := map[AlphaNumType][]byte{
		All:             []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%\""),
		AllCVS:          []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:*?#/\\()[]{}<>%\""),
		AlphaNum:        []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890 "),
		AlphaNumNoSpace: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890"),
		Alpha:           []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz "),
		AlphaNoSpace:    []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
		Caps:            []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		Small:           []byte("abcdefghijklmnopqrstuvwxyz"),
	}
	if size == 0 {
		size = Rng.Intn(256) + 1
	}

	var buffer []byte
	choice, ok := conv[t]
	if !ok {
		return ""
	}

	choiceSizee := len(choice)
	for i := 0; i < size; i++ {
		// generates the characters
		s := Rng.Intn(choiceSizee)
		buffer = append(buffer, choice[s])
	}
	return string(buffer)
}

// RandomCSVFile generates a file `name` that is a CSV table
// of `columns` x `rows` using as separator `sep` with
// random size fields.
func RandomCSVFile(name string, columns int, rows int, sep rune) error {
	name = toolbox.SetExtension(name, "csv")
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	wr.Comma = sep
	for i := 0; i < rows; i++ {
		var rec []string
		for j := 0; j < columns; j++ {
			// uses a complete character set without potential delimiters
			rec = append(rec, RandomAlphaString(0, AllCVS))
		}
		err = wr.Write(rec)
		if err != nil {
			return err
		}
	}
	wr.Flush() // do not forget to flush :(,
	return nil
}

// RandomFile generates a random binary file of `size` K bytes with
// a random name and extension `ext`.  If `inTestdata` is true, the file
// is in "testdata/" subdirectory, else in the current directory.  It
// returns the name of the generated file.
//
// DEPRECATED:  should be replaced by RandomFileWithDir.
func RandomFile(size int, ext string, inTestdata bool) (string, error) {
	name := RandomID() + "." + ext
	if inTestdata {
		name = "testdata/" + name
	}
	f, err := os.Create(name)
	if err != nil {
		return "", nil
	}
	defer f.Close()

	for i := 0; i < size; i++ {
		_, err = f.Write(RandomSlice(1024))
		if err != nil {
			return "", err
		}
	}
	return f.Name(), nil
}

// RandomFileWithDir generates a random binary file of `size` K bytes with
// a random name and extension `ext`.  `path` defines the relative path of
// the location where to store the generated file. If empty string,
// it stores locally.
//
// It returns the name of the generated file (without) the path.
func RandomFileWithDir(size int, ext string, path string) (string, error) {
	name := toolbox.SetExtension(RandomID(), ext)
	if path != "" {
		name = filepath.Join(path, name)
	}
	f, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for i := 0; i < size; i++ {
		_, err = f.Write(RandomSlice(1024))
		if err != nil {
			return "", err
		}
	}
	return filepath.Base(f.Name()), nil
}
