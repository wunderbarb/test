// V0.10.0
// Author: Diehl E.
// © Nov 2024

package test

import (
	rand1 "crypto/rand"
	"encoding/csv"
	"math/rand/v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// AlphaNumType represents the kind of characters
// that will be generated.
type AlphaNumType int

const (
	// All requests all the characters from the character set
	// ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%"
	All AlphaNumType = iota
	// AllCVS requests the same characters as 'All' at the exception to ; and ,.  It is to be
	// used when applying to CVS files.
	AllCVS
	// AlphaNum requests only characters that are alphanumerical with space included.
	AlphaNum
	// AlphaNumNoSpace requests only characters that are alphanumerical without space.
	AlphaNumNoSpace
	// Alpha requests only characters that are alphabetical with space included.
	Alpha
	// AlphaNoSpace requests only characters that are alphabetical without space.
	AlphaNoSpace
	// Caps requests only upper characters without space.
	Caps
	// Small requests only minor characters without space.
	Small
)

// RandomID returns a random 16-character, alphanumeric, ID.
func RandomID() string {
	const sizeID = 16
	return RandomAlphaString(sizeID, AlphaNumNoSpace)
}

// RandomName returns a random string with size characters.
// If size is null, then the length of the string is random in the range
// 1 to 256 characters.
//
// # The character set is ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating passphrases.
func RandomName(size int) string {

	return RandomAlphaString(size, AlphaNoSpace)
}

// RandomSlice returns a random slice with size bytes.
// If size is zero or negative, then the number of bytes in the slice is random in the range
// 1 to 256 characters.
func RandomSlice(size int) []byte {
	const size0 = 256 // max number of bytes for random set.
	if size <= 0 {
		size = rand.IntN(size0) + 1
	}
	buffer := make([]byte, size)
	_, _ = rand1.Read(buffer)
	return buffer
}

// RandomString returns a random string with size characters.
// If size is null, then the length of the string is random in the range
// 1 to 256 characters.
//
// The character set is ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890 @.!$&_+-:;*?#/\\,()[]{}<>%
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.
func RandomString(size int) string {

	return RandomAlphaString(size, All)
}

// RandomAlphaString generates a size-character random string which character
// set depends on the value of t.  if t is not a proper value, the returned value
// is the empty string.
// If size is zero or negative, then the length of the string is random in the range
// 1 to 256 characters.
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.
func RandomAlphaString(size int, t AlphaNumType) string {
	const size0 = 256 // max number of bytes for random set.
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
	if size <= 0 {
		size = rand.IntN(size0) + 1
	}
	var buffer []byte
	choice, ok := conv[t]
	if !ok {
		return ""
	}
	choiceSize := len(choice)
	for i := 0; i < size; i++ {
		// generates the characters
		s := rand.IntN(choiceSize)
		buffer = append(buffer, choice[s])
	}
	return string(buffer)
}

// RandomCSVFile generates a file `name` that is a CSV table
// of `columns` x `rows` using as separator `sep` with
// random size fields.
func RandomCSVFile(name string, columns int, rows int, sep rune) error {
	name = setExtension(name, "csv")
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
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

// RandomFileWithDir generates a random binary file of `size` K bytes with
// a random name and extension `ext`.  `path` defines the relative path of
// the location where to store the generated file. If empty string,
// it stores locally.
//
// It returns the name of the generated file (without) the path.
func RandomFileWithDir(size int, ext string, path string) (string, error) {
	const sizeOfSlices = 1024
	name := setExtension(RandomID(), ext)
	if path != "" {
		name = filepath.Join(path, name)
	}
	f, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	p := make([]byte, sizeOfSlices)
	for i := 0; i < size; i++ {
		_, err = rand1.Read(p)
		if err != nil {
			return "", err
		}
		_, err = f.Write(p)
		if err != nil {
			return "", err
		}
	}
	return filepath.Base(f.Name()), nil
}

// SwapCase randomly changes each character to upper or lower case
func SwapCase(s string) string {
	const dice = 3
	var sb strings.Builder
	for _, r := range s {
		switch rand.IntN(dice) { //nolint:gosec
		case 0:
			sb.WriteRune(toLower(r))
		case 1:
			sb.WriteRune(toUpper(r))
		case dice - 1:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// setExtension ensures that the extension ext is present at the end of the file.  If ext does not have a
// trailing '.', it adds the proper extension.
//
// Is the same as mypkg/setExtension but need to break cyclic mod.
func setExtension(name string, ext string) string {
	if ext == "" {
		return name
	}
	if ext[:1] != "." {
		ext = "." + ext
	}
	return strip(name, ext) + ext
}

// strip removes from `name` the extension `ext` if present.
func strip(name string, ext string) string {
	if ext == "" {
		return name
	}
	if ext[:1] != "." {
		ext = "." + ext
	}
	if path.Ext(name) == ext {

		s := name
		l := len(name) - len(ext)
		name = s[:l]
	}
	return name
}

// toUpper converts a rune to upper case if it's a letter
func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	}
	return r
}

// toLower converts a rune to lower case if it's a letter
func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r - 'A' + 'a'
	}
	return r
}
