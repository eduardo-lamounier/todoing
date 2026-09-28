package main

import "unicode"

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
