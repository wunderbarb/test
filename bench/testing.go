// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bench provides support for benchmarking of Go packages.
// It is intended to be used in concert with the "go test" command, which automates
// execution of any function of the form
//     func TestXxx(*testing.T)
// where Xxx does not start with a lowercase letter. The function name
// serves to identify the test routine.
//

// Benchmarks
//
// Functions of the form
//     func BenchmarkXxx(*testing.B)
// are considered benchmarks, and are executed by the "go test" command when
// its -bench flag is provided. Benchmarks are run sequentially.
//
// For a description of the testing flags, see
// https://golang.org/cmd/go/#hdr-Testing_flags
//
// A sample benchmark function looks like this:
//     func BenchmarkRandInt(b *testing.B) {
//         for i := 0; i < b.N; i++ {
//             rand.Int()
//         }
//     }
//
// The benchmark function must run the target code b.N times.
// During benchmark execution, b.N is adjusted until the benchmark function lasts
// long enough to be timed reliably. The output
//     BenchmarkRandInt-8   	68453040	        17.8 ns/op
// means that the loop ran 68453040 times at a speed of 17.8 ns per loop.
//
// If a benchmark needs some expensive setup before running, the timer
// may be reset:
//
//     func BenchmarkBigLen(b *testing.B) {
//         big := NewBig()
//         b.ResetTimer()
//         for i := 0; i < b.N; i++ {
//             big.Len()
//         }
//     }
//
// If a benchmark needs to test performance in a parallel setting, it may use
// the RunParallel helper function; such benchmarks are intended to be used with
// the go test -cpu flag:
//
//     func BenchmarkTemplateParallel(b *testing.B) {
//         templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
//         b.RunParallel(func(pb *testing.PB) {
//             var buf bytes.Buffer
//             for pb.Next() {
//                 buf.Reset()
//                 templ.Execute(&buf, "World")
//             }
//         })
//     }
//
// Examples
//
// The package also runs and verifies example code. Example functions may
// include a concluding line comment that begins with "Output:" and is compared with
// the standard output of the function when the tests are run. (The comparison
// ignores leading and trailing space.) These are examples of an example:
//
//     func ExampleHello() {
//         fmt.Println("hello")
//         // Output: hello
//     }
//
//     func ExampleSalutations() {
//         fmt.Println("hello, and")
//         fmt.Println("goodbye")
//         // Output:
//         // hello, and
//         // goodbye
//     }
//
// The comment prefix "Unordered output:" is like "Output:", but matches any
// line order:
//
//     func ExamplePerm() {
//         for _, value := range Perm(5) {
//             fmt.Println(value)
//         }
//         // Unordered output: 4
//         // 2
//         // 1
//         // 3
//         // 0
//     }
//
// Example functions without output comments are compiled but not executed.
//
// The naming convention to declare examples for the package, a function F, a type T and
// method M on type T are:
//
//     func Example() { ... }
//     func ExampleF() { ... }
//     func ExampleT() { ... }
//     func ExampleT_M() { ... }
//
// Multiple example functions for a package/type/function/method may be provided by
// appending a distinct suffix to the name. The suffix must start with a
// lower-case letter.
//
//     func Example_suffix() { ... }
//     func ExampleF_suffix() { ... }
//     func ExampleT_suffix() { ... }
//     func ExampleT_M_suffix() { ... }
//
// The entire test file is presented as the example when it contains a single
// example function, at least one other function, type, variable, or constant
// declaration, and no test or benchmark functions.
//
// Skipping
//
// Tests or benchmarks may be skipped at run time with a call to
// the Skip method of *T or *B:
//
//     func TestTimeConsuming(t *testing.T) {
//         if testing.Short() {
//             t.Skip("skipping test in short mode.")
//         }
//         ...
//     }
//
// Subtests and Sub-benchmarks
//
package bench

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/trace"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var initRan bool

