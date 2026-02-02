package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/data/yq"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	yqVerbose     bool
	yqInPlace     bool
	yqInputFormat string
	yqOutFormat   string
)

// yqCmd represents the yq command group
var yqCmd = &cobra.Command{
	Use:   "yq",
	Short: "YAML processing commands",
	Long: `YAML processing commands using yq.

Provides filtering, transformation, conversion, and merging of YAML data.

Examples:
  acorn yq status                     # Show yq status
  acorn yq filter '.name' file.yaml   # Extract field
  acorn yq validate file.yaml         # Validate YAML
  acorn yq convert file.yaml -o json  # Convert to JSON
  acorn yq merge base.yaml overlay.yaml # Merge files`,
}

// yqStatusCmd shows yq status
var yqStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show yq installation status",
	Long: `Display yq installation status and version.

Examples:
  acorn yq status
  acorn yq status -o json`,
	RunE: runYqStatus,
}

// yqFilterCmd applies a filter
var yqFilterCmd = &cobra.Command{
	Use:   "filter [expression] [file]",
	Short: "Apply yq expression to YAML",
	Long: `Apply a yq expression to YAML input.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn yq filter '.name' file.yaml
  acorn yq filter '.items[]' file.yaml
  cat file.yaml | acorn yq filter '.name'
  acorn yq filter -i '.version = "2.0"' file.yaml  # in-place edit`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runYqFilter,
}

// yqValidateCmd validates YAML
var yqValidateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate YAML",
	Long: `Check if input is valid YAML.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn yq validate file.yaml
  cat file.yaml | acorn yq validate`,
	Args: cobra.MaximumNArgs(1),
	RunE: runYqValidate,
}

// yqFormatCmd formats YAML
var yqFormatCmd = &cobra.Command{
	Use:   "format [file]",
	Short: "Pretty print YAML",
	Long: `Format YAML with proper indentation.

Reads from file if provided, otherwise from stdin.

Examples:
  acorn yq format file.yaml
  cat file.yaml | acorn yq format`,
	Aliases: []string{"fmt", "pretty"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runYqFormat,
}

// yqConvertCmd converts between formats
var yqConvertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert between formats",
	Long: `Convert between YAML, JSON, XML, TOML, CSV, and other formats.

Supported formats: yaml, json, xml, toml, csv, tsv, props

Examples:
  acorn yq convert file.yaml -o json         # YAML to JSON
  acorn yq convert file.json -p json -o yaml # JSON to YAML
  acorn yq convert file.xml -p xml -o yaml   # XML to YAML
  acorn yq convert file.yaml -o toml         # YAML to TOML`,
	Args: cobra.MaximumNArgs(1),
	RunE: runYqConvert,
}

// yqMergeCmd merges files
var yqMergeCmd = &cobra.Command{
	Use:   "merge [base] [overlay...]",
	Short: "Merge YAML files",
	Long: `Merge multiple YAML files together.

Later files override earlier ones.

Examples:
  acorn yq merge base.yaml overlay.yaml
  acorn yq merge base.yaml env/dev.yaml secrets.yaml`,
	Args: cobra.MinimumNArgs(2),
	RunE: runYqMerge,
}

// yqGetCmd gets a value
var yqGetCmd = &cobra.Command{
	Use:   "get [path] [file]",
	Short: "Get value at path",
	Long: `Extract a value at a specific path.

Examples:
  acorn yq get '.metadata.name' file.yaml
  acorn yq get '.spec.replicas' deployment.yaml`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runYqGet,
}

// yqSetCmd sets a value
var yqSetCmd = &cobra.Command{
	Use:   "set [path] [value] [file]",
	Short: "Set value at path",
	Long: `Set a value at a specific path.

Examples:
  acorn yq set '.metadata.name' 'myapp' file.yaml
  acorn yq set '.spec.replicas' '3' file.yaml
  acorn yq set -i '.version' '"2.0"' file.yaml  # in-place`,
	Args: cobra.RangeArgs(2, 3),
	RunE: runYqSet,
}

