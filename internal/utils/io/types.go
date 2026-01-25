// Package io provides application-layer I/O handling for the CLI.
// It enforces structured JSON/YAML input/output while components remain strongly typed.
package io

import (
	"context"
	"io"
)

// Format represents an I/O format type.
type Format string

const (
	// FormatJSON outputs as pretty-printed JSON.
	FormatJSON Format = "json"
	// FormatYAML outputs as YAML.
	FormatYAML Format = "yaml"
	// FormatNDJSON outputs as newline-delimited JSON (JSON Lines).
	FormatNDJSON Format = "ndjson"
	// FormatTable outputs as human-readable table (default).
	FormatTable Format = "table"
	// FormatRaw outputs data as-is without marshaling.
	FormatRaw Format = "raw"
	// FormatAuto auto-detects format from content.
	FormatAuto Format = "auto"
)

// ParseFormat parses a string into a Format type.
func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "yaml", "yml":
		return FormatYAML
	case "ndjson", "jsonl", "jsonlines":
		return FormatNDJSON
	case "table", "":
		return FormatTable
	case "raw":
		return FormatRaw
	case "auto":
		return FormatAuto
	default:
		return FormatTable
	}
}

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// IsStructured returns true if the format is a structured data format (JSON/YAML/NDJSON).
func (f Format) IsStructured() bool {
	switch f {
	case FormatJSON, FormatYAML, FormatNDJSON:
		return true
	default:
		return false
	}
}

// IOConfig holds I/O configuration for a command.
type IOConfig struct {
	// Input configuration
	InputFormat Format    // Format of input data (auto, json, yaml, ndjson)
	InputFile   string    // File to read from (empty = stdin)
	InputReader io.Reader // Underlying reader (set by middleware)

	// Output configuration
	OutputFormat Format    // Format for output data (table, json, yaml, ndjson, raw)
	OutputFile   string    // File to write to (empty = stdout)
	OutputWriter io.Writer // Underlying writer (set by middleware)

	// Behavior flags
	Pretty    bool // Pretty-print JSON/YAML output
	Streaming bool // Enable streaming mode (NDJSON)
	NoColor   bool // Disable ANSI colors (auto-detected for non-TTY)
}

// NewIOConfig creates a new IOConfig with defaults.
func NewIOConfig() *IOConfig {
	return &IOConfig{
		InputFormat:  FormatAuto,
		OutputFormat: FormatTable,
		Pretty:       true,
		Streaming:    false,
		NoColor:      false,
	}
}

// IOContext is the central I/O state passed to commands via context.
type IOContext struct {
	Config *IOConfig
	Reader *Reader
	Writer *Writer
	Errors *ErrorHandler
	ctx    context.Context
}

// Context returns the underlying context.
func (c *IOContext) Context() context.Context {
	return c.ctx
}

// Close closes all I/O resources.
func (c *IOContext) Close() error {
	var firstErr error

	if c.Reader != nil {
		if err := c.Reader.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if c.Writer != nil {
		if err := c.Writer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// Result wraps command output with optional metadata.
type Result[T any] struct {
	Data     T               `json:"data,omitempty" yaml:"data,omitempty"`
	Metadata *ResultMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// ResultMetadata contains optional result metadata.
type ResultMetadata struct {
	Count      *int   `json:"count,omitempty" yaml:"count,omitempty"`
	Page       *int   `json:"page,omitempty" yaml:"page,omitempty"`
	TotalPages *int   `json:"total_pages,omitempty" yaml:"total_pages,omitempty"`
	Truncated  bool   `json:"truncated,omitempty" yaml:"truncated,omitempty"`
	Message    string `json:"message,omitempty" yaml:"message,omitempty"`
}

// StreamItem represents a single item in a streaming response.
type StreamItem[T any] struct {
	Index int      `json:"index" yaml:"index"`
	Item  T        `json:"item" yaml:"item"`
	Error *IOError `json:"error,omitempty" yaml:"error,omitempty"`
}

// ContextKey is the key type for context values.
type ContextKey string

const (
	// IOContextKey is the context key for IOContext.
	IOContextKey ContextKey = "io_context"
)
