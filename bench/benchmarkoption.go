// v0.2.0
// Author: DIEHL E.
// © Sony Pictures Entertainment, Feb 2021

package bench

//BenchmarkOption allows to parameterize Benchmark function.
type BenchmarkOption func(opts *benchmarkOptions)

type benchmarkOptions struct {
	atLeast int
}

// WithAtLeast informs that the benchmark must do at least `n` iterations.
func WithAtLeast(n int) BenchmarkOption {
	return func(bo *benchmarkOptions) {
		bo.atLeast = n
	}
}

func collectOptions(options ...BenchmarkOption) *benchmarkOptions {
	opts := &benchmarkOptions{atLeast: 0}
	for _, option := range options {
		option(opts)
	}
	return opts
}