// yqDeleteCmd deletes a key
var yqDeleteCmd = &cobra.Command{
	Use:   "delete [path] [file]",
	Short: "Delete key at path",
	Long: `Delete a key at a specific path.

Examples:
  acorn yq delete '.metadata.annotations' file.yaml
  acorn yq delete -i '.spec.nodeSelector' file.yaml`,
	Aliases: []string{"del", "rm"},
	Args:    cobra.RangeArgs(1, 2),
	RunE:    runYqDelete,
}

// yqKeysCmd gets keys
var yqKeysCmd = &cobra.Command{
	Use:   "keys [file]",
	Short: "Get keys from YAML object",
	Long: `Extract all keys from a YAML object.

Examples:
  acorn yq keys file.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runYqKeys,
}

// yqSortCmd sorts keys
var yqSortCmd = &cobra.Command{
	Use:   "sort [file]",
	Short: "Sort YAML keys alphabetically",
	Long: `Sort all keys in YAML alphabetically.

Examples:
  acorn yq sort file.yaml
  acorn yq sort -i file.yaml  # in-place`,
	Args: cobra.MaximumNArgs(1),
	RunE: runYqSort,
}

// yqDocsCmd handles multi-document YAML
var yqDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Multi-document YAML operations",
	Long:  `Commands for handling multi-document YAML files.`,
}

// yqDocsCountCmd counts documents
var yqDocsCountCmd = &cobra.Command{
	Use:   "count [file]",
	Short: "Count documents in YAML",
	Long: `Count the number of documents in a multi-document YAML file.

Examples:
  acorn yq docs count manifests.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runYqDocsCount,
}

// yqDocsGetCmd gets a specific document
var yqDocsGetCmd = &cobra.Command{
	Use:   "get [index] [file]",
	Short: "Get document by index",
	Long: `Extract a specific document from a multi-document YAML file.

Index is 0-based.

Examples:
  acorn yq docs get 0 manifests.yaml  # First document
  acorn yq docs get 2 manifests.yaml  # Third document`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runYqDocsGet,
}

