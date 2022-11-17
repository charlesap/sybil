package main

import (
        "os"
	"fmt"
        "path/filepath"

	"github.com/ncw/directio"
        "github.com/charlesap/sybil/pkg/lodge"
)

func Attach() {

	fmt.Println("\nStarting up Lodge\n")

	baseName := filepath.Base(os.Args[0])

	lodge.Emit(baseName)

	storeA, err := directio.OpenFile("blockstore1", os.O_RDONLY, 0666)
	storeB, err := directio.OpenFile("blockstore2", os.O_RDONLY, 0666)

	fmt.Println("\nLodge Initialized\n")

	if err!= nil {

		storeA.Close()
		storeB.Close()

	}


}
