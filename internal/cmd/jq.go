package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/data/jq"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	jqOutputFormat string
	jqVerbose      bool
	jqCompact      bool
	jqRaw          bool
	jqSlurp        bool
)

// jqCmd represents the jq command group
var jqCmd = &cobra.Command{
	Use:   "jq",
	Short: "JSON processing commands",
	Long: `JSON processing commands using jq.

Provides filtering, transformation, and validation of JSON data.

Examples:
  acorn jq status                    # Show jq status
  acorn jq filter '.name' file.json  # Extract field
  acorn jq validate file.json        # Validate JSON
  acorn jq format file.json          # Pretty print
  acorn jq keys file.json            # Get keys`,
}

// jqStatusCmd shows jq status
var jqStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show jq installation status",
	Long: `Display jq installation status and version.

Examples:
  acorn jq status
  acorn jq status -o json`,
	RunE: runJqStatus,
}

// jqFilterCmd applies a filter
var jqFilterCmd = &cobra.Command{
	Use:   "filter [expression] [file]",
	Short: "Apply jq filter to JSON",
	Long: `Apply a jq filter expression to JSON input.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn jq filter '.name' file.json
  acorn jq filter '.items[]' file.json
  cat file.json | acorn jq filter '.name'
  acorn jq filter -c '.name' file.json  # compact output
  acorn jq filter -r '.name' file.json  # raw string output`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runJqFilter,
}

// jqValidateCmd validates JSON
var jqValidateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate JSON",
	Long: `Check if input is valid JSON.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn jq validate file.json
  cat file.json | acorn jq validate`,
	Args: cobra.MaximumNArgs(1),
	RunE: runJqValidate,
}

