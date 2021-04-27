# test
The package holds some functions I use a lot when unit testing.

It encapsulates the `github.com/stretch/testify` package in a neat way. 
```golang
func Test_InRAMReader_Read(t *testing.T) {
	require, assert := Describe(t)

	p := RandomSlice(Rng.Intn(10000) + 256)
	rd := NewInRAMReader(p)
	b := make([]byte, 128)

	n, err := rd.Read(b)
	require.NoError(err)
	assert.Equal(128, n)
	assert.Equal(p[:128], b)

	rd.Close()
	_, err = rd.Read(b)
	assert.EqualError(err, ErrClosed.Error())

}
```
The function `Describe` returns a `require` and `assert`.  These functions are identical to the tetsify assertins and requirements but do not ask the `t` in the call.  They increase readibility of the code.  Furthermode, `Describe` prints on the terminal a test sequence number with the name of the tested function.  For instance, in this example:
```
> Test 1: InRAMReader_Read
```

The package provides:

- a set of functions, such as `RandomSlice`, `RandomId`, or `RandomText`, that generate random information.  The randomness is not cryptographically secure.
- a set of stuctures, `InRAMWriter` and `InRAMReader` that emulates `RAM` buffers.
- a faulty `FaultyReader` that is an `io.Reader` which fails.
- and many other functions useful for unit testing.

`bench` is not ready for use.
