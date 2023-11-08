package repos

import (
	"fmt"
	"testing"
	"time"

	_ "gitmic/test_init"

	"github.com/fatih/color"
)

func BenchmarkRunSimple(b *testing.B) {
	t := time.Now()

	_ = RunSimple(false)

	elapsed := time.Since(t)

	testName := color.YellowString("%s", "RunSimple")
	sinceTime := color.New(color.FgHiRed, color.Bold).Sprintf("%s", elapsed)
	fmt.Printf("%s elapsed time: %s\n", testName, sinceTime)
}

func BenchmarkRunConcurrency(b *testing.B) {
	t := time.Now()

	_ = RunConcurrency(false)

	elapsed := time.Since(t)

	testName := color.YellowString("%s", "RunConcurrency")
	sinceTime := color.New(color.FgHiRed, color.Bold).Sprintf("%s", elapsed)
	fmt.Printf("%s elapsed time: %s\n", testName, sinceTime)
}
