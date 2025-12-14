package main

import "os"

func main() {
	os.Exit(1) // want "os.Exit call within main package main func"

	defer func() {
		os.Exit(2) // want "os.Exit call within main package main func"
	}()
	x()
}

func x() {
	os.Exit(3) // want
}
