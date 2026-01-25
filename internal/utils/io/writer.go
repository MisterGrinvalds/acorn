package io

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Writer handles output serialization with streaming support.
type Writer struct {
	config   *IOConfig
	writer   io.Writer
	buffered *bufio.Writer
	mu       sync.Mutex
	closed   bool
	file     *os.File // If we opened a file, track it for closing

	// For JSON array streaming
	arrayStarted bool
	itemCount    int
}

// NewWriter creates a Writer from IOConfig.
func NewWriter(cfg *IOConfig) (*Writer, error) {
	var writer io.Writer

	if cfg.OutputFile != "" {
		f, err := os.Create(cfg.OutputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}
		writer = f
		return &Writer{
			config:   cfg,
			writer:   f,
			buffered: bufio.NewWriterSize(f, 64*1024),
			file:     f,
		}, nil
	}

	if cfg.OutputWriter != nil {
		writer = cfg.OutputWriter
	} else {
		writer = os.Stdout
	}

	return &Writer{
		config:   cfg,
		writer:   writer,
		buffered: bufio.NewWriterSize(writer, 64*1024),
	}, nil
}

// Write serializes and writes data in the configured format.
func (w *Writer) Write(data interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return fmt.Errorf("writer is closed")
	}

	switch w.config.OutputFormat {
	case FormatJSON:
		return w.writeJSON(data)
	case FormatYAML:
		return w.writeYAML(data)
	case FormatNDJSON:
		return w.writeNDJSON(data)
	case FormatRaw:
		return w.writeRaw(data)
	case FormatTable:
		return fmt.Errorf("table format must be handled by command")
	default:
		return w.writeJSON(data) // Default to JSON
	}
}

// WriteStream writes items one at a time in streaming mode.
// For NDJSON: writes each item on a separate line.
// For JSON: writes as array elements (call StartArray/EndArray for proper formatting).
func (w *Writer) WriteStream(item interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return fmt.Errorf("writer is closed")
	}

	switch w.config.OutputFormat {
	case FormatNDJSON:
		// Write as single line JSON
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		if _, err := w.buffered.Write(data); err != nil {
			return err
		}
		if _, err := w.buffered.WriteString("\n"); err != nil {
			return err
		}
		return w.buffered.Flush()

	case FormatJSON:
		// Write as array element
		if !w.arrayStarted {
			if _, err := w.buffered.WriteString("[\n"); err != nil {
				return err
			}
			w.arrayStarted = true
		} else {
			if _, err := w.buffered.WriteString(",\n"); err != nil {
				return err
			}
		}

		// Indent the item
		data, err := json.MarshalIndent(item, "  ", "  ")
		if err != nil {
			return err
		}
		if _, err := w.buffered.WriteString("  "); err != nil {
			return err
		}
		if _, err := w.buffered.Write(data); err != nil {
			return err
		}
		w.itemCount++
		return w.buffered.Flush()

	case FormatYAML:
		// Write as YAML document with separator
		if w.itemCount > 0 {
			if _, err := w.buffered.WriteString("---\n"); err != nil {
				return err
			}
		}

		enc := yaml.NewEncoder(w.buffered)
		enc.SetIndent(2)
		if err := enc.Encode(item); err != nil {
			return err
		}
		if err := enc.Close(); err != nil {
			return err
		}
		w.itemCount++
		return w.buffered.Flush()

	default:
		return w.writeJSON(item)
	}
}

// writeJSON writes pretty or compact JSON.
func (w *Writer) writeJSON(data interface{}) error {
	enc := json.NewEncoder(w.buffered)
	if w.config.Pretty {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(data); err != nil {
		return err
	}
	return w.buffered.Flush()
}

// writeYAML writes YAML output.
func (w *Writer) writeYAML(data interface{}) error {
	enc := yaml.NewEncoder(w.buffered)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}
	return w.buffered.Flush()
}

// writeNDJSON writes data as NDJSON (handles arrays specially).
func (w *Writer) writeNDJSON(data interface{}) error {
	// If data is a slice, write each element on a separate line
	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			itemData, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if _, err := w.buffered.Write(itemData); err != nil {
				return err
			}
			if _, err := w.buffered.WriteString("\n"); err != nil {
				return err
			}
		}
		return w.buffered.Flush()
	default:
		// Single item
		itemData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if _, err := w.buffered.Write(itemData); err != nil {
			return err
		}
		if _, err := w.buffered.WriteString("\n"); err != nil {
			return err
		}
		return w.buffered.Flush()
	}
}

// writeRaw writes data as-is (for []byte or string).
func (w *Writer) writeRaw(data interface{}) error {
	switch v := data.(type) {
	case []byte:
		if _, err := w.buffered.Write(v); err != nil {
			return err
		}
	case string:
		if _, err := w.buffered.WriteString(v); err != nil {
			return err
		}
	case fmt.Stringer:
		if _, err := w.buffered.WriteString(v.String()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("raw format requires []byte, string, or fmt.Stringer, got %T", data)
	}
	return w.buffered.Flush()
}

// Flush flushes any buffered data.
func (w *Writer) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffered.Flush()
}

// Close closes the writer and underlying file if applicable.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil
	}
	w.closed = true

	// Finalize JSON array if started
	if w.arrayStarted && w.config.OutputFormat == FormatJSON {
		if _, err := w.buffered.WriteString("\n]\n"); err != nil {
			return err
		}
	}

	if err := w.buffered.Flush(); err != nil {
		return err
	}

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// IsStructured returns true if output format is structured (JSON/YAML/NDJSON).
func (w *Writer) IsStructured() bool {
	return w.config.OutputFormat.IsStructured()
}

// Format returns the configured output format.
func (w *Writer) Format() Format {
	return w.config.OutputFormat
}

// Printf writes formatted text (for non-structured output).
func (w *Writer) Printf(format string, args ...interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := fmt.Fprintf(w.buffered, format, args...); err != nil {
		return err
	}
	return w.buffered.Flush()
}

// Println writes a line of text (for non-structured output).
func (w *Writer) Println(args ...interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := fmt.Fprintln(w.buffered, args...); err != nil {
		return err
	}
	return w.buffered.Flush()
}

// Underlying returns the underlying io.Writer for direct access.
// Use with caution - prefer Write/WriteStream for structured output.
func (w *Writer) Underlying() io.Writer {
	return w.writer
}
