package linters

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzer(t *testing.T) {

	analysistest.Run(t, analysistest.TestData(), NewExitAnalyzer(), "./osexit")
}
