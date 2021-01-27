// v0.2.3
// Author: DIEHL E.
// (C) Sony Pictures Entertainment, Jan 2021

package test

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	limitBench = 100
)

var (
	// ErrExceedLimit occurs when the iteration exceeds 100.
	ErrExceedLimit = errors.New("exceed the limit")
)

// FuncToBench is the signature of the functions that can be benchmarked with
// function Bench.
type FuncToBench func() error

// BenchResult is the information returned by Bench on a succesful benchmark
type BenchResult struct {
	N     int
	Speed time.Duration
}

func (br BenchResult) String() string {
	return fmt.Sprintf(
		"%d iterations %f sec/ops", br.N, br.Speed.Seconds())
}

// Bench benchmarks the function `f`.  It iteraates the function
// until it reaches a variation of the avreage of less than `precision` percents.
// If the number of iterations exceeds 100, then it returns an error.
func Bench(precision int, f FuncToBench) (*BenchResult, error) {
	br := &BenchResult{N: 0}

	startTime := time.Now()
	for i := 0; i < 5; i++ {
		br.N++
		err := f()
		if err != nil {
			return br,
				errors.Wrapf(err, "benchmarked function failed at iteration %d", br.N)
		}
	}

	value := time.Since(startTime).Nanoseconds() / int64(br.N)
	for i := 5; i < limitBench; i++ {
		br.N++
		err := f()
		if err != nil {
			return br,
				errors.Wrapf(err, "benchmarked function failed at iteration %d", br.N)
		}

		newValue := time.Since(startTime).Nanoseconds() / int64(br.N)

		prec := abs((100 * (newValue - value)) / value)
		// fmt.Printf("%d p- %d\n", newValue, prec)
		if prec <= int64(precision) {
			br.Speed = time.Duration(newValue)
			return br, nil
		}
		value = newValue
	}
	br.Speed = time.Duration(value)
	return br, ErrExceedLimit
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
