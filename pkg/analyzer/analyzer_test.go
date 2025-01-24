package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/go-fxlint/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	// Get the absolute path to the testdata directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")

	// Run all test cases
	analysistest.Run(t, testdata, analyzer.Analyzer, "a/...")
}
