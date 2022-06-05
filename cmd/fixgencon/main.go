package main

import (
	"os"

	"github.com/arkuchy/gencon"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	// on fix flag forcibly
	os.Args = append([]string{os.Args[0], "-fix"}, os.Args[1:]...)
	singlechecker.Main(gencon.Analyzer)
}
