// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Modified  by DIEHL E.
// v0.2.3

package bench

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

// BenchmarkResult contains the results of a benchmark run.
type BenchmarkResult struct {
	N         int           // The number of iterations.
	T         time.Duration // The total time taken.
	Bytes     int64         // Bytes processed in one iteration.
	MemAllocs uint64        // The total number of memory allocations.
	MemBytes  uint64        // The total number of bytes allocated.

	// Extra records additional metrics reported by ReportMetric.
	Extra map[string]float64
}

// NsPerOp returns the "ns/op" metric.
func (r BenchmarkResult) NsPerOp() int64 {
	if v, ok := r.Extra["ns/op"]; ok {
		return int64(v)
	}
	if r.N <= 0 {
		return 0
	}
	return r.T.Nanoseconds() / int64(r.N)
}

// MbPerSec returns the "MB/s" metric.
func (r BenchmarkResult) MbPerSec() float64 {
	if v, ok := r.Extra["MB/s"]; ok {
		return v
	}
	if r.Bytes <= 0 || r.T <= 0 || r.N <= 0 {
		return 0
	}
	return (float64(r.Bytes) * float64(r.N) / 1e6) / r.T.Seconds()
}

// AllocsPerOp returns the "allocs/op" metric,
// which is calculated as r.MemAllocs / r.N.
func (r BenchmarkResult) AllocsPerOp() int64 {
	if v, ok := r.Extra["allocs/op"]; ok {
		return int64(v)
	}
	if r.N <= 0 {
		return 0
	}
	return int64(r.MemAllocs) / int64(r.N)
}

// AllocedBytesPerOp returns the "B/op" metric,
// which is calculated as r.MemBytes / r.N.
func (r BenchmarkResult) AllocedBytesPerOp() int64 {
	if v, ok := r.Extra["B/op"]; ok {
		return int64(v)
	}
	if r.N <= 0 {
		return 0
	}
	return int64(r.MemBytes) / int64(r.N)
}

// CSV writes in the CSV file the follwoing record: `name`, number of iterations, number of ns per
// operation, number of allocated bytes per operation and number of Allocs per operation.
func (r BenchmarkResult) CSV(name string, wr *csv.Writer) error {
	var record []string
	record = append(record, name)
	record = append(record, fmt.Sprintf("%d", r.N))
	record = append(record, fmt.Sprintf("%d", r.NsPerOp()))
	record = append(record, fmt.Sprintf("%d", r.AllocedBytesPerOp()))
	record = append(record, fmt.Sprintf("%d", r.AllocsPerOp()))
	return wr.Write(record)
}

// String returns a summary of the benchmark results.
// It follows the benchmark result line format from
// https://golang.org/design/14313-benchmark-format, not including the
// benchmark name.
// Extra metrics override built-in metrics of the same name.
// String does not include allocs/op or B/op, since those are reported
// by MemString.
func (r BenchmarkResult) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "%8d\t", r.N)
	ns := r.NsPerOp()
	switch {
	case ns <= 1000:
		fmt.Fprintf(buf, "%d ns/op", ns)
	case ns <= 1000000:
		nn := float64(ns) / 1000.0
		fmt.Fprintf(buf, "%.3f us/op", nn)
	case ns <= 1000000000:
		nn := float64(ns) / 1000000.0
		fmt.Fprintf(buf, "%.3f ms/op", nn)
	default:
		nn := float64(ns) / 1000000000.0
		fmt.Fprintf(buf, "%.3f s/op", nn)
	}

	return buf.String()
}

func prettyPrint(w io.Writer, x float64, unit string) {
	// Print all numbers with 10 places before the decimal point
	// and small numbers with three sig figs.
	var format string
	switch y := math.Abs(x); {
	case y == 0 || y >= 99.95:
		format = "%10.0f %s"
	case y >= 9.995:
		format = "%12.1f %s"
	case y >= 0.9995:
		format = "%13.2f %s"
	case y >= 0.09995:
		format = "%14.3f %s"
	case y >= 0.009995:
		format = "%15.4f %s"
	case y >= 0.0009995:
		format = "%16.5f %s"
	default:
		format = "%17.6f %s"
	}
	fmt.Fprintf(w, format, x, unit)
}

// MemString returns r.AllocedBytesPerOp and r.AllocsPerOp in the same format as 'go test'.
func (r BenchmarkResult) MemString() string {
	return fmt.Sprintf("%8d B/op\t%8d allocs/op",
		r.AllocedBytesPerOp(), r.AllocsPerOp())
}

// benchmarkName returns full name of benchmark including procs suffix.
func benchmarkName(name string, n int) string {
	if n != 1 {
		return fmt.Sprintf("%s-%d", name, n)
	}
	return name
}
