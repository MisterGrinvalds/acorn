// Package yq provides yq YAML processing helper functionality.
package yq

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents yq installation status.
type Status struct {
	Installed bool   `json:"installed" yaml:"installed"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
	Location  string `json:"location,omitempty" yaml:"location,omitempty"`
}

// FilterResult represents the result of a yq filter operation.
type FilterResult struct {
	Output string `json:"output" yaml:"output"`
	Valid  bool   `json:"valid" yaml:"valid"`
	Error  string `json:"error,omitempty" yaml:"error,omitempty"`
}

// ConvertFormat represents supported conversion formats.
type ConvertFormat string

const (
	FormatYAML  ConvertFormat = "yaml"
	FormatJSON  ConvertFormat = "json"
	FormatXML   ConvertFormat = "xml"
	FormatTOML  ConvertFormat = "toml"
	FormatCSV   ConvertFormat = "csv"
	FormatTSV   ConvertFormat = "tsv"
	FormatProps ConvertFormat = "props"
)

// Helper provides yq helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new yq Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// IsInstalled checks if yq is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("yq")
	return err == nil
}

// GetStatus returns yq installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	yqPath, err := exec.LookPath("yq")
	if err != nil {
		return status
	}

	status.Installed = true
	status.Location = yqPath

	// Get version
	cmd := exec.Command("yq", "--version")
	out, err := cmd.Output()
	if err == nil {
		// yq version format: "yq (https://github.com/mikefarah/yq/) version v4.x.x"
		version := strings.TrimSpace(string(out))
		if idx := strings.Index(version, "version"); idx != -1 {
			status.Version = strings.TrimSpace(version[idx+7:])
		} else {
			status.Version = version
		}
	}

	return status
}

// Filter applies a yq expression to YAML input.
func (h *Helper) Filter(input []byte, expression string) (*FilterResult, error) {
	cmd := exec.Command("yq", expression)
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

// FilterFile applies a yq expression to a YAML file.
func (h *Helper) FilterFile(filepath, expression string) (*FilterResult, error) {
	cmd := exec.Command("yq", expression, filepath)

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

// EditInPlace modifies a YAML file in place.
func (h *Helper) EditInPlace(filepath, expression string) error {
	cmd := exec.Command("yq", "-i", expression, filepath)
	return cmd.Run()
}

// Validate checks if input is valid YAML.
func (h *Helper) Validate(input []byte) bool {
	cmd := exec.Command("yq", ".")
	cmd.Stdin = bytes.NewReader(input)
	return cmd.Run() == nil
}

// ValidateFile checks if a file contains valid YAML.
func (h *Helper) ValidateFile(filepath string) (bool, error) {
	cmd := exec.Command("yq", ".", filepath)
	err := cmd.Run()
	return err == nil, err
}

// Format pretty-prints YAML.
func (h *Helper) Format(input []byte) (*FilterResult, error) {
	return h.Filter(input, ".")
}

// Convert converts between formats.
func (h *Helper) Convert(input []byte, inputFormat, outputFormat ConvertFormat) (*FilterResult, error) {
	args := []string{}

	// Input format
	if inputFormat != FormatYAML {
		args = append(args, "-p="+string(inputFormat))
	}

	// Output format
	if outputFormat != FormatYAML {
		args = append(args, "-o="+string(outputFormat))
	}

	args = append(args, ".")

	cmd := exec.Command("yq", args...)
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

// ConvertFile converts a file between formats.
func (h *Helper) ConvertFile(filepath string, inputFormat, outputFormat ConvertFormat) (*FilterResult, error) {
	args := []string{}

	// Input format
	if inputFormat != FormatYAML {
		args = append(args, "-p="+string(inputFormat))
	}

	// Output format
	if outputFormat != FormatYAML {
		args = append(args, "-o="+string(outputFormat))
	}

	args = append(args, ".", filepath)

	cmd := exec.Command("yq", args...)

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

// YAMLToJSON converts YAML to JSON.
func (h *Helper) YAMLToJSON(input []byte) (*FilterResult, error) {
	return h.Convert(input, FormatYAML, FormatJSON)
}

// JSONToYAML converts JSON to YAML.
func (h *Helper) JSONToYAML(input []byte) (*FilterResult, error) {
	return h.Convert(input, FormatJSON, FormatYAML)
}

// YAMLToXML converts YAML to XML.
func (h *Helper) YAMLToXML(input []byte) (*FilterResult, error) {
	return h.Convert(input, FormatYAML, FormatXML)
}

// XMLToYAML converts XML to YAML.
func (h *Helper) XMLToYAML(input []byte) (*FilterResult, error) {
	return h.Convert(input, FormatXML, FormatYAML)
}

// YAMLToTOML converts YAML to TOML.
func (h *Helper) YAMLToTOML(input []byte) (*FilterResult, error) {
	return h.Convert(input, FormatYAML, FormatTOML)
}

// TOMLToYAML converts TOML to YAML.
func (h *Helper) TOMLToYAML(input []byte) (*FilterResult, error) {
	return h.Convert(input, FormatTOML, FormatYAML)
}

// Merge merges YAML files/inputs using the * operator.
func (h *Helper) Merge(base, overlay []byte) (*FilterResult, error) {
	// Write overlay to temp file since yq load() needs a file path
	tmpFile, err := os.CreateTemp("", "yq-overlay-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(overlay); err != nil {
		return nil, fmt.Errorf("failed to write overlay: %w", err)
	}
	tmpFile.Close()

	expression := fmt.Sprintf(`. * load("%s")`, tmpFile.Name())
	return h.Filter(base, expression)
}

// MergeDeep performs deep merge with array append (*+).
func (h *Helper) MergeDeep(base, overlay []byte) (*FilterResult, error) {
	tmpFile, err := os.CreateTemp("", "yq-overlay-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(overlay); err != nil {
		return nil, fmt.Errorf("failed to write overlay: %w", err)
	}
	tmpFile.Close()

	expression := fmt.Sprintf(`. *+ load("%s")`, tmpFile.Name())
	return h.Filter(base, expression)
}

// MergeFiles merges multiple YAML files.
func (h *Helper) MergeFiles(files ...string) (*FilterResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	// Read first file as base
	base, err := os.ReadFile(files[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", files[0], err)
	}

	result := &FilterResult{
		Output: string(base),
		Valid:  true,
	}

	// Merge each subsequent file
	for i := 1; i < len(files); i++ {
		overlay, err := os.ReadFile(files[i])
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", files[i], err)
		}

		result, err = h.Merge([]byte(result.Output), overlay)
		if err != nil {
			return nil, fmt.Errorf("failed to merge %s: %w", files[i], err)
		}
		if !result.Valid {
			return result, nil
		}
	}

	return result, nil
}

// Get extracts a value at a path.
func (h *Helper) Get(input []byte, path string) (*FilterResult, error) {
	return h.Filter(input, path)
}

// Set sets a value at a path.
func (h *Helper) Set(input []byte, path, value string) (*FilterResult, error) {
	expression := fmt.Sprintf("%s = %s", path, value)
	return h.Filter(input, expression)
}

// SetString sets a string value at a path.
func (h *Helper) SetString(input []byte, path, value string) (*FilterResult, error) {
	expression := fmt.Sprintf(`%s = "%s"`, path, value)
	return h.Filter(input, expression)
}

// Delete removes a key from YAML.
func (h *Helper) Delete(input []byte, path string) (*FilterResult, error) {
	expression := fmt.Sprintf("del(%s)", path)
	return h.Filter(input, expression)
}

// Select filters elements matching a condition.
func (h *Helper) Select(input []byte, condition string) (*FilterResult, error) {
	expression := fmt.Sprintf(".[] | select(%s)", condition)
	return h.Filter(input, expression)
}

// EvalAll processes multiple documents with eval-all.
func (h *Helper) EvalAll(input []byte, expression string) (*FilterResult, error) {
	cmd := exec.Command("yq", "eval-all", expression)
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

// CountDocuments counts the number of documents in multi-doc YAML.
func (h *Helper) CountDocuments(input []byte) (int, error) {
	result, err := h.EvalAll(input, "[.] | length")
	if err != nil {
		return 0, err
	}
	if !result.Valid {
		return 0, fmt.Errorf("failed to count documents: %s", result.Error)
	}

	var count int
	if _, err := fmt.Sscanf(strings.TrimSpace(result.Output), "%d", &count); err != nil {
		return 0, err
	}

	return count, nil
}

// GetDocument extracts a specific document from multi-doc YAML.
func (h *Helper) GetDocument(input []byte, index int) (*FilterResult, error) {
	expression := fmt.Sprintf("select(documentIndex == %d)", index)
	return h.Filter(input, expression)
}

// SplitDocuments splits multi-doc YAML into separate documents.
func (h *Helper) SplitDocuments(input []byte) ([]string, error) {
	result, err := h.EvalAll(input, ".")
	if err != nil {
		return nil, err
	}
	if !result.Valid {
		return nil, fmt.Errorf("failed to split documents: %s", result.Error)
	}

	// Split on document separator
	docs := strings.Split(result.Output, "---\n")
	var filtered []string
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc != "" {
			filtered = append(filtered, doc)
		}
	}

	return filtered, nil
}

// GetKeys returns all keys from a YAML object.
func (h *Helper) GetKeys(input []byte) ([]string, error) {
	result, err := h.Filter(input, "keys | .[]")
	if err != nil {
		return nil, err
	}
	if !result.Valid {
		return nil, fmt.Errorf("failed to get keys: %s", result.Error)
	}

	keys := strings.Split(strings.TrimSpace(result.Output), "\n")
	return keys, nil
}

// Comment adds a comment to a path.
func (h *Helper) Comment(input []byte, path, comment string) (*FilterResult, error) {
	expression := fmt.Sprintf(`%s line_comment="%s"`, path, comment)
	return h.Filter(input, expression)
}

// Sort sorts YAML keys alphabetically.
func (h *Helper) Sort(input []byte) (*FilterResult, error) {
	return h.Filter(input, "sort_keys(.)")
}

// Flatten flattens nested YAML to dot notation.
func (h *Helper) Flatten(input []byte) (*FilterResult, error) {
	return h.Filter(input, `[.. | select(type == "string" or type == "number" or type == "boolean") | {path: path | join("."), value: .}]`)
}

// EnvSubst substitutes environment variables in YAML values.
func (h *Helper) EnvSubst(input []byte, envVars map[string]string) (*FilterResult, error) {
	// Build yq expression with strenv
	var parts []string
	for k, v := range envVars {
		os.Setenv(k, v)
		parts = append(parts, fmt.Sprintf(`.. | select(. == "$%s") |= strenv(%s)`, k, k))
	}

	if len(parts) == 0 {
		return h.Filter(input, ".")
	}

	expression := strings.Join(parts, " | ")
	return h.Filter(input, expression)
}
