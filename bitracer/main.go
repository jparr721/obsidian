package main

import (
	"fmt"
	"os"
)

func main() {
	arglen := len(os.Args)
	if arglen > 2 {
		fmt.Println("usage: bitracer [script]")
		os.Exit(64)
	} else if arglen == 2 {
		staticFileRunner(os.Args[1])
	} else {
		new(repl).start()
	}
}