// jqFormatCmd formats JSON
var jqFormatCmd = &cobra.Command{
	Use:   "format [file]",
	Short: "Pretty print JSON",
	Long: `Format JSON with proper indentation.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn jq format file.json
  cat file.json | acorn jq format`,
	Aliases: []string{"fmt", "pretty"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runJqFormat,
}

// jqCompactCmd compacts JSON
var jqCompactCmd = &cobra.Command{
	Use:   "compact [file]",
	Short: "Compact JSON to single line",
	Long: `Remove whitespace and compact JSON to a single line.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn jq compact file.json
  cat file.json | acorn jq compact`,
	Aliases: []string{"minify"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runJqCompact,
}

// jqKeysCmd gets keys
var jqKeysCmd = &cobra.Command{
	Use:   "keys [file]",
	Short: "Get keys from JSON object",
	Long: `Extract all keys from a JSON object.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn jq keys file.json
  cat file.json | acorn jq keys`,
	Args: cobra.MaximumNArgs(1),
	RunE: runJqKeys,
}

// jqTypeCmd gets type
var jqTypeCmd = &cobra.Command{
	Use:   "type [file]",
	Short: "Get JSON value type",
	Long: `Get the type of a JSON value (object, array, string, number, boolean, null).

Reads from file if provided, otherwise from stdin.

Examples:
  acorn jq type file.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runJqType,
}

// jqLengthCmd gets length
var jqLengthCmd = &cobra.Command{
	Use:   "length [file]",
	Short: "Get length of array or object",
	Long: `Get the length of a JSON array or number of keys in an object.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn jq length file.json`,
	Aliases: []string{"len"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runJqLength,
}

// jqSelectCmd selects elements
var jqSelectCmd = &cobra.Command{
	Use:   "select [condition] [file]",
	Short: "Select array elements matching condition",
	Long: `Filter array elements that match a condition.

Examples:
  acorn jq select '.active == true' file.json
  acorn jq select '.age > 18' file.json`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runJqSelect,
}

// jqSortCmd sorts array
var jqSortCmd = &cobra.Command{
	Use:   "sort [file]",
	Short: "Sort JSON array",
	Long: `Sort a JSON array.

Examples:
  acorn jq sort file.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runJqSort,
}

// jqUniqueCmd removes duplicates
var jqUniqueCmd = &cobra.Command{
	Use:   "unique [file]",
	Short: "Remove duplicates from array",
	Long: `Remove duplicate elements from a JSON array.

Examples:
  acorn jq unique file.json`,
	Aliases: []string{"uniq"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runJqUnique,
}

// jqFlattenCmd flattens nested JSON
var jqFlattenCmd = &cobra.Command{
	Use:   "flatten [file]",
	Short: "Flatten nested JSON",
	Long: `Flatten nested JSON structure to path-value pairs.

Examples:
  acorn jq flatten file.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runJqFlatten,
}

func init() {
	dataCmd.AddCommand(jqCmd)

	// Add subcommands
	jqCmd.AddCommand(jqStatusCmd)
	jqCmd.AddCommand(jqFilterCmd)
	jqCmd.AddCommand(jqValidateCmd)
	jqCmd.AddCommand(jqFormatCmd)
	jqCmd.AddCommand(jqCompactCmd)
	jqCmd.AddCommand(jqKeysCmd)
	jqCmd.AddCommand(jqTypeCmd)
	jqCmd.AddCommand(jqLengthCmd)
	jqCmd.AddCommand(jqSelectCmd)
	jqCmd.AddCommand(jqSortCmd)
	jqCmd.AddCommand(jqUniqueCmd)
	jqCmd.AddCommand(jqFlattenCmd)

	// Persistent flags
	jqCmd.PersistentFlags().StringVarP(&jqOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	jqCmd.PersistentFlags().BoolVarP(&jqVerbose, "verbose", "v", false,
		"Show verbose output")

	// Filter-specific flags
	jqFilterCmd.Flags().BoolVarP(&jqCompact, "compact", "c", false, "Compact output")
	jqFilterCmd.Flags().BoolVarP(&jqRaw, "raw", "r", false, "Raw string output")
	jqFilterCmd.Flags().BoolVarP(&jqSlurp, "slurp", "s", false, "Slurp input into array")
}

func getJqInput(args []string, argIndex int) ([]byte, error) {
	if len(args) > argIndex {
		return os.ReadFile(args[argIndex])
	}
	return io.ReadAll(os.Stdin)
}

func runJqStatus(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)
	status := helper.GetStatus()

	format, err := output.ParseFormat(jqOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	fmt.Fprintf(os.Stdout, "Installed: %v\n", status.Installed)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:   %s\n", status.Version)
	}
	if status.Location != "" {
		fmt.Fprintf(os.Stdout, "Location:  %s\n", status.Location)
	}

	return nil
}

func runJqFilter(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	expression := args[0]
	input, err := getJqInput(args, 1)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	var result *jq.FilterResult
	if jqCompact {
		result, err = helper.FilterCompact(input, expression)
	} else if jqRaw {
		result, err = helper.FilterRaw(input, expression)
	} else {
		result, err = helper.Filter(input, expression)
	}

	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("filter error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runJqValidate(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if helper.Validate(input) {
		fmt.Fprintf(os.Stdout, "%s Valid JSON\n", output.Success("✓"))
		return nil
	}

	return fmt.Errorf("invalid JSON")
}

func runJqFormat(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	formatted, err := helper.Format(input)
	if err != nil {
		return err
	}

	fmt.Println(string(formatted))
	return nil
}

func runJqCompact(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	compacted, err := helper.Compact(input)
	if err != nil {
		return err
	}

	fmt.Println(string(compacted))
	return nil
}

func runJqKeys(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	keys, err := helper.GetKeys(input)
	if err != nil {
		return err
	}

	for _, k := range keys {
		fmt.Println(k)
	}
	return nil
}

func runJqType(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	t, err := helper.GetType(input)
	if err != nil {
		return err
	}

	fmt.Println(t)
	return nil
}

func runJqLength(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	length, err := helper.GetLength(input)
	if err != nil {
		return err
	}

	fmt.Println(length)
	return nil
}

func runJqSelect(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	condition := args[0]
	input, err := getJqInput(args, 1)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Select(input, condition)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("select error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runJqSort(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Sort(input)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("sort error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runJqUnique(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Unique(input)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("unique error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runJqFlatten(cmd *cobra.Command, args []string) error {
	helper := jq.NewHelper(jqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("jq is not installed")
	}

	input, err := getJqInput(args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Flatten(input)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("flatten error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}
