package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/fatih/color"
)

// Returns the total width of the resulting table
func totalWidth(rows [][]string, colsName []string, padding int) int {
	n := len(colsName)
	widthCols := make([]int, n)

	for i, colName := range colsName {
		widthCols[i] = len(colName)
	}

	for _, row := range rows {
		for i, s := range row {
			widthCols[i] = max(widthCols[i], len(s))
		}
	}

	totalWidth := (n - 1) * padding
	for _, widthCol := range widthCols {
		totalWidth += widthCol
	}

	return totalWidth
}

func PrintTable(title string,
	noRowMessage string,
	rows [][]string,
	colsName []string,
) {
	const (
		minWidth = 0
		tabWidth = 4
		padding  = 4
		padChar  = ' '
		flags    = 0
	)

	n := len(colsName)

	Assert(n > 0, "The column names should be specified.")

	for _, row := range rows {
		Assert(len(row) == n, "All rows must have the same amount of columns")
	}

	totalWidth := totalWidth(rows, colsName, padding)

	colored := color.New(color.FgHiWhite)

	var buff bytes.Buffer
	tabWriter := tabwriter.NewWriter(
		&buff, minWidth, tabWidth, padding, padChar, flags,
	)

	colored.Printf("%*s\n\n", totalWidth/2+len(title)/2, title)

	if len(rows) == 0 {
		fmt.Printf("%*s\n", totalWidth/2+len(noRowMessage)/2, noRowMessage)
		return
	}

	// Printing the col names:
	for _, colName := range colsName[:n-1] {
		fmt.Fprintf(tabWriter, "%s\t", colName)
	}
	fmt.Fprintf(tabWriter, "%s\n", colsName[n-1])

	// Printing the rows:
	for _, row := range rows {
		for _, s := range row[:len(row)-1] {
			fmt.Fprintf(tabWriter, "%v\t", s)
		}
		fmt.Fprintf(tabWriter, "%v\n", row[len(row)-1])
	}

	tabWriter.Flush()

	// Printing into stdout:
	lines := strings.Split(buff.String(), "\n")

	colored.Println(lines[0])
	for _, line := range lines[1:] {
		fmt.Println(line)
	}
}

// WrapPrompt It's guaranteed that even in case of an input error, the returned
// selected index will be in [low, high[.
//
// If `begin` < `low`, then the beggining is `low`.
func WrapPrompt(begin int, low int, high int,
	getInputText func(int) string, newline bool,
) (int, error) {
	cursor.Hide()
	current := max(begin, low)

	showInputText := func() {
		cursor.StartOfLine()
		cursor.ClearLine()
		inputText := getInputText(current)
		inputText = strings.Trim(inputText, "\n") // Makes sure there isn't any
		// newline
		fmt.Print(inputText)
	}

	showInputText()

	err := keyboard.Listen(func(key keys.Key) (bool, error) {
		switch key.Code {
		case keys.CtrlC:
			return true, fmt.Errorf("the program got interrupted")
		case keys.Down:
			current = low + Mod(current-1-low, high-low)
		case keys.Up:
			current = low + Mod(current+1-low, high-low)
		case keys.Enter:
			return true, nil
		}

		showInputText()
		return false, nil
	})

	if newline {
		fmt.Println()
	}

	cursor.Show()
	return current, err
}

func YesNoPrompt(message string, yesByDefault bool, newline bool) (bool, error) {
	cursor.Hide()
	message = strings.Trim(message, "\n") // Makes sure there isn't any newline
	fmt.Println(message)

	confirmed := yesByDefault
	inputTextLen := len("[ N ] | [ Y ]")

	showInputText := func() {
		cursor.StartOfLine()
		cursor.ClearLine()
		if confirmed {
			fmt.Printf("%*s", len(message)/2+inputTextLen/2, "   N   | [ Y ] ")
		} else {
			fmt.Printf("%*s", len(message)/2+inputTextLen/2, " [ N ] |   Y   ")
		}
	}

	showInputText()

	err := keyboard.Listen(func(key keys.Key) (bool, error) {
		switch key.Code {
		case keys.CtrlC:
			return true, fmt.Errorf("the program got interrupted")
		case keys.Left:
			confirmed = false
		case keys.Right:
			confirmed = true
		case keys.Enter:
			return true, nil
		}

		showInputText()
		return false, nil
	})

	if newline {
		fmt.Println()
	}

	cursor.Show()
	return confirmed, err
}
