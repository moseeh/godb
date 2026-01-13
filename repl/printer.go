package repl

import (
	"fmt"
	"godb/engine"
	"strings"
)

// PrintRows formats and prints rows in a table format
func PrintRows(rows []engine.Row) {
	if len(rows) == 0 {
		fmt.Println("No rows returned.")
		return
	}

	// Collect all column names
	columnSet := make(map[string]bool)
	for _, row := range rows {
		for col := range row {
			columnSet[col] = true
		}
	}

	// Convert to sorted list for consistent ordering
	columns := make([]string, 0, len(columnSet))
	for col := range columnSet {
		columns = append(columns, col)
	}

	// Calculate column widths
	widths := make(map[string]int)
	for _, col := range columns {
		widths[col] = len(col)
	}

	for _, row := range rows {
		for _, col := range columns {
			if val, ok := row[col]; ok {
				valStr := fmt.Sprintf("%v", val)
				if len(valStr) > widths[col] {
					widths[col] = len(valStr)
				}
			}
		}
	}

	// Print header
	var headerParts []string
	for _, col := range columns {
		headerParts = append(headerParts, padRight(col, widths[col]))
	}
	fmt.Println(strings.Join(headerParts, " | "))

	// Print separator
	var separatorParts []string
	for _, col := range columns {
		separatorParts = append(separatorParts, strings.Repeat("-", widths[col]))
	}
	fmt.Println(strings.Join(separatorParts, "-+-"))

	// Print rows
	for _, row := range rows {
		var rowParts []string
		for _, col := range columns {
			val := ""
			if v, ok := row[col]; ok && v != nil {
				val = fmt.Sprintf("%v", v)
			}
			rowParts = append(rowParts, padRight(val, widths[col]))
		}
		fmt.Println(strings.Join(rowParts, " | "))
	}

	fmt.Printf("\n%d row(s) returned.\n", len(rows))
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}

// PrintError prints an error message
func PrintError(err error) {
	fmt.Printf("✗ Error: %v\n", err)
}

// padRight pads a string to a given width with spaces on the right
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
