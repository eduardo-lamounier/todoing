package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	msgProgramInterrupted = "the program got interrupted"
)

func PrintFatalError(err error) {
	for errMsg := range strings.SplitSeq(fmt.Sprintf("%s", err), "\n") {
		fmt.Fprintf(os.Stderr, "ERROR: %s.\n", errMsg)
	}
}
