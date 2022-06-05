package main

import (
	"github.com/arkuchy/gencon"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(gencon.Analyzer) }
