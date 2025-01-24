package plugin

import (
	"github.com/jonesrussell/go-fxlint/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin provides analyzers as a plugin.
// It is loaded by golangci-lint when built with the plugin build tag.
type AnalyzerPlugin struct{}

// GetAnalyzers returns analyzers implemented in this plugin.
func (*AnalyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}

// New creates a new AnalyzerPlugin for golangci-lint.
func New() *AnalyzerPlugin {
	return &AnalyzerPlugin{}
}