// Init registers testing flags. These flags are automatically registered by
// the "go test" command before running test functions, so Init is only needed
// when calling functions such as Benchmark without using "go test".
//
// Init has no effect if it was already called.
func Init() {
	if initRan {
		return
	}
	initRan = true
	// The short flag requests that tests run more quickly, but its functionality
	// is provided by test writers themselves. The testing package is just its
	// home. The all.bash installation script sets it to make installation more
	// efficient, but by default the flag is off so a plain "go test" will do a
	// full test of the package.
	short = flag.Bool("test.short", false, "run smaller test suite to save time")

	// The failfast flag requests that test execution stop after the first test failure.
	failFast = flag.Bool("test.failfast", false, "do not start new tests after the first test failure")

	// The directory in which to create profile files and the like. When run from
	// "go test", the binary always runs in the source directory for the package;
	// this flag lets "go test" tell the binary to write the files in the directory where
	// the "go test" command is run.
	outputDir = flag.String("test.outputdir", "", "write profiles to `dir`")
	// Report as tests are run; default is silent for success.
	chatty = flag.Bool("test.v", false, "verbose: print additional output")
	count = flag.Uint("test.count", 1, "run tests and benchmarks `n` times")
	coverProfile = flag.String("test.coverprofile", "", "write a coverage profile to `file`")
	matchList = flag.String("test.list", "", "list tests, examples, and benchmarks matching `regexp` then exit")
	match = flag.String("test.run", "", "run only tests and examples matching `regexp`")
	memProfile = flag.String("test.memprofile", "", "write an allocation profile to `file`")
	memProfileRate = flag.Int("test.memprofilerate", 0, "set memory allocation profiling `rate` (see runtime.MemProfileRate)")
	cpuProfile = flag.String("test.cpuprofile", "", "write a cpu profile to `file`")
	blockProfile = flag.String("test.blockprofile", "", "write a goroutine blocking profile to `file`")
	blockProfileRate = flag.Int("test.blockprofilerate", 1, "set blocking profile `rate` (see runtime.SetBlockProfileRate)")
	mutexProfile = flag.String("test.mutexprofile", "", "write a mutex contention profile to the named file after execution")
	mutexProfileFraction = flag.Int("test.mutexprofilefraction", 1, "if >= 0, calls runtime.SetMutexProfileFraction()")
	traceFile = flag.String("test.trace", "", "write an execution trace to `file`")
	timeout = flag.Duration("test.timeout", 0, "panic test binary after duration `d` (default 0, timeout disabled)")
	cpuListStr = flag.String("test.cpu", "", "comma-separated `list` of cpu counts to run each test with")
	parallel = flag.Int("test.parallel", runtime.GOMAXPROCS(0), "run at most `n` tests in parallel")
	testlog = flag.String("test.testlogfile", "", "write test action log to `file` (for use only by cmd/go)")

	initBenchmarkFlags()
}

var (
	// Flags, registered during Init.
	short                *bool
	failFast             *bool
	outputDir            *string
	chatty               *bool
	count                *uint
	coverProfile         *string
	matchList            *string
	match                *string
	memProfile           *string
	memProfileRate       *int
	cpuProfile           *string
	blockProfile         *string
	blockProfileRate     *int
	mutexProfile         *string
	mutexProfileFraction *int
	traceFile            *string
	timeout              *time.Duration
	cpuListStr           *string
	parallel             *int
	testlog              *string

	cpuList     []int
	testlogFile *os.File

	// numFailed uint32 // number of test failures
)

type chattyPrinter struct {
	w          io.Writer
	lastNameMu sync.Mutex // guards lastName
	lastName   string     // last printed test name in chatty mode
}

func newChattyPrinter(w io.Writer) *chattyPrinter {
	return &chattyPrinter{w: w}
}

// Updatef prints a message about the status of the named test to w.
//
// The formatted message must include the test name itself.
func (p *chattyPrinter) Updatef(testName, format string, args ...interface{}) {
	p.lastNameMu.Lock()
	defer p.lastNameMu.Unlock()

	// Since the message already implies an association with a specific new test,
	// we don't need to check what the old test name was or log an extra CONT line
	// for it. (We're updating it anyway, and the current message already includes
	// the test name.)
	p.lastName = testName
	fmt.Fprintf(p.w, format, args...)
}

// Printf prints a message, generated by the named test, that does not
// necessarily mention that tests's name itself.
func (p *chattyPrinter) Printf(testName, format string, args ...interface{}) {
	p.lastNameMu.Lock()
	defer p.lastNameMu.Unlock()

	if p.lastName == "" {
		p.lastName = testName
	} else if p.lastName != testName {
		fmt.Fprintf(p.w, "=== CONT  %s\n", testName)
		p.lastName = testName
	}

	fmt.Fprintf(p.w, format, args...)
}

