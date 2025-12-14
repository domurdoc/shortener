package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"

	"github.com/kisielk/errcheck/errcheck"

	"github.com/domurdoc/shortener/cmd/staticlint/exitcheckanalyzer"
)

func main() {

	checks := []*analysis.Analyzer{
		exitcheckanalyzer.ExitCheckAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		errcheck.Analyzer,
	}
	for _, a := range staticcheck.Analyzers {
		checks = append(checks, a.Analyzer)
	}
	for _, a := range quickfix.Analyzers {
		checks = append(checks, a.Analyzer)
	}
	multichecker.Main(
		checks...,
	)
}
