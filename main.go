package main

import (
	"fmt"
	"os"

	"httpreplay/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'capture', 'replay', or 'inspect' subcommands")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "capture":
		cmd.RunCapture()
	case "replay":
		replay()
	case "inspect":
		inspect()
	default:
		fmt.Println("Expected 'capture', 'replay', or 'inspect' subcommands")
		os.Exit(1)
	}
}

func replay() {
	fmt.Println("Replaying...")
}

func inspect() {
	fmt.Println("Inspecting...")
}