// The maximum number of stack frames to go through when skipping helper functions for
// the purpose of decorating log messages.
const maxStackLen = 50

// common holds the elements common between T and B and
// captures common methods such as Errorf.
type common struct {
	mu          sync.RWMutex        // guards this group of fields
	output      []byte              // Output generated by test or benchmark.
	w           io.Writer           // For flushToParent.
	ran         bool                // Test or benchmark (or one of its subtests) was executed.
	failed      bool                // Test or benchmark has failed.
	skipped     bool                // Test of benchmark has been skipped.
	done        bool                // Test is finished and all subtests have completed.
	helpers     map[string]struct{} // functions to be skipped when writing file/line info
	cleanup     func()              // optional function to be called at the end of the test
	cleanupName string              // Name of the cleanup function.
	cleanupPc   []uintptr           // The stack trace at the point where Cleanup was called.

	chatty   *chattyPrinter // A copy of chattyPrinter, if the chatty flag is set.
	bench    bool           // Whether the current test is a benchmark.
	finished bool           // Test function has completed.
	hasSub   int32          // Written atomically.
	// raceErrors int            // Number of races detected during test.
	runner string // Function name of tRunner running the test.

	parent   *common
	level    int       // Nesting depth of test or benchmark.
	creator  []uintptr // If level > 0, the stack trace at the point where the parent called t.Run.
	name     string    // Name of test or benchmark.
	start    time.Time // Time test or benchmark started
	duration time.Duration
	// barrier  chan bool // To signal parallel subtests they may start.
	signal chan bool // To signal a test is done.
	// sub      []*T      // Queue of subtests to be run in parallel.

	tempDirOnce sync.Once
	tempDir     string
	tempDirErr  error
	tempDirSeq  int32
}

// Short reports whether the -test.short flag is set.
func Short() bool {
	if short == nil {
		panic("testing: Short called before Init")
	}
	// Catch code that calls this from TestMain without first calling flag.Parse.
	if !flag.Parsed() {
		panic("testing: Short called before Parse")
	}

	return *short
}

// CoverMode reports what the test coverage mode is set to. The
// values are "set", "count", or "atomic". The return value will be
// empty if test coverage is not enabled.
func CoverMode() string {
	return cover.Mode
}

// Verbose reports whether the -test.v flag is set.
func Verbose() bool {
	// Same as in Short.
	if chatty == nil {
		panic("testing: Verbose called before Init")
	}
	if !flag.Parsed() {
		panic("testing: Verbose called before Parse")
	}
	return *chatty
}

