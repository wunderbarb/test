// v0.3.1
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

// BenchOption allows to parameterize the Bench function.
type BenchOption func(opts *benchOptions)

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
func Bench(precision int, f FuncToBench, options ...BenchOption) (*BenchResult, error) {
	opts := benchOptions{fVerboseIteration: false}
	for _, option := range options {
		option(&opts)
	}

	br := &BenchResult{N: 0}

	startTime := time.Now()
	for i := 0; i < 5; i++ {
		br.N++
		err := f()
		if err != nil {
			return br,
				errors.Wrapf(err, "benchmarked function failed at iteration %d", br.N)
		}
		if opts.fVerboseIteration {
			fmt.Printf("Iter %d - mS %f\n", br.N,
				float64(average(startTime, br.N))/1000)
		}
	}

	value := average(startTime, br.N)
	for i := 5; i < limitBench; i++ {
		br.N++
		err := f()
		if err != nil {
			return br,
				errors.Wrapf(err, "benchmarked function failed at iteration %d", br.N)
		}

		newValue := average(startTime, br.N)
		if opts.fVerboseIteration {
			fmt.Printf("Iter %d - mS %f\n", br.N, float64(newValue)/1000)
		}

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

// WithVerboseIteration requests to print on stdout the result for each iteration.
func WithVerboseIteration() BenchOption {
	return func(b *benchOptions) {
		b.fVerboseIteration = true
	}
}

// ===========================================
// ===========================================

type benchOptions struct {
	fVerboseIteration bool
}

// -------------------------------------------------------
// --------------------------------------------------------

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// calculate the average tsime since `t` for `iter` values.
func average(t time.Time, iter int) int64 {
	return time.Since(t).Nanoseconds() / int64(iter)
}
