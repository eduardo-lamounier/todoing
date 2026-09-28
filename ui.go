package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

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
) error {
	const (
		minWidth = 0
		tabWidth = 4
		padding  = 4
		padChar  = ' '
		flags    = 0
	)

	n := len(colsName)

	if n == 0 {
		return fmt.Errorf("no column names passed")
	}

	for _, row := range rows {
		if len(row) != n {
			return fmt.Errorf("all rows should have the same amount of columns")
		}
	}

	totalWidth := totalWidth(rows, colsName, padding)

	colored := color.New(color.FgHiWhite)

	var buff bytes.Buffer
	tabWriter := tabwriter.NewWriter(
		&buff, minWidth, tabWidth, padding, padChar, flags,
	)

	colored.Printf("%*s\n\n", totalWidth/2+len(title)/2, title)

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

	if len(rows) == 0 {
		fmt.Printf("%*s\n", totalWidth/2+len(noRowMessage)/2, noRowMessage)
	}

	tabWriter.Flush()

	// Printing into stdout:
	lines := strings.Split(buff.String(), "\n")

	colored.Println(lines[0])
	for _, line := range lines[1:] {
		fmt.Println(line)
	}

	return nil
}
