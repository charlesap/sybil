package main

import (
	"os"
	"path/filepath"

	"github.com/charlesap/sybil/pkg/lodge"
)

func main() {

	baseName := filepath.Base(os.Args[0])

	lodge.Emit(baseName)
}
