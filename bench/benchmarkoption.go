// v0.2.0
// Author: DIEHL E.
// © Sony Pictures Entertainment, Feb 2021

package bench

import "encoding/csv"

//BenchmarkOption allows to parameterize Benchmark function.
type BenchmarkOption func(opts *benchmarkOptions)

type benchmarkOptions struct {
	atLeast   int
	csv       *csv.Writer // pointer to a csv writer.   Initialized via WithCSV
	benchName string
}

// WithAtLeast informs that the benchmark must do at least `n` iterations.
func WithAtLeast(n int) BenchmarkOption {
	return func(bo *benchmarkOptions) {
		bo.atLeast = n
	}
}

// WithCSV informs the benchmark that it must write the results in the csv writer.
func WithCSV(wr *csv.Writer) BenchmarkOption {
	return func(bo *benchmarkOptions) {
		bo.csv = wr
	}
}

func WithBenchName(name string) BenchmarkOption {
	return func(bo *benchmarkOptions) {
		bo.benchName = name
	}
}

func collectOptions(options ...BenchmarkOption) *benchmarkOptions {
	opts := &benchmarkOptions{atLeast: 0, csv: nil, benchName: ""}
	for _, option := range options {
		option(opts)
	}
	return opts
}
