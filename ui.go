package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

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

// ScrollInput It's guaranteed that even in case of an input error, the returned
// selected index will be in [0, maximum[.
func ScrollInput(maximum int, getInputText func(int) string, newline bool,
) (int, error) {
	current := 0

	showInputText := func() {
		ConsoleClearLine()
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
			current = Mod(current-1, maximum)
		case keys.Up:
			current = Mod(current+1, maximum)
		case keys.Enter:
			return true, nil
		}

		showInputText()
		return false, nil
	})

	if newline {
		fmt.Println()
	}

	return current, err
}
