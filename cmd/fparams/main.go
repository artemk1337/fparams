package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/artemk1337/fparams/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}
