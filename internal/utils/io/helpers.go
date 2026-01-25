package io

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// CommandIO provides convenient I/O operations for command implementations.
type CommandIO struct {
	ctx *IOContext
}

// IO returns a CommandIO helper from the command.
// If IOContext is not available, returns a helper that works with defaults.
func IO(cmd *cobra.Command) *CommandIO {
	return &CommandIO{ctx: GetIOContext(cmd)}
}

// ReadInput reads and deserializes input into target.
func (c *CommandIO) ReadInput(target any) error {
	if c.ctx == nil || c.ctx.Reader == nil {
		return NewError(ErrCodeInternal, "I/O context not initialized")
	}
	return c.ctx.Reader.Read(target)
}

// ReadRaw reads the input as raw bytes without deserialization.
func (c *CommandIO) ReadRaw() ([]byte, error) {
	if c.ctx == nil || c.ctx.Reader == nil {
		return nil, NewError(ErrCodeInternal, "I/O context not initialized")
	}
	return c.ctx.Reader.ReadRaw()
}

// ReadStream returns a channel that yields input items one at a time.
func (c *CommandIO) ReadStream() (<-chan *StreamRecord, error) {
	if c.ctx == nil || c.ctx.Reader == nil {
		return nil, NewError(ErrCodeInternal, "I/O context not initialized")
	}
	return c.ctx.Reader.ReadStream()
}

// WriteOutput serializes and writes output.
func (c *CommandIO) WriteOutput(data any) error {
	if c.ctx == nil || c.ctx.Writer == nil {
		return NewError(ErrCodeInternal, "I/O context not initialized")
	}
	return c.ctx.Writer.Write(data)
}

// WriteStreamItem writes a single item in streaming mode.
func (c *CommandIO) WriteStreamItem(item any) error {
	if c.ctx == nil || c.ctx.Writer == nil {
		return NewError(ErrCodeInternal, "I/O context not initialized")
	}
	return c.ctx.Writer.WriteStream(item)
}

// Error handles an error with structured output.
// Returns the error for use as command return value.
func (c *CommandIO) Error(err error) error {
	if c.ctx == nil || c.ctx.Errors == nil {
		return err
	}
	return c.ctx.Errors.HandleError(err)
}

// IsStructured returns true if output is structured (JSON/YAML/NDJSON).
func (c *CommandIO) IsStructured() bool {
	if c.ctx == nil || c.ctx.Writer == nil {
		return false
	}
	return c.ctx.Writer.IsStructured()
}

// IsTable returns true if output format is table (human-readable).
func (c *CommandIO) IsTable() bool {
	if c.ctx == nil {
		return true
	}
	return c.ctx.Config.OutputFormat == FormatTable
}

// HasInput returns true if there's input available.
func (c *CommandIO) HasInput() bool {
	if c.ctx == nil || c.ctx.Reader == nil {
		return false
	}
	return c.ctx.Reader.HasInput()
}

// Format returns the configured output format.
func (c *CommandIO) Format() Format {
	if c.ctx == nil {
		return FormatTable
	}
	return c.ctx.Config.OutputFormat
}

// InputFormat returns the configured input format.
func (c *CommandIO) InputFormat() Format {
	if c.ctx == nil {
		return FormatJSON
	}
	return c.ctx.Config.InputFormat
}

// Config returns the I/O configuration.
func (c *CommandIO) Config() *IOConfig {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Config
}

// Context returns the IOContext.
func (c *CommandIO) Context() *IOContext {
	return c.ctx
}

// Printf writes formatted text output.
// Only use for non-structured output (table format).
func (c *CommandIO) Printf(format string, args ...any) error {
	if c.ctx == nil || c.ctx.Writer == nil {
		_, err := fmt.Fprintf(os.Stdout, format, args...)
		return err
	}
	return c.ctx.Writer.Printf(format, args...)
}

// Println writes a line of text output.
// Only use for non-structured output (table format).
func (c *CommandIO) Println(args ...any) error {
	if c.ctx == nil || c.ctx.Writer == nil {
		_, err := fmt.Fprintln(os.Stdout, args...)
		return err
	}
	return c.ctx.Writer.Println(args...)
}

// Writer returns the underlying io.Writer for direct access.
// Use with caution - prefer WriteOutput for structured output.
func (c *CommandIO) Writer() io.Writer {
	if c.ctx == nil || c.ctx.Writer == nil {
		return os.Stdout
	}
	return c.ctx.Writer.Underlying()
}

// Stderr returns os.Stderr for writing error messages.
func (c *CommandIO) Stderr() io.Writer {
	return os.Stderr
}

// NoColor returns true if color output should be disabled.
func (c *CommandIO) NoColor() bool {
	if c.ctx == nil {
		return false
	}
	return c.ctx.Config.NoColor
}

// Flush flushes any buffered output.
func (c *CommandIO) Flush() error {
	if c.ctx == nil || c.ctx.Writer == nil {
		return nil
	}
	return c.ctx.Writer.Flush()
}
