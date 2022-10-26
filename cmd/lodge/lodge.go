package main

import (
	"os"
	"path/filepath"
	"fmt"

//	"github.com/charlesap/sybil/pkg/cmd"
//	"github.com/charlesap/sybil/pkg/cmd/lodge"
)

func main() {

	baseName := filepath.Base(os.Args[0])

	fmt.Println(baseName)
}
