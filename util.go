package main

import (
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

func DaysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
