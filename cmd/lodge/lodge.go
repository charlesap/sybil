package main

import (
	"os"
	"fmt"
        "github.com/charlesap/sybil/pkg/lodge"
)

func main() {

	base := lodge.Base {0,nil,"","",""}

	e := base.Init("blockstore1")
	if e != nil {
		fmt.Printf("lodge error: %s\n", e)
		os.Exit(1)
	}

	go HandleUDP()

	Webmain()

}