// yqDocsSplitCmd splits documents
var yqDocsSplitCmd = &cobra.Command{
	Use:   "split [file]",
	Short: "Split multi-document YAML",
	Long: `Split a multi-document YAML file and output each document.

Examples:
  acorn yq docs split manifests.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runYqDocsSplit,
}

func init() {
	dataCmd.AddCommand(yqCmd)

	// Add subcommands
	yqCmd.AddCommand(yqStatusCmd)
	yqCmd.AddCommand(yqFilterCmd)
	yqCmd.AddCommand(yqValidateCmd)
	yqCmd.AddCommand(yqFormatCmd)
	yqCmd.AddCommand(yqConvertCmd)
	yqCmd.AddCommand(yqMergeCmd)
	yqCmd.AddCommand(yqGetCmd)
	yqCmd.AddCommand(yqSetCmd)
	yqCmd.AddCommand(yqDeleteCmd)
	yqCmd.AddCommand(yqKeysCmd)
	yqCmd.AddCommand(yqSortCmd)
	yqCmd.AddCommand(yqDocsCmd)
	yqCmd.AddCommand(configcmd.NewConfigRouter("yq"))

	// Docs subcommands
	yqDocsCmd.AddCommand(yqDocsCountCmd)
	yqDocsCmd.AddCommand(yqDocsGetCmd)
	yqDocsCmd.AddCommand(yqDocsSplitCmd)

	// Persistent flags (output format is inherited from root command)
	yqCmd.PersistentFlags().BoolVarP(&yqVerbose, "verbose", "v", false,
		"Show verbose output")

	// Command-specific flags
	yqFilterCmd.Flags().BoolVarP(&yqInPlace, "in-place", "i", false, "Edit file in place")
	yqSetCmd.Flags().BoolVarP(&yqInPlace, "in-place", "i", false, "Edit file in place")
	yqDeleteCmd.Flags().BoolVarP(&yqInPlace, "in-place", "i", false, "Edit file in place")
	yqSortCmd.Flags().BoolVarP(&yqInPlace, "in-place", "i", false, "Edit file in place")

	// Convert flags
	yqConvertCmd.Flags().StringVarP(&yqInputFormat, "parse", "p", "yaml",
		"Input format (yaml|json|xml|toml|csv|tsv|props)")
	yqConvertCmd.Flags().StringVar(&yqOutFormat, "to", "json",
		"Output format (yaml|json|xml|toml|csv|tsv|props)")
}

func getYqInput(cmd *cobra.Command, args []string, argIndex int) ([]byte, error) {
	ioHelper := ioutils.IO(cmd)

	// Use positional file argument if provided
	if len(args) > argIndex {
		return os.ReadFile(args[argIndex])
	}

	// Check for --input-file flag or piped stdin
	if ioHelper.HasInput() {
		return ioHelper.ReadRaw()
	}

	return nil, fmt.Errorf("no input provided: specify a file or pipe input")
}

func runYqStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := yq.NewHelper(yqVerbose)
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
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

func runYqFilter(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	expression := args[0]

	// In-place edit
	if yqInPlace && len(args) > 1 {
		return helper.EditInPlace(args[1], expression)
	}

	input, err := getYqInput(cmd, args, 1)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Filter(input, expression)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("filter error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqValidate(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	if len(args) > 0 {
		valid, err := helper.ValidateFile(args[0])
		if err != nil {
			return err
		}
		if valid {
			fmt.Fprintf(os.Stdout, "%s Valid YAML\n", output.Success("✓"))
			return nil
		}
		return fmt.Errorf("invalid YAML")
	}

	ioHelper := ioutils.IO(cmd)
	if !ioHelper.HasInput() {
		return fmt.Errorf("no input provided: specify a file or pipe input")
	}

	input, err := ioHelper.ReadRaw()
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if helper.Validate(input) {
		fmt.Fprintf(os.Stdout, "%s Valid YAML\n", output.Success("✓"))
		return nil
	}

	return fmt.Errorf("invalid YAML")
}

func runYqFormat(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	input, err := getYqInput(cmd, args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Format(input)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("format error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqConvert(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	input, err := getYqInput(cmd, args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Convert(input, yq.ConvertFormat(yqInputFormat), yq.ConvertFormat(yqOutFormat))
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("convert error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqMerge(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	result, err := helper.MergeFiles(args...)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("merge error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqGet(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	path := args[0]
	input, err := getYqInput(cmd, args, 1)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Get(input, path)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("get error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqSet(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	path := args[0]
	value := args[1]

	// In-place edit
	if yqInPlace && len(args) > 2 {
		expression := fmt.Sprintf("%s = %s", path, value)
		return helper.EditInPlace(args[2], expression)
	}

	input, err := getYqInput(cmd, args, 2)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Set(input, path, value)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("set error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqDelete(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	path := args[0]

	// In-place edit
	if yqInPlace && len(args) > 1 {
		expression := fmt.Sprintf("del(%s)", path)
		return helper.EditInPlace(args[1], expression)
	}

	input, err := getYqInput(cmd, args, 1)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.Delete(input, path)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("delete error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqKeys(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	input, err := getYqInput(cmd, args, 0)
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

func runYqSort(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	// In-place edit
	if yqInPlace && len(args) > 0 {
		return helper.EditInPlace(args[0], "sort_keys(.)")
	}

	input, err := getYqInput(cmd, args, 0)
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

func runYqDocsCount(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	input, err := getYqInput(cmd, args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	count, err := helper.CountDocuments(input)
	if err != nil {
		return err
	}

	fmt.Println(count)
	return nil
}

func runYqDocsGet(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	var index int
	if _, err := fmt.Sscanf(args[0], "%d", &index); err != nil {
		return fmt.Errorf("invalid index: %s", args[0])
	}

	input, err := getYqInput(cmd, args, 1)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	result, err := helper.GetDocument(input, index)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("get document error: %s", result.Error)
	}

	fmt.Print(result.Output)
	return nil
}

func runYqDocsSplit(cmd *cobra.Command, args []string) error {
	helper := yq.NewHelper(yqVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("yq is not installed")
	}

	input, err := getYqInput(cmd, args, 0)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	docs, err := helper.SplitDocuments(input)
	if err != nil {
		return err
	}

	for i, doc := range docs {
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Println(doc)
	}
	return nil
}
