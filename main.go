package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	Command := os.Args[1]
	switch Command {
	case "copy":
		return
	case "paste":
		fmt.Printf("cmm")
	case "serve":
	default:
	}
}
