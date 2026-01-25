package io

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"unicode"
)

// InputMode represents the detected input format.
type InputMode string

const (
	// ModeUnknown indicates format could not be detected.
	ModeUnknown InputMode = "unknown"
	// ModeBatchJSON indicates a JSON array or single object.
	ModeBatchJSON InputMode = "batch_json"
	// ModeNDJSON indicates newline-delimited JSON.
	ModeNDJSON InputMode = "ndjson"
	// ModeBatchYAML indicates a single YAML document.
	ModeBatchYAML InputMode = "batch_yaml"
	// ModeMultiYAML indicates multiple YAML documents (--- separated).
	ModeMultiYAML InputMode = "multi_yaml"
)

// DetectInputMode peeks at input to determine the format.
// Returns the detected mode and a new reader that includes the peeked bytes.
func DetectInputMode(r io.Reader) (InputMode, io.Reader, error) {
	bufReader := bufio.NewReaderSize(r, 4096)

	// Peek first bytes to determine format
	peek, err := bufReader.Peek(512)
	if err != nil && err != io.EOF {
		return ModeUnknown, bufReader, err
	}

	if len(peek) == 0 {
		return ModeUnknown, bufReader, nil
	}

	// Skip leading whitespace
	content := skipWhitespace(peek)
	if len(content) == 0 {
		return ModeUnknown, bufReader, nil
	}

	// Detect based on first character
	switch content[0] {
	case '[':
		// JSON array
		return ModeBatchJSON, bufReader, nil
	case '{':
		// Could be single JSON object or NDJSON
		if containsNDJSONPattern(content) {
			return ModeNDJSON, bufReader, nil
		}
		return ModeBatchJSON, bufReader, nil
	case '-':
		// Check for YAML document separator "---"
		if len(content) >= 3 && string(content[:3]) == "---" {
			// Check if there are multiple --- separators
			if bytes.Count(content, []byte("\n---")) > 0 {
				return ModeMultiYAML, bufReader, nil
			}
			return ModeBatchYAML, bufReader, nil
		}
		return ModeBatchYAML, bufReader, nil
	default:
		// Assume YAML (can start with alphanumeric keys)
		return ModeBatchYAML, bufReader, nil
	}
}

// skipWhitespace returns the slice with leading whitespace removed.
func skipWhitespace(data []byte) []byte {
	for i, b := range data {
		if !unicode.IsSpace(rune(b)) {
			return data[i:]
		}
	}
	return nil
}

// containsNDJSONPattern checks for {}\n{ pattern indicating NDJSON.
func containsNDJSONPattern(data []byte) bool {
	depth := 0
	for i, b := range data {
		switch b {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 && i+1 < len(data) {
				// Check for newline followed by another object
				rest := skipWhitespace(data[i+1:])
				if len(rest) > 0 && rest[0] == '{' {
					return true
				}
			}
		case '"':
			// Skip string content to avoid counting braces in strings
			for j := i + 1; j < len(data); j++ {
				if data[j] == '\\' && j+1 < len(data) {
					j++ // Skip escaped character
					continue
				}
				if data[j] == '"' {
					break
				}
			}
		}
	}
	return false
}

// IsTerminal checks if the file is a terminal (interactive) vs pipe.
func IsTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// HasStdinData returns true if stdin has data available (is a pipe or file).
func HasStdinData() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	// Not a terminal means it's a pipe or file
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// FormatFromInputMode converts InputMode to Format.
func FormatFromInputMode(mode InputMode) Format {
	switch mode {
	case ModeBatchJSON:
		return FormatJSON
	case ModeNDJSON:
		return FormatNDJSON
	case ModeBatchYAML, ModeMultiYAML:
		return FormatYAML
	default:
		return FormatJSON // Default to JSON for unknown
	}
}
