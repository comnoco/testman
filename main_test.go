package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	closer := MustCaptureStdoutAndStderr()
	err := run([]string{"testman"})
	require.Equal(t, err, flag.ErrHelp)
	_ = closer()
}

func ExampleList_examples() {
	os.Stderr = os.Stdout // required hack to run the test as an example
	os.Chdir("examples")
	defer os.Chdir("..")
	err := run([]string{"testman", "list", "./..."})
	checkErr(err)
	// Output:
	// github.com/comnoco/testman/v2/examples/testpkg
	//   TestStableAlwaysSucceed
	//   TestUnstableMaySucceed
	//   TestBrokenAlwaysFailing
	//   ExampleAlwaysSucceed
}

func ExampleList() {
	err := run([]string{"testman", "list", "./..."})
	if err != nil {
		panic(err)
	}
	// Output:
	// github.com/comnoco/testman/v2
	//   TestRun
	//   ExampleList_examples
	//   ExampleList
}

func TestMain(m *testing.M) {
	// goleak.VerifyTestMain(m)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
