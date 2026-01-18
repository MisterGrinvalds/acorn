// Package agentic provides utilities for auditing CLI commands for JSON-first compliance.
package agentic

// AuditResult contains the complete audit results for all commands.
type AuditResult struct {
	Summary  AuditSummary   `json:"summary"`
	Commands []CommandAudit `json:"commands"`
}

// AuditSummary contains aggregate statistics about audit compliance.
type AuditSummary struct {
	Total       int     `json:"total"`
	Compliant   int     `json:"compliant"`
	Percentage  float64 `json:"percentage"`
	TotalIssues int     `json:"total_issues"`
}

// CommandAudit represents the audit results for a single command file.
type CommandAudit struct {
	File            string   `json:"file"`
	Command         string   `json:"command"`
	Subcommands     []string `json:"subcommands"`
	Compliant       bool     `json:"compliant"`
	HasOutputFlag   bool     `json:"has_output_flag"`
	UsesParseFormat bool     `json:"uses_parse_format"`
	UsesNewPrinter  bool     `json:"uses_new_printer"`
	BypassCount     int      `json:"bypass_count"`
	Issues          []string `json:"issues"`
}

// ValidationResult contains deep validation results for a single file.
type ValidationResult struct {
	File           string          `json:"file"`
	Command        string          `json:"command"`
	Valid          bool            `json:"valid"`
	OutputFlag     *FlagInfo       `json:"output_flag,omitempty"`
	FormatBranches []FormatBranch  `json:"format_branches"`
	PrinterUsage   []PrinterUsage  `json:"printer_usage"`
	Bypasses       []BypassInfo    `json:"bypasses"`
	Issues         []string        `json:"issues"`
	Recommendations []string       `json:"recommendations"`
}

// FlagInfo describes an output format flag declaration.
type FlagInfo struct {
	VarName   string `json:"var_name"`
	FlagName  string `json:"flag_name"`
	ShortFlag string `json:"short_flag"`
	Default   string `json:"default"`
	Line      int    `json:"line"`
}

// FormatBranch represents a format check branch in the code.
type FormatBranch struct {
	Line      int    `json:"line"`
	Condition string `json:"condition"`
	HasTable  bool   `json:"has_table"`
	HasJSON   bool   `json:"has_json"`
}

// PrinterUsage represents a call to output.NewPrinter.
type PrinterUsage struct {
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// BypassInfo represents a direct stdout write that bypasses the printer.
type BypassInfo struct {
	Line     int    `json:"line"`
	Call     string `json:"call"`
	InFormat bool   `json:"in_format_branch"`
}

// PatternExample contains code pattern examples for compliance.
type PatternExample struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

// Patterns contains all pattern examples.
type Patterns struct {
	OutputFlag    PatternExample `json:"output_flag"`
	ParseFormat   PatternExample `json:"parse_format"`
	FormatSwitch  PatternExample `json:"format_switch"`
	TableOutput   PatternExample `json:"table_output"`
	StructuredData PatternExample `json:"structured_data"`
}
