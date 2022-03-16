package main

import (
	"github.com/ari1021/gencon"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(gencon.Analyzer) }
