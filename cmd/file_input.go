package cmd

import (
	"fmt"
	"io"
	"os"
)

// readContentFromFile reads the entire content of a file and returns it as a string.
// If path is "-", it reads from stdin.
func readContentFromFile(path string) (string, error) {
	var data []byte
	var err error

	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}

	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", path, err)
	}

	return string(data), nil
}

// resolveBodyFromFlags resolves the text content from either a direct string flag or a file flag.
// flagName is the name of the direct text flag (e.g. "body", "description", "content").
// fileFlagName is the name of the file flag (e.g. "body-file", "description-file", "content-file").
// Returns the resolved text and any error.
func resolveBodyFromFlags(flagValue string, flagChanged bool, filePath string, flagName string, fileFlagName string) (string, error) {
	if flagChanged && filePath != "" {
		return "", fmt.Errorf("cannot use both --%s and --%s", flagName, fileFlagName)
	}

	if filePath != "" {
		content, err := readContentFromFile(filePath)
		if err != nil {
			return "", err
		}
		return content, nil
	}

	return flagValue, nil
}
