package analyzer_test

import (
	"testing"

	"github.com/jonesrussell/go-fxlint/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	// Run all test cases
	analysistest.Run(t, testdata, analyzer.Analyzer, "a/...")
}