// frameSkip searches, starting after skip frames, for the first caller frame
// in a function not marked as a helper and returns that frame.
// The search stops if it finds a tRunner function that
// was the entry point into the test and the test is not a subtest.
// This function must be called with c.mu held.
func (c *common) frameSkip(skip int) runtime.Frame {
	// If the search continues into the parent test, we'll have to hold
	// its mu temporarily. If we then return, we need to unlock it.
	shouldUnlock := false
	defer func() {
		if shouldUnlock {
			c.mu.Unlock()
		}
	}()
	var pc [maxStackLen]uintptr
	// Skip two extra frames to account for this function
	// and runtime.Callers itself.
	n := runtime.Callers(skip+2, pc[:])
	if n == 0 {
		panic("testing: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	var firstFrame, prevFrame, frame runtime.Frame
	for more := true; more; prevFrame = frame {
		frame, more = frames.Next()
		if frame.Function == c.cleanupName {
			frames = runtime.CallersFrames(c.cleanupPc)
			continue
		}
		if firstFrame.PC == 0 {
			firstFrame = frame
		}
		if frame.Function == c.runner {
			// We've gone up all the way to the tRunner calling
			// the test function (so the user must have
			// called tb.Helper from inside that test function).
			// If this is a top-level test, only skip up to the test function itself.
			// If we're in a subtest, continue searching in the parent test,
			// starting from the point of the call to Run which created this subtest.
			if c.level > 1 {
				frames = runtime.CallersFrames(c.creator)
				parent := c.parent
				// We're no longer looking at the current c after this point,
				// so we should unlock its mu, unless it's the original receiver,
				// in which case our caller doesn't expect us to do that.
				if shouldUnlock {
					c.mu.Unlock()
				}
				c = parent
				// Remember to unlock c.mu when we no longer need it, either
				// because we went up another nesting level, or because we
				// returned.
				shouldUnlock = true
				c.mu.Lock()
				continue
			}
			return prevFrame
		}
		if _, ok := c.helpers[frame.Function]; !ok {
			// Found a frame that wasn't inside a helper function.
			return frame
		}
	}
	return firstFrame
}

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation spaces for formatting.
// This function must be called with c.mu held.
func (c *common) decorate(s string, skip int) string {
	frame := c.frameSkip(skip)
	file := frame.File
	line := frame.Line
	if file != "" {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
	}
	if line == 0 {
		line = 1
	}
	buf := new(strings.Builder)
	// Every line is indented at least 4 spaces.
	buf.WriteString("    ")
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an additional 4 spaces.
			buf.WriteString("\n        ")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}

// TB is the interface common to T and B.
type TB interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate Go 1 compatibility.
	private()
}

var _ TB = (*B)(nil)

func (c *common) private() {}

// Name returns the name of the running test or benchmark.
func (c *common) Name() string {
	return c.name
}

// func (c *common) setRan() {
// 	if c.parent != nil {
// 		c.parent.setRan()
// 	}
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.ran = true
// }

// Fail marks the function as having failed but continues execution.
func (c *common) Fail() {
	if c.parent != nil {
		c.parent.Fail()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// c.done needs to be locked to synchronize checks to c.done in parent tests.
	if c.done {
		panic("Fail in goroutine after " + c.name + " has completed")
	}
	c.failed = true
}

// Failed reports whether the function has failed.
func (c *common) Failed() bool {
	c.mu.RLock()
	failed := c.failed
	c.mu.RUnlock()
	// return failed || c.raceErrors+race.Errors() > 0
	return failed
}

// FailNow marks the function as having failed and stops its execution
// by calling runtime.Goexit (which then runs all deferred calls in the
// current goroutine).
// Execution will continue at the next test or benchmark.
// FailNow must be called from the goroutine running the
// test or benchmark function, not from other goroutines
// created during the test. Calling FailNow does not stop
// those other goroutines.
func (c *common) FailNow() {
	c.Fail()
	runtime.Goexit()
}

// log generates the output. It's always at the same stack depth.
func (c *common) log(s string) {
	c.logDepth(s, 3) // logDepth + log + public function
}

// logDepth generates the output at an arbitrary stack depth.
func (c *common) logDepth(s string, depth int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.done {
		// This test has already finished. Try and log this message
		// with our parent. If we don't have a parent, panic.
		for parent := c.parent; parent != nil; parent = parent.parent {
			parent.mu.Lock()
			defer parent.mu.Unlock()
			if !parent.done {
				parent.output = append(parent.output, parent.decorate(s, depth+1)...)
				return
			}
		}
		panic("Log in goroutine after " + c.name + " has completed")
	} else {
		if c.chatty != nil {
			if c.bench {
				// Benchmarks don't print === CONT, so we should skip the test
				// printer and just print straight to stdout.
				fmt.Print(c.decorate(s, depth+1))
			} else {
				c.chatty.Printf(c.name, "%s", c.decorate(s, depth+1))
			}

			return
		}
		c.output = append(c.output, c.decorate(s, depth+1)...)
	}
}

// Log formats its arguments using default formatting, analogous to Println,
// and records the text in the error log. For tests, the text will be printed only if
// the test fails or the -test.v flag is set. For benchmarks, the text is always
// printed to avoid having performance depend on the value of the -test.v flag.
func (c *common) Log(args ...interface{}) { c.log(fmt.Sprintln(args...)) }

// Logf formats its arguments according to the format, analogous to Printf, and
// records the text in the error log. A final newline is added if not provided. For
// tests, the text will be printed only if the test fails or the -test.v flag is
// set. For benchmarks, the text is always printed to avoid having performance
// depend on the value of the -test.v flag.
func (c *common) Logf(format string, args ...interface{}) { c.log(fmt.Sprintf(format, args...)) }

// Error is equivalent to Log followed by Fail.
func (c *common) Error(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.Fail()
}

// Errorf is equivalent to Logf followed by Fail.
func (c *common) Errorf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.Fail()
}

// Fatal is equivalent to Log followed by FailNow.
func (c *common) Fatal(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.FailNow()
}

// Fatalf is equivalent to Logf followed by FailNow.
func (c *common) Fatalf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.FailNow()
}

