package main

import (
        "os"
        "path/filepath"
        "github.com/charlesap/sybil/pkg/lodge"
)

func Attach() {

	baseName := filepath.Base(os.Args[0])

	lodge.Emit(baseName)
	
	
}
