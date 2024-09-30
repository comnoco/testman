package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var opts Opts

// Define a custom flag type for string slices
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		if err != flag.ErrHelp {
			log.Fatalf("error: %v", err)
		}
		// os.Exit(1)
	}
}

func run(args []string) error {
	// flags
	testFlags := flag.NewFlagSet("testman test", flag.ExitOnError)
	testFlags.BoolVar(&opts.Verbose, "v", false, "verbose")
	testFlags.Var(&opts.Run, "run", "regex to run tests and examples (can be specified multiple times)")
	testFlags.Var(&opts.Skip, "skip", "regex to skip tests and examples (can be specified multiple times)")
	testFlags.IntVar(&opts.Retry, "retry", 0, "fail after N retries")
	testFlags.DurationVar(&opts.Timeout, "timeout", 0, "program max duration")
	testFlags.BoolVar(&opts.ContinueOnError, "continue-on-error", false, "continue on error (but still fails at the end)")
	testFlags.DurationVar(&opts.TestTimeout, "test.timeout", 0, "`go test -timeout=VAL`")
	testFlags.IntVar(&opts.TestCount, "test.count", 1, "`go test -count=VAL`")
	testFlags.BoolVar(&opts.TestV, "test.v", false, "`go test -v`")
	testFlags.BoolVar(&opts.TestRace, "test.race", false, "`go test -race`")
	testFlags.BoolVar(&opts.RegexCaseInsensitive, "i", false, "case insensitive regex matching")

	listFlags := flag.NewFlagSet("testman list", flag.ExitOnError)
	listFlags.BoolVar(&opts.Verbose, "v", false, "verbose")
	listFlags.Var(&opts.Skip, "skip", "regex to skip tests and examples (can be specified multiple times)")
	listFlags.Var(&opts.Run, "run", "regex to run tests and examples (can be specified multiple times)")
	listFlags.BoolVar(&opts.RegexCaseInsensitive, "i", false, "case insensitive regex matching")

	root := &ffcli.Command{
		ShortUsage: "testman <subcommand> [flags]",
		ShortHelp:  "Advanced testing workflows for Go projects.",
		Exec: func(ctx context.Context, args []string) error {
			_ = ctx
			_ = args
			return flag.ErrHelp
		},
		Subcommands: []*ffcli.Command{
			{
				Name:       "test",
				FlagSet:    testFlags,
				ShortHelp:  "advanced go test workflows",
				ShortUsage: "testman test [flags] [packages]",
				LongHelp:   testLongHelp,
				Exec:       runTest,
			}, {
				Name:       "list",
				FlagSet:    listFlags,
				ShortHelp:  "list available tests",
				ShortUsage: "testman list [packages]",
				LongHelp:   listLongHelp,
				Exec:       runList,
			},
		},
	}

	return root.ParseAndRun(context.Background(), args[1:])
}

const (
	testLongHelp = `EXAMPLES
   testman test ./...
   testman test -v ./...
   testman test -skip ^TestUnstable -timeout=300s -retry=50 ./...
   testman test -skip ^TestBroken -test.timeout=30s -retry=10 --continue-on-error ./...
	 testman test -skip slow -run stable -i ./...
	 testman test -run ^TestUnstable -timeout=300s -retry=50 ./...
   testman test -run ^TestBroken -test.timeout=30s -retry=10 --continue-on-error ./...
   testman test -test.timeout=10s -test.v -test.count=2 -test.race`
	listLongHelp = `EXAMPLES
   testman list ./...
   testman list -v ./...
   testman list -skip ^TestStable ./...
	 testman list -run stable -i ./...
	 testman list -run ^TestStable ./...`
)

func runList(ctx context.Context, args []string) error {
	_ = ctx
	if len(args) == 0 {
		return flag.ErrHelp
	}
	cleanup, err := preRun()
	if err != nil {
		return err
	}
	defer cleanup()

	// list packages
	pkgs, err := listPackagesWithTests(args)
	if err != nil {
		return err
	}

	// list tests
	for _, pkg := range pkgs {
		tests, err := listDirTests(pkg.Dir)
		if err != nil {
			return err
		}
		if len(tests) == 0 {
			continue
		}

		fmt.Println(pkg.ImportPath)
		for _, test := range tests {
			fmt.Printf("  %s\n", test)
		}
	}
	return nil
}

