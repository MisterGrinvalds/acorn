// Package jq provides jq JSON processing helper functionality.
package jq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents jq installation status.
type Status struct {
	Installed bool   `json:"installed" yaml:"installed"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
	Location  string `json:"location,omitempty" yaml:"location,omitempty"`
}

// FilterResult represents the result of a jq filter operation.
type FilterResult struct {
	Output string `json:"output" yaml:"output"`
	Valid  bool   `json:"valid" yaml:"valid"`
	Error  string `json:"error,omitempty" yaml:"error,omitempty"`
}

// Helper provides jq helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new jq Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// IsInstalled checks if jq is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("jq")
	return err == nil
}

// GetStatus returns jq installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	jqPath, err := exec.LookPath("jq")
	if err != nil {
		return status
	}

	status.Installed = true
	status.Location = jqPath

	// Get version
	cmd := exec.Command("jq", "--version")
	out, err := cmd.Output()
	if err == nil {
		status.Version = strings.TrimSpace(string(out))
	}

	return status
}

// Filter applies a jq filter to JSON input.
func (h *Helper) Filter(input []byte, filter string) (*FilterResult, error) {
	cmd := exec.Command("jq", filter)
	cmd.Stdin = bytes.NewReader(input)

	out, err := cmd.Output()
	if err != nil {
		return &FilterResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &FilterResult{
		Output: string(out),
		Valid:  true,
	}, nil
}

// FilterCompact applies a jq filter and returns compact output.
func (h *Helper) FilterCompact(input []byte, filter string) (*FilterResult, error) {
	cmd := exec.Command("jq", "-c", filter)
	cmd.Stdin = bytes.NewReader(input)

	out, err := cmd.Output()
	if err != nil {
		return &FilterResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &FilterResult{
		Output: string(out),
		Valid:  true,
	}, nil
}

// FilterRaw applies a jq filter with raw string output.
func (h *Helper) FilterRaw(input []byte, filter string) (*FilterResult, error) {
	cmd := exec.Command("jq", "-r", filter)
	cmd.Stdin = bytes.NewReader(input)

	out, err := cmd.Output()
	if err != nil {
		return &FilterResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &FilterResult{
		Output: string(out),
		Valid:  true,
	}, nil
}

// FilterFile applies a jq filter to a JSON file.
func (h *Helper) FilterFile(filepath, filter string) (*FilterResult, error) {
	cmd := exec.Command("jq", filter, filepath)

	out, err := cmd.Output()
	if err != nil {
		return &FilterResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &FilterResult{
		Output: string(out),
		Valid:  true,
	}, nil
}

// Validate checks if input is valid JSON.
func (h *Helper) Validate(input []byte) bool {
	return json.Valid(input)
}

// ValidateFile checks if a file contains valid JSON.
func (h *Helper) ValidateFile(filepath string) (bool, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}
	return json.Valid(data), nil
}

// Format pretty-prints JSON.
func (h *Helper) Format(input []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return json.MarshalIndent(data, "", "  ")
}

// Compact compresses JSON to a single line.
func (h *Helper) Compact(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, input); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return buf.Bytes(), nil
}

// GetKeys returns all keys from a JSON object.
func (h *Helper) GetKeys(input []byte) ([]string, error) {
	result, err := h.FilterRaw(input, "keys[]")
	if err != nil {
		return nil, err
	}
	if !result.Valid {
		return nil, fmt.Errorf("failed to get keys: %s", result.Error)
	}

	keys := strings.Split(strings.TrimSpace(result.Output), "\n")
	return keys, nil
}

// GetLength returns the length of an array or object.
func (h *Helper) GetLength(input []byte) (int, error) {
	result, err := h.FilterRaw(input, "length")
	if err != nil {
		return 0, err
	}
	if !result.Valid {
		return 0, fmt.Errorf("failed to get length: %s", result.Error)
	}

	var length int
	if _, err := fmt.Sscanf(strings.TrimSpace(result.Output), "%d", &length); err != nil {
		return 0, err
	}

	return length, nil
}

// GetType returns the type of a JSON value.
func (h *Helper) GetType(input []byte) (string, error) {
	result, err := h.FilterRaw(input, "type")
	if err != nil {
		return "", err
	}
	if !result.Valid {
		return "", fmt.Errorf("failed to get type: %s", result.Error)
	}

	return strings.TrimSpace(result.Output), nil
}

// Flatten flattens nested JSON structure.
func (h *Helper) Flatten(input []byte) (*FilterResult, error) {
	return h.Filter(input, "[paths(scalars) as $p | {path: $p | join(\".\"), value: getpath($p)}]")
}

// Select filters array elements matching a condition.
func (h *Helper) Select(input []byte, condition string) (*FilterResult, error) {
	filter := fmt.Sprintf(".[] | select(%s)", condition)
	return h.Filter(input, filter)
}

// Map applies a transformation to each array element.
func (h *Helper) Map(input []byte, transform string) (*FilterResult, error) {
	filter := fmt.Sprintf("[.[] | %s]", transform)
	return h.Filter(input, filter)
}

// Sort sorts an array.
func (h *Helper) Sort(input []byte) (*FilterResult, error) {
	return h.Filter(input, "sort")
}

// SortBy sorts an array by a key.
func (h *Helper) SortBy(input []byte, key string) (*FilterResult, error) {
	filter := fmt.Sprintf("sort_by(%s)", key)
	return h.Filter(input, filter)
}

// Unique removes duplicates from an array.
func (h *Helper) Unique(input []byte) (*FilterResult, error) {
	return h.Filter(input, "unique")
}

// Group groups array elements by a key.
func (h *Helper) Group(input []byte, key string) (*FilterResult, error) {
	filter := fmt.Sprintf("group_by(%s)", key)
	return h.Filter(input, filter)
}

// Slurp reads multiple JSON values into an array.
func (h *Helper) Slurp(inputs ...[]byte) (*FilterResult, error) {
	var combined bytes.Buffer
	for _, input := range inputs {
		combined.Write(input)
		combined.WriteByte('\n')
	}

	cmd := exec.Command("jq", "-s", ".")
	cmd.Stdin = &combined

	out, err := cmd.Output()
	if err != nil {
		return &FilterResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &FilterResult{
		Output: string(out),
		Valid:  true,
	}, nil
}

// ToCSV converts JSON array to CSV format.
func (h *Helper) ToCSV(input []byte) (*FilterResult, error) {
	// Get headers from first object
	filter := `(.[0] | keys_unsorted) as $keys | $keys, (.[] | [.[$keys[]]] | @csv) | @csv`
	result, err := h.FilterRaw(input, filter)
	if err != nil {
		return nil, err
	}

	// Simpler approach if that fails
	if !result.Valid {
		return h.FilterRaw(input, `.[] | [.[]] | @csv`)
	}

	return result, nil
}

// FromCSV converts CSV to JSON (requires explicit headers).
func (h *Helper) FromCSV(input []byte, headers []string) (*FilterResult, error) {
	// This is a simplified implementation
	// Full CSV parsing would require more complex jq expressions
	filter := fmt.Sprintf(`split("\n") | .[1:] | map(split(",") | {%s})`, buildCSVMapping(headers))
	return h.Filter(input, filter)
}

// buildCSVMapping builds a jq object mapping for CSV headers.
func buildCSVMapping(headers []string) string {
	var parts []string
	for i, h := range headers {
		parts = append(parts, fmt.Sprintf(`"%s": .[%d]`, h, i))
	}
	return strings.Join(parts, ", ")
}

// Merge merges multiple JSON objects.
func (h *Helper) Merge(inputs ...[]byte) (*FilterResult, error) {
	// Combine inputs for slurping
	var combined bytes.Buffer
	for _, input := range inputs {
		combined.Write(input)
		combined.WriteByte('\n')
	}

	cmd := exec.Command("jq", "-s", "add")
	cmd.Stdin = &combined

	out, err := cmd.Output()
	if err != nil {
		return &FilterResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &FilterResult{
		Output: string(out),
		Valid:  true,
	}, nil
}

// Diff compares two JSON values and returns differences.
func (h *Helper) Diff(a, b []byte) (*FilterResult, error) {
	// Combine both inputs
	var combined bytes.Buffer
	combined.Write(a)
	combined.WriteByte('\n')
	combined.Write(b)

	filter := `-s 'def diff($a; $b):
		($a | type) as $at | ($b | type) as $bt |
		if $at != $bt then {type_mismatch: {a: $at, b: $bt}}
		elif $at == "object" then
			(($a | keys) + ($b | keys) | unique) as $keys |
			[$keys[] as $k | {key: $k, diff: diff($a[$k]; $b[$k])} | select(.diff != null)] |
			if length == 0 then null else from_entries end
		elif $at == "array" then
			if $a == $b then null else {a: $a, b: $b} end
		else
			if $a == $b then null else {a: $a, b: $b} end
		end;
	diff(.[0]; .[1])'`

	cmd := exec.Command("jq", filter)
	cmd.Stdin = &combined

	out, err := cmd.Output()
	if err != nil {
		// Fallback to simple comparison
		return h.Filter(combined.Bytes(), `-s '.[0] == .[1]'`)
	}

	return &FilterResult{
		Output: string(out),
		Valid:  true,
	}, nil
}

// Path extracts a value at a specific path.
func (h *Helper) Path(input []byte, path string) (*FilterResult, error) {
	return h.Filter(input, path)
}

// Set sets a value at a specific path.
func (h *Helper) Set(input []byte, path, value string) (*FilterResult, error) {
	filter := fmt.Sprintf("%s = %s", path, value)
	return h.Filter(input, filter)
}

// Delete removes a key from JSON.
func (h *Helper) Delete(input []byte, path string) (*FilterResult, error) {
	filter := fmt.Sprintf("del(%s)", path)
	return h.Filter(input, filter)
}
