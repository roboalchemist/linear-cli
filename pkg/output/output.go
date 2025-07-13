package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// TableData represents data for table output
type TableData struct {
	Headers []string
	Rows    [][]string
}

// JSON outputs data as JSON
func JSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

// Error outputs an error message
func Error(message string, plaintext, jsonOut bool) {
	if jsonOut {
		JSON(map[string]interface{}{
			"error": message,
		})
	} else if plaintext {
		fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", color.New(color.FgRed).Sprint("❌"), message)
	}
}

// Success outputs a success message
func Success(message string, plaintext, jsonOut bool) {
	if jsonOut {
		JSON(map[string]interface{}{
			"status":  "success",
			"message": message,
		})
	} else if plaintext {
		fmt.Println(message)
	} else {
		fmt.Printf("%s %s\n", color.New(color.FgGreen).Sprint("✅"), message)
	}
}

// Table outputs data in table format
func Table(data TableData, plaintext, jsonOut bool) {
	if jsonOut {
		// Convert table data to JSON
		jsonData := make([]map[string]interface{}, len(data.Rows))
		for i, row := range data.Rows {
			item := make(map[string]interface{})
			for j, header := range data.Headers {
				if j < len(row) {
					item[strings.ToLower(header)] = row[j]
				}
			}
			jsonData[i] = item
		}
		JSON(jsonData)
		return
	}

	if plaintext {
		// Simple plaintext output
		if len(data.Headers) > 0 {
			fmt.Println(strings.Join(data.Headers, "\t"))
		}
		for _, row := range data.Rows {
			fmt.Println(strings.Join(row, "\t"))
		}
		return
	}

	// Rich table output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(data.Headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	// Add color to headers
	coloredHeaders := make([]string, len(data.Headers))
	for i, header := range data.Headers {
		coloredHeaders[i] = color.New(color.FgCyan, color.Bold).Sprint(header)
	}
	table.SetHeader(coloredHeaders)

	for _, row := range data.Rows {
		table.Append(row)
	}
	table.Render()
}

// Info outputs an informational message
func Info(message string, plaintext, jsonOut bool) {
	if jsonOut {
		JSON(map[string]interface{}{
			"info": message,
		})
	} else if plaintext {
		fmt.Println(message)
	} else {
		fmt.Printf("%s %s\n", color.New(color.FgBlue).Sprint("ℹ️"), message)
	}
}