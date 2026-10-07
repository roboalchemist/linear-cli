package cmd

import (
	"encoding/json"
	"fmt"
)

// mergeJSONInput parses a JSON object supplied via a flag (e.g. --input-json)
// and merges its keys into input. Keys from the JSON object override values
// already present in input. An empty string is a no-op.
func mergeJSONInput(input map[string]interface{}, raw string) error {
	if raw == "" {
		return nil
	}
	var extra map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &extra); err != nil {
		return fmt.Errorf("invalid JSON object: %w", err)
	}
	for k, v := range extra {
		input[k] = v
	}
	return nil
}

// parseJSONValue parses an arbitrary JSON value (object, array, string, etc.)
// from a flag argument.
func parseJSONValue(raw string) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return v, nil
}
