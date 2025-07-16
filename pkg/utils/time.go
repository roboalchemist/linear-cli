package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTimeExpression converts time expressions like "3_weeks_ago" into ISO8601 datetime strings
// Returns empty string for "all_time"
// Default is "6_months_ago" if empty string is provided
func ParseTimeExpression(expr string) (string, error) {
	// Handle empty input - use default
	if expr == "" {
		expr = "6_months_ago"
	}

	// Handle special case
	if expr == "all_time" {
		return "", nil
	}

	// Try to parse as a date first (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", expr); err == nil {
		return expr + "T00:00:00Z", nil
	}

	// Try to parse as ISO8601
	if _, err := time.Parse(time.RFC3339, expr); err == nil {
		return expr, nil
	}

	// Parse relative time expressions
	parts := strings.Split(expr, "_")
	if len(parts) < 3 || parts[len(parts)-1] != "ago" {
		return "", fmt.Errorf("invalid time expression: %s (expected format like '3_weeks_ago' or 'all_time')", expr)
	}

	// Get the number
	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid number in time expression: %s", parts[0])
	}

	// Get the unit (handle both singular and plural)
	unit := strings.Join(parts[1:len(parts)-1], "_")

	// Calculate the time
	now := time.Now()
	var targetTime time.Time

	switch strings.TrimSuffix(unit, "s") {
	case "minute":
		targetTime = now.Add(-time.Duration(num) * time.Minute)
	case "hour":
		targetTime = now.Add(-time.Duration(num) * time.Hour)
	case "day":
		targetTime = now.AddDate(0, 0, -num)
	case "week":
		targetTime = now.AddDate(0, 0, -num*7)
	case "month":
		targetTime = now.AddDate(0, -num, 0)
	case "year":
		targetTime = now.AddDate(-num, 0, 0)
	default:
		return "", fmt.Errorf("invalid time unit: %s (valid units: minute, hour, day, week, month, year)", unit)
	}

	// Return as ISO8601 string
	return targetTime.Format(time.RFC3339), nil
}