// Skip is equivalent to Log followed by SkipNow.
func (c *common) Skip(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.SkipNow()
}

// Skipf is equivalent to Logf followed by SkipNow.
func (c *common) Skipf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.SkipNow()
}

// SkipNow marks the test as having been skipped and stops its execution
// by calling runtime.Goexit.
// If a test fails (see Error, Errorf, Fail) and is then skipped,
// it is still considered to have failed.
// Execution will continue at the next test or benchmark. See also FailNow.
// SkipNow must be called from the goroutine running the test, not from
// other goroutines created during the test. Calling SkipNow does not stop
// those other goroutines.
func (c *common) SkipNow() {
	c.skip()
	c.finished = true
	runtime.Goexit()
}

func (c *common) skip() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.skipped = true
}

// Skipped reports whether the test was skipped.
func (c *common) Skipped() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.skipped
}

// Helper marks the calling function as a test helper function.
// When printing file and line information, that function will be skipped.
// Helper may be called simultaneously from multiple goroutines.
func (c *common) Helper() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.helpers == nil {
		c.helpers = make(map[string]struct{})
	}
	c.helpers[callerName(1)] = struct{}{}
}

// Cleanup registers a function to be called when the test and all its
// subtests complete. Cleanup functions will be called in last added,
// first called order.
func (c *common) Cleanup(f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	oldCleanup := c.cleanup
	oldCleanupPc := c.cleanupPc
	c.cleanup = func() {
		if oldCleanup != nil {
			defer func() {
				c.mu.Lock()
				c.cleanupPc = oldCleanupPc
				c.mu.Unlock()
				oldCleanup()
			}()
		}
		c.mu.Lock()
		c.cleanupName = callerName(0)
		c.mu.Unlock()
		f()
	}
	var pc [maxStackLen]uintptr
	// Skip two extra frames to account for this function and runtime.Callers itself.
	n := runtime.Callers(2, pc[:])
	c.cleanupPc = pc[:n]
}

var tempDirReplacer struct {
	sync.Once
	r *strings.Replacer
}

// TempDir returns a temporary directory for the test to use.
// The directory is automatically removed by Cleanup when the test and
// all its subtests complete.
// Each subsequent call to t.TempDir returns a unique directory;
// if the directory creation fails, TempDir terminates the test by calling Fatal.
func (c *common) TempDir() string {
	// Use a single parent directory for all the temporary directories
	// created by a test, each numbered sequentially.
	c.tempDirOnce.Do(func() {
		c.Helper()

		// ioutil.TempDir doesn't like path separators in its pattern,
		// so mangle the name to accommodate subtests.
		tempDirReplacer.Do(func() {
			tempDirReplacer.r = strings.NewReplacer("/", "_", "\\", "_", ":", "_")
		})
		pattern := tempDirReplacer.r.Replace(c.Name())

		c.tempDir, c.tempDirErr = ioutil.TempDir("", pattern)
		if c.tempDirErr == nil {
			c.Cleanup(func() {
				if err := os.RemoveAll(c.tempDir); err != nil {
					c.Errorf("TempDir RemoveAll cleanup: %v", err)
				}
			})
		}
	})
	if c.tempDirErr != nil {
		c.Fatalf("TempDir: %v", c.tempDirErr)
	}
	seq := atomic.AddInt32(&c.tempDirSeq, 1)
	dir := fmt.Sprintf("%s%c%03d", c.tempDir, os.PathSeparator, seq)
	if err := os.Mkdir(dir, 0777); err != nil {
		c.Fatalf("TempDir: %v", err)
	}
	return dir
}

// panicHanding is an argument to runCleanup.
type panicHandling int

const (
	normalPanic panicHandling = iota
	recoverAndReturnPanic
)

// runCleanup is called at the end of the test.
// If catchPanic is true, this will catch panics, and return the recovered
// value if any.
func (c *common) runCleanup(ph panicHandling) (panicVal interface{}) {
	c.mu.Lock()
	cleanup := c.cleanup
	c.cleanup = nil
	c.mu.Unlock()
	if cleanup == nil {
		return nil
	}

	if ph == recoverAndReturnPanic {
		defer func() {
			panicVal = recover()
		}()
	}

	cleanup()
	return nil
}

