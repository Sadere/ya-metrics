package main

import "os"

func main() {
	os.Exit(1) // want "usage of os.Exit inside main function of main package is prohibited"
}