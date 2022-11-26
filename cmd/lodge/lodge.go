package main

import (
	"os"
	"flag"
	"fmt"
        "github.com/charlesap/sybil/pkg/lodge"
)

func parseFlags() {
}

func main() {

	var verb, prod, red, reinit bool
	var store string

	flag.BoolVar(&verb, "v", false, "if true, more logging")
	flag.BoolVar(&prod, "p", false, "if true, we start HTTPS server")
	flag.BoolVar(&red, "r", false, "if true, we redirect HTTP to HTTPS")
	flag.BoolVar(&reinit, "reinit", false, "if true, reinitialize the blockstore")
	flag.StringVar(&store, "b", "blockstore1", "blockstore to use")

	flag.Parse()


	base := lodge.Base {0,nil,0,"","",""}

	e := base.Init(store,reinit)

	if e != nil {
		fmt.Printf("lodge error: %s\n", e)
		os.Exit(1)
	}

	go HandleUDP()

	Webmain(verb,prod,red)

}
