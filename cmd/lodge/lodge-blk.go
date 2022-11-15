package main

import (
        "os"
        "path/filepath"

	"github.com/ncw/directio"
        "github.com/charlesap/sybil/pkg/lodge"
)

func Attach() {

	storage, err := directio.OpenFile("blockstore", os.O_RDONLY, 0666)

	if err!= nil {

		baseName := filepath.Base(os.Args[0])

		lodge.Emit(baseName)

		storage.Close()

	}
	
	
}
