package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/k9withabone/fluttui/tui"
)

func main() {
	_, err := exec.LookPath("flutter")
	if err != nil {
		fmt.Println(
			"Error: could not locate `flutter`.",
			"\nPlease make sure `flutter` is installed and",
			"available in your PATH before using this tool.",
		)
		os.Exit(1)
	}

	tui.StartTea()
}
