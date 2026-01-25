package io

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// AddOutputFlag adds just the output format flag to a command.
// Use this for subcommand groups that need their own output flag.
func AddOutputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("output", "o", "table",
		"Output format (table|json|yaml)")
}

// BindFlags adds I/O flags to a command (typically root command).
// These flags are inherited by all subcommands.
func BindFlags(cmd *cobra.Command, cfg *IOConfig) {
	// Output flags
	cmd.PersistentFlags().StringVarP((*string)(&cfg.OutputFormat), "output", "o", "table",
		"Output format (table|json|yaml|ndjson|raw)")
	cmd.PersistentFlags().StringVar(&cfg.OutputFile, "output-file", "",
		"Write output to file instead of stdout")

	// Input flags
	cmd.PersistentFlags().StringVarP((*string)(&cfg.InputFormat), "input-format", "I", "auto",
		"Input format (json|yaml|ndjson|auto)")
	cmd.PersistentFlags().StringVarP(&cfg.InputFile, "input-file", "f", "",
		"Read input from file instead of stdin")

	// Behavior flags
	cmd.PersistentFlags().BoolVar(&cfg.Pretty, "pretty", true,
		"Pretty-print JSON/YAML output")
	cmd.PersistentFlags().BoolVar(&cfg.Streaming, "streaming", false,
		"Enable streaming mode (NDJSON for output)")
	cmd.PersistentFlags().BoolVar(&cfg.NoColor, "no-color", false,
		"Disable ANSI color codes in output")
}

// Middleware returns Cobra PersistentPreRunE and PersistentPostRunE functions
// that set up and tear down the IOContext.
func Middleware(cfg *IOConfig) (preRun, postRun func(*cobra.Command, []string) error) {
	var ioCtx *IOContext

	preRun = func(cmd *cobra.Command, args []string) error {
		// Auto-detect NoColor for non-TTY
		if !cfg.NoColor && cfg.OutputFile == "" {
			if !term.IsTerminal(int(os.Stdout.Fd())) {
				cfg.NoColor = true
			}
		}

		// Auto-switch to JSON for non-TTY if table format
		if cfg.OutputFormat == FormatTable && cfg.OutputFile == "" {
			if !term.IsTerminal(int(os.Stdout.Fd())) {
				cfg.OutputFormat = FormatJSON
			}
		}

		// Enable streaming output for NDJSON
		if cfg.OutputFormat == FormatNDJSON {
			cfg.Streaming = true
		}

		// Create IOContext
		var err error
		ioCtx, err = NewIOContext(cmd.Context(), cfg)
		if err != nil {
			return err
		}

		// Store in command context
		ctx := context.WithValue(cmd.Context(), IOContextKey, ioCtx)
		cmd.SetContext(ctx)

		return nil
	}

	postRun = func(cmd *cobra.Command, args []string) error {
		if ioCtx != nil {
			return ioCtx.Close()
		}
		return nil
	}

	return preRun, postRun
}

// NewIOContext creates a fully initialized IOContext.
func NewIOContext(ctx context.Context, cfg *IOConfig) (*IOContext, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	reader, err := NewReader(cfg)
	if err != nil {
		return nil, err
	}

	writer, err := NewWriter(cfg)
	if err != nil {
		reader.Close()
		return nil, err
	}

	// Error handler writes to stderr by default, or to the output if structured
	errWriter := os.Stderr
	if cfg.OutputFormat.IsStructured() {
		errWriter = os.Stderr // Keep errors on stderr even for structured output
	}

	return &IOContext{
		Config: cfg,
		Reader: reader,
		Writer: writer,
		Errors: NewErrorHandler(cfg, errWriter),
		ctx:    ctx,
	}, nil
}

// GetIOContext retrieves IOContext from Cobra command context.
func GetIOContext(cmd *cobra.Command) *IOContext {
	if ctx := cmd.Context(); ctx != nil {
		if ioCtx, ok := ctx.Value(IOContextKey).(*IOContext); ok {
			return ioCtx
		}
	}
	return nil
}

// MustGetIOContext retrieves IOContext or panics.
func MustGetIOContext(cmd *cobra.Command) *IOContext {
	ioCtx := GetIOContext(cmd)
	if ioCtx == nil {
		panic("IOContext not found in command context - did you add I/O middleware?")
	}
	return ioCtx
}

// ChainPreRunE chains multiple PersistentPreRunE functions.
// Use this if you have existing PreRunE functions to preserve.
func ChainPreRunE(funcs ...func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range funcs {
			if fn != nil {
				if err := fn(cmd, args); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// ChainPostRunE chains multiple PersistentPostRunE functions.
// Use this if you have existing PostRunE functions to preserve.
func ChainPostRunE(funcs ...func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var firstErr error
		for _, fn := range funcs {
			if fn != nil {
				if err := fn(cmd, args); err != nil && firstErr == nil {
					firstErr = err
				}
			}
		}
		return firstErr
	}
}