// callerName gives the function name (qualified with a package path)
// for the caller after skip frames (where 0 means the current function).
func callerName(skip int) string {
	// Make room for the skip PC.
	var pc [1]uintptr
	n := runtime.Callers(skip+2, pc[:]) // skip + runtime.Callers + callerName
	if n == 0 {
		panic("testing: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

// No one should be using func Main anymore.
// See the doc comment on func Main and use MainStart instead.
var errMain = errors.New("testing: unexpected use of func Main")

type matchStringOnly func(pat, str string) (bool, error)

func (f matchStringOnly) MatchString(pat, str string) (bool, error)   { return f(pat, str) }
func (f matchStringOnly) StartCPUProfile(w io.Writer) error           { return errMain }
func (f matchStringOnly) StopCPUProfile()                             {}
func (f matchStringOnly) WriteProfileTo(string, io.Writer, int) error { return errMain }
func (f matchStringOnly) ImportPath() string                          { return "" }
func (f matchStringOnly) StartTestLog(io.Writer)                      {}
func (f matchStringOnly) StopTestLog() error                          { return errMain }

// Main is an internal function, part of the implementation of the "go test" command.
// It was exported because it is cross-package and predates "internal" packages.
// It is no longer used by "go test" but preserved, as much as possible, for other
// systems that simulate "go test" using Main, but Main sometimes cannot be updated as
// new functionality is added to the testing package.
// Systems simulating "go test" should be updated to use MainStart.
// func Main(matchString func(pat, str string) (bool, error), tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample) {
func Main(matchString func(pat, str string) (bool, error), benchmarks []InternalBenchmark) {
	os.Exit(MainStart(matchStringOnly(matchString), benchmarks).Run())
}

// M is a type passed to a TestMain function to run the actual tests.
type M struct {
	deps testDeps
	// tests      []InternalTest
	benchmarks []InternalBenchmark
	// examples   []InternalExample

	// timer     *time.Timer
	afterOnce sync.Once

	numRun int

	// value to pass to os.Exit, the outer test func main
	// harness calls os.Exit with this code. See #34129.
	exitCode int
}

// testDeps is an internal interface of functionality that is
// passed into this package by a test's generated main package.
// The canonical implementation of this interface is
// testing/internal/testdeps's TestDeps.
type testDeps interface {
	ImportPath() string
	MatchString(pat, str string) (bool, error)
	StartCPUProfile(io.Writer) error
	StopCPUProfile()
	StartTestLog(io.Writer)
	StopTestLog() error
	WriteProfileTo(string, io.Writer, int) error
}

// MainStart is meant for use by tests generated by 'go test'.
// It is not meant to be called directly and is not subject to the Go 1 compatibility document.
// It may change signature from release to release.
// func MainStart(deps testDeps, tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample) *M {

func MainStart(deps testDeps, benchmarks []InternalBenchmark) *M {

	Init()
	return &M{
		deps: deps,
		// tests:      tests,
		benchmarks: benchmarks,
		// examples:   examples,
	}
}

// Run runs the tests. It returns an exit code to pass to os.Exit.
func (m *M) Run() (code int) {
	defer func() {
		code = m.exitCode
	}()

	// Count the number of calls to m.Run.
	// We only ever expected 1, but we didn't enforce that,
	// and now there are tests in the wild that call m.Run multiple times.
	// Sigh. golang.org/issue/23129.
	m.numRun++

	// TestMain may have already called flag.Parse.
	if !flag.Parsed() {
		flag.Parse()
	}

	if *parallel < 1 {
		fmt.Fprintln(os.Stderr, "testing: -parallel can only be given a positive integer")
		flag.Usage()
		m.exitCode = 2
		return
	}

	if len(*matchList) != 0 {
		// listTests(m.deps.MatchString, m.tests, m.benchmarks, m.examples)
		listTests(m.deps.MatchString, m.benchmarks)
		m.exitCode = 0
		return
	}

	parseCpuList()

	m.before()
	defer m.after()
	// deadline := m.startAlarm()
	// haveExamples = len(m.examples) > 0
	// testRan, testOk := runTests(m.deps.MatchString, m.tests, deadline)
	// exampleRan, exampleOk := runExamples(m.deps.MatchString, m.examples)
	// m.stopAlarm()
	// if !testRan && !exampleRan && *matchBenchmarks == "" {
	// 	fmt.Fprintln(os.Stderr, "testing: warning: no tests to run")
	// }
	// if !testOk || !exampleOk || !runBenchmarks(m.deps.ImportPath(), m.deps.MatchString, m.benchmarks) || race.Errors() > 0 {
	if !runBenchmarks(m.deps.ImportPath(), m.deps.MatchString, m.benchmarks) {
		fmt.Println("FAIL")
		m.exitCode = 1
		return
	}

	fmt.Println("PASS")
	m.exitCode = 0
	return
}

// func (t *T) report() {
// 	if t.parent == nil {
// 		return
// 	}
// 	dstr := fmtDuration(t.duration)
// 	format := "--- %s: %s (%s)\n"
// 	if t.Failed() {
// 		t.flushToParent(t.name, format, "FAIL", t.name, dstr)
// 	} else if t.chatty != nil {
// 		if t.Skipped() {
// 			t.flushToParent(t.name, format, "SKIP", t.name, dstr)
// 		} else {
// 			t.flushToParent(t.name, format, "PASS", t.name, dstr)
// 		}
// 	}
// }

// func listTests(matchString func(pat, str string) (bool, error), tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample) {
func listTests(matchString func(pat, str string) (bool, error), benchmarks []InternalBenchmark) {
	if _, err := matchString(*matchList, "non-empty"); err != nil {
		fmt.Fprintf(os.Stderr, "testing: invalid regexp in -test.list (%q): %s\n", *matchList, err)
		os.Exit(1)
	}

	// for _, test := range tests {
	// 	if ok, _ := matchString(*matchList, test.Name); ok {
	// 		fmt.Println(test.Name)
	// 	}
	// }
	for _, bench := range benchmarks {
		if ok, _ := matchString(*matchList, bench.Name); ok {
			fmt.Println(bench.Name)
		}
	}
	// for _, example := range examples {
	// 	if ok, _ := matchString(*matchList, example.Name); ok {
	// 		fmt.Println(example.Name)
	// 	}
	// }
}

// // RunTests is an internal function but exported because it is cross-package;
// // it is part of the implementation of the "go test" command.
// func RunTests(matchString func(pat, str string) (bool, error), tests []InternalTest) (ok bool) {
// 	var deadline time.Time
// 	if *timeout > 0 {
// 		deadline = time.Now().Add(*timeout)
// 	}
// 	ran, ok := runTests(matchString, tests, deadline)
// 	if !ran && !haveExamples {
// 		fmt.Fprintln(os.Stderr, "testing: warning: no tests to run")
// 	}
// 	return ok
// }

// func runTests(matchString func(pat, str string) (bool, error), tests []InternalTest, deadline time.Time) (ran, ok bool) {
// 	ok = true
// 	for _, procs := range cpuList {
// 		runtime.GOMAXPROCS(procs)
// 		for i := uint(0); i < *count; i++ {
// 			if shouldFailFast() {
// 				break
// 			}
// 			ctx := newTestContext(*parallel, newMatcher(matchString, *match, "-test.run"))
// 			ctx.deadline = deadline
// 			t := &T{
// 				common: common{
// 					signal:  make(chan bool),
// 					barrier: make(chan bool),
// 					w:       os.Stdout,
// 				},
// 				context: ctx,
// 			}
// 			if Verbose() {
// 				t.chatty = newChattyPrinter(t.w)
// 			}
// 			tRunner(t, func(t *T) {
// 				for _, test := range tests {
// 					t.Run(test.Name, test.F)
// 				}
// 				// Run catching the signal rather than the tRunner as a separate
// 				// goroutine to avoid adding a goroutine during the sequential
// 				// phase as this pollutes the stacktrace output when aborting.
// 				go func() { <-t.signal }()
// 			})
// 			ok = ok && !t.Failed()
// 			ran = ran || t.ran
// 		}
// 	}
// 	return ran, ok
// }

// before runs before all testing.
func (m *M) before() {
	if *memProfileRate > 0 {
		runtime.MemProfileRate = *memProfileRate
	}
	if *cpuProfile != "" {
		f, err := os.Create(toOutputDir(*cpuProfile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			return
		}
		if err := m.deps.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't start cpu profile: %s\n", err)
			f.Close()
			return
		}
		// Could save f so after can call f.Close; not worth the effort.
	}
	if *traceFile != "" {
		f, err := os.Create(toOutputDir(*traceFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			return
		}
		if err := trace.Start(f); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't start tracing: %s\n", err)
			f.Close()
			return
		}
		// Could save f so after can call f.Close; not worth the effort.
	}
	if *blockProfile != "" && *blockProfileRate >= 0 {
		runtime.SetBlockProfileRate(*blockProfileRate)
	}
	if *mutexProfile != "" && *mutexProfileFraction >= 0 {
		runtime.SetMutexProfileFraction(*mutexProfileFraction)
	}
	if *coverProfile != "" && cover.Mode == "" {
		fmt.Fprintf(os.Stderr, "testing: cannot use -test.coverprofile because test binary was not built with coverage enabled\n")
		os.Exit(2)
	}
	if *testlog != "" {
		// Note: Not using toOutputDir.
		// This file is for use by cmd/go, not users.
		var f *os.File
		var err error
		if m.numRun == 1 {
			f, err = os.Create(*testlog)
		} else {
			f, err = os.OpenFile(*testlog, os.O_WRONLY, 0)
			if err == nil {
				f.Seek(0, io.SeekEnd)
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			os.Exit(2)
		}
		m.deps.StartTestLog(f)
		testlogFile = f
	}
}

// after runs after all testing.
func (m *M) after() {
	m.afterOnce.Do(func() {
		m.writeProfiles()
	})
}

func (m *M) writeProfiles() {
	if *testlog != "" {
		if err := m.deps.StopTestLog(); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't write %s: %s\n", *testlog, err)
			os.Exit(2)
		}
		if err := testlogFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't write %s: %s\n", *testlog, err)
			os.Exit(2)
		}
	}
	if *cpuProfile != "" {
		m.deps.StopCPUProfile() // flushes profile to disk
	}
	if *traceFile != "" {
		trace.Stop() // flushes trace to disk
	}
	if *memProfile != "" {
		f, err := os.Create(toOutputDir(*memProfile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			os.Exit(2)
		}
		runtime.GC() // materialize all statistics
		if err = m.deps.WriteProfileTo("allocs", f, 0); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't write %s: %s\n", *memProfile, err)
			os.Exit(2)
		}
		f.Close()
	}
	if *blockProfile != "" && *blockProfileRate >= 0 {
		f, err := os.Create(toOutputDir(*blockProfile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			os.Exit(2)
		}
		if err = m.deps.WriteProfileTo("block", f, 0); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't write %s: %s\n", *blockProfile, err)
			os.Exit(2)
		}
		f.Close()
	}
	if *mutexProfile != "" && *mutexProfileFraction >= 0 {
		f, err := os.Create(toOutputDir(*mutexProfile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			os.Exit(2)
		}
		if err = m.deps.WriteProfileTo("mutex", f, 0); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't write %s: %s\n", *mutexProfile, err)
			os.Exit(2)
		}
		f.Close()
	}
	if cover.Mode != "" {
		coverReport()
	}
}

// toOutputDir returns the file name relocated, if required, to outputDir.
// Simple implementation to avoid pulling in path/filepath.
func toOutputDir(path string) string {
	if *outputDir == "" || path == "" {
		return path
	}
	// On Windows, it's clumsy, but we can be almost always correct
	// by just looking for a drive letter and a colon.
	// Absolute paths always have a drive letter (ignoring UNC).
	// Problem: if path == "C:A" and outputdir == "C:\Go" it's unclear
	// what to do, but even then path/filepath doesn't help.
	// TODO: Worth doing better? Probably not, because we're here only
	// under the management of go test.
	if runtime.GOOS == "windows" && len(path) >= 2 {
		letter, colon := path[0], path[1]
		if ('a' <= letter && letter <= 'z' || 'A' <= letter && letter <= 'Z') && colon == ':' {
			// If path starts with a drive letter we're stuck with it regardless.
			return path
		}
	}
	if os.IsPathSeparator(path[0]) {
		return path
	}
	return fmt.Sprintf("%s%c%s", *outputDir, os.PathSeparator, path)
}

func parseCpuList() {
	for _, val := range strings.Split(*cpuListStr, ",") {
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		cpu, err := strconv.Atoi(val)
		if err != nil || cpu <= 0 {
			fmt.Fprintf(os.Stderr, "testing: invalid value %q for -test.cpu\n", val)
			os.Exit(1)
		}
		cpuList = append(cpuList, cpu)
	}
	if cpuList == nil {
		cpuList = append(cpuList, runtime.GOMAXPROCS(-1))
	}
}
