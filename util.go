package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}

	chars := []rune(s)
	chars[0] = unicode.ToUpper(chars[0])
	return string(chars)
}

func Assert(cond bool, v ...any) {
	if !cond {
		panic(v)
	}
}

func Mod(n int, d int) int {
	return (n%d + d) % d
}

func ConsoleClearLine() {
	fmt.Print("\r\033[K")
}

func logFatalError(err error) {
	for errMsg := range strings.SplitSeq(fmt.Sprintf("%s", err), "\n") {
		fmt.Fprintf(os.Stderr, "ERROR: %s.\n", errMsg)
	}
	os.Exit(1)
}

func DaysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
