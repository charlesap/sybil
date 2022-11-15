package main

import (
        "os"
        "path/filepath"

	"github.com/ncw/directio"
        "github.com/charlesap/sybil/pkg/lodge"
)

func Attach() {

	storeA, err := directio.OpenFile("blockstore1", os.O_RDONLY, 0666)
	storeB, err := directio.OpenFile("blockstore2", os.O_RDONLY, 0666)

	if err!= nil {

		baseName := filepath.Base(os.Args[0])

		lodge.Emit(baseName)

		storeA.Close()
		storeB.Close()

	}
	
	
}