func runTest(ctx context.Context, args []string) error {
	_ = ctx
	if len(args) == 0 {
		return flag.ErrHelp
	}
	cleanup, err := preRun()
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Printf("runTest opts=%s args=%s", JSON(opts), JSON(args))
	start := time.Now()

	if opts.Timeout > 0 {
		go func() {
			<-time.After(opts.Timeout)
			fmt.Printf("FAIL: timed out after %s\n", time.Since(start))
			panic(fmt.Sprintf("timed out after %s", time.Since(start)))
		}()
	}

	// list packages
	pkgs, err := listPackagesWithTests(args)
	if err != nil {
		return err
	}

	atLeastOneFailure := false
	// list tests
	for _, pkg := range pkgs {
		tests, err := listDirTests(pkg.Dir)
		if err != nil {
			return err
		}
		if len(tests) == 0 {
			continue
		}

		pkgStart := time.Now()
		// compile test binary
		bin, err := compileTestBin(pkg, opts.TmpDir)
		if err != nil {
			fmt.Printf("FAIL\t%s\t[compile error: %v]\n", pkg.ImportPath, err)
			return err
		}

		isPackageOK := true
		for _, test := range tests {
			args := []string{
				fmt.Sprintf("-test.count=%d", opts.TestCount),
			}
			timeout := opts.TestTimeout
			if opts.Timeout > 0 && opts.Timeout < opts.TestTimeout {
				timeout = opts.Timeout
			}
			if timeout > 0 {
				args = append(args, fmt.Sprintf("-test.timeout=%s", timeout))
			}
			if opts.TestV {
				args = append(args, "-test.v")
			}
			args = append(args, "-test.run", fmt.Sprintf("^%s$", test))
			for i := opts.Retry; i >= 0; i-- {
				cmd := exec.Command(bin, args...)
				fmt.Println(cmd.String())
				out, err := cmd.CombinedOutput()
				if err != nil {
					if i == 0 {
						fmt.Printf("FAIL\t%s.%s\t[test error: %v]\n", pkg.ImportPath, test, err)
						isPackageOK = false
						atLeastOneFailure = true
					} else if opts.Verbose {
						fmt.Printf("RETRY\t%s.%s\t[test error: %v]\n", pkg.ImportPath, test, err)
					}
					if opts.Verbose {
						fmt.Println(string(out))
					}
				} else {
					fmt.Printf("ok\t%s.%s\n", pkg.ImportPath, test)
					break
				}
			}
		}
		if isPackageOK {
			fmt.Printf("ok\t%s\t%s\n", pkg.ImportPath, time.Since(pkgStart))
		}
	}

	fmt.Printf("total: %s\n", time.Since(start))
	if atLeastOneFailure {
		return errors.New("at least one failure occurred")
	}
	return nil
}

func preRun() (func(), error) {
	if !opts.Verbose {
		log.SetOutput(io.Discard)
	}

	// create temp dir
	var err error
	opts.TmpDir, err = os.MkdirTemp("", "testman")
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		os.RemoveAll(opts.TmpDir)
	}
	return cleanup, nil
}

func compileTestBin(pkg Package, tempdir string) (string, error) {
	name := strings.ReplaceAll(pkg.ImportPath, "/", "~")
	bin := filepath.Join(tempdir, name)
	args := []string{"test", "-c"}
	if opts.TestV {
		args = append(args, "-v")
	}
	if opts.TestRace {
		args = append(args, "-race")
	}
	args = append(args, "-o", bin)
	cmd := exec.Command("go", args...)
	cmd.Dir = pkg.Dir
	fmt.Println(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return "", err
	}

	return bin, nil
}

func listDirTests(dir string) ([]string, error) {
	cmd := exec.Command("go", "test", "-list", ".")
	cmd.Dir = dir
	fmt.Println(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(out)) == "" {
		return nil, nil
	}
	tests := []string{}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ok ") {
			continue
		}

		tests = append(tests, line)
	}

	if len(opts.Skip) > 0 {
		fmt.Println("skip", opts.Skip)
		var filteredTests []string
		for _, test := range tests {
			shouldKeep := true
			for _, skip := range opts.Skip {
				if opts.RegexCaseInsensitive {
					skip = "(?i)" + skip
				}

				matched, err := regexp.MatchString(skip, test)
				if err != nil {
					return nil, err
				}
				if matched {
					shouldKeep = false
					break
				}
			}
			if shouldKeep {
				filteredTests = append(filteredTests, test)
			}
		}
		tests = filteredTests
	}

	if len(opts.Run) > 0 {
		fmt.Println("run", opts.Run)
		var filteredTests []string
		for _, test := range tests {
			shouldKeep := false
			for _, run := range opts.Run {
				if opts.RegexCaseInsensitive {
					run = "(?i)" + run
				}

				matched, err := regexp.MatchString(run, test)
				if err != nil {
					return nil, err
				}
				if matched {
					shouldKeep = true
					break
				}
			}
			if shouldKeep {
				filteredTests = append(filteredTests, test)
			}
		}
		tests = filteredTests
	}

	return tests, nil
}

func listPackagesWithTests(patterns []string) ([]Package, error) {
	cmdArgs := append([]string{"list", "-test", "-f", "{{.ImportPath}} {{.Dir}}"}, patterns...)
	cmd := exec.Command("go", cmdArgs...)
	fmt.Println(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	pkgs := []Package{}
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if !strings.HasSuffix(parts[0], ".test") {
			continue
		}
		pkgs = append(pkgs, Package{
			ImportPath: strings.TrimSuffix(parts[0], ".test"),
			Dir:        parts[1],
		})
	}
	return pkgs, nil
}

type Package struct {
	Dir        string
	ImportPath string
}

type Opts struct {
	Verbose              bool
	Skip                 stringSliceFlag
	Run                  stringSliceFlag
	Timeout              time.Duration
	Retry                int
	TmpDir               string
	ContinueOnError      bool
	RegexCaseInsensitive bool

	// test.*
	TestTimeout time.Duration
	TestV       bool
	TestCount   int
	TestRace    bool
}

// Insignificant change to trigger a build
// Insignificant change to trigger a build
