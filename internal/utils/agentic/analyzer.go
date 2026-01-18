package agentic

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Analyzer checks Go command files for agentic/JSON-first compliance.
type Analyzer struct {
	cmdDir string
}

// NewAnalyzer creates a new Analyzer for the given command directory.
func NewAnalyzer(cmdDir string) *Analyzer {
	return &Analyzer{cmdDir: cmdDir}
}

// Regular expressions for pattern matching
var (
	// Match output flag declarations
	outputFlagRe = regexp.MustCompile(`StringVarP?\s*\(\s*&\w*[Oo]utput[Ff]ormat\w*\s*,\s*"output"\s*,\s*"o"`)

	// Match output.ParseFormat calls
	parseFormatRe = regexp.MustCompile(`output\.ParseFormat\s*\(`)

	// Match output.NewPrinter calls
	newPrinterRe = regexp.MustCompile(`output\.NewPrinter\s*\(`)

	// Match fmt.Fprintf(os.Stdout, ...) or fmt.Printf(...)
	fmtBypassRe = regexp.MustCompile(`fmt\.(Fprintf\s*\(\s*os\.Stdout|Printf|Fprintln\s*\(\s*os\.Stdout|Println)`)

	// Match cobra.Command declarations
	cobraCommandRe = regexp.MustCompile(`var\s+(\w+)\s*=\s*&cobra\.Command\s*\{`)

	// Match format != output.FormatTable check
	formatCheckRe = regexp.MustCompile(`format\s*!=\s*output\.FormatTable`)

	// Match AddCommand calls to find subcommands
	addCommandRe = regexp.MustCompile(`(\w+)\.AddCommand\s*\(\s*(\w+)\s*\)`)
)

// AuditAll scans all command files and returns audit results.
func (a *Analyzer) AuditAll() (*AuditResult, error) {
	files, err := filepath.Glob(filepath.Join(a.cmdDir, "*.go"))
	if err != nil {
		return nil, err
	}

	result := &AuditResult{
		Commands: make([]CommandAudit, 0),
	}

	for _, file := range files {
		// Skip test files and group files
		baseName := filepath.Base(file)
		if strings.HasSuffix(baseName, "_test.go") {
			continue
		}
		if strings.HasPrefix(baseName, "group_") {
			continue
		}
		if baseName == "root.go" {
			continue
		}

		audit, err := a.auditFile(file)
		if err != nil {
			continue // Skip files with errors
		}

		result.Commands = append(result.Commands, *audit)
	}

	// Calculate summary
	result.Summary = a.calculateSummary(result.Commands)

	return result, nil
}

// auditFile analyzes a single command file.
func (a *Analyzer) auditFile(filePath string) (*CommandAudit, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	text := string(content)
	baseName := filepath.Base(filePath)

	audit := &CommandAudit{
		File:        baseName,
		Command:     strings.TrimSuffix(baseName, ".go"),
		Subcommands: make([]string, 0),
		Issues:      make([]string, 0),
	}

	// Find the main command name
	if matches := cobraCommandRe.FindStringSubmatch(text); len(matches) > 1 {
		audit.Command = matches[1]
	}

	// Count subcommands from AddCommand calls
	subCmds := make(map[string]bool)
	for _, match := range addCommandRe.FindAllStringSubmatch(text, -1) {
		if len(match) > 2 {
			subCmds[match[2]] = true
		}
	}
	for cmd := range subCmds {
		audit.Subcommands = append(audit.Subcommands, cmd)
	}

	// Check for output flag
	audit.HasOutputFlag = outputFlagRe.MatchString(text)

	// Check for ParseFormat usage
	audit.UsesParseFormat = parseFormatRe.MatchString(text)

	// Check for NewPrinter usage
	audit.UsesNewPrinter = newPrinterRe.MatchString(text)

	// Count bypasses (fmt.Printf/Fprintf to stdout)
	audit.BypassCount = len(fmtBypassRe.FindAllString(text, -1))

	// Determine compliance and issues
	a.determineCompliance(audit, text)

	return audit, nil
}

// determineCompliance evaluates compliance and generates issues.
func (a *Analyzer) determineCompliance(audit *CommandAudit, text string) {
	// A command is compliant if it has the output flag and uses ParseFormat
	hasFormatCheck := formatCheckRe.MatchString(text)

	if !audit.HasOutputFlag {
		audit.Issues = append(audit.Issues, "Missing --output/-o flag for format selection")
	}

	if !audit.UsesParseFormat {
		audit.Issues = append(audit.Issues, "Does not use output.ParseFormat() for format parsing")
	}

	if !audit.UsesNewPrinter && audit.HasOutputFlag {
		audit.Issues = append(audit.Issues, "Has output flag but doesn't use output.NewPrinter()")
	}

	if !hasFormatCheck && audit.HasOutputFlag {
		audit.Issues = append(audit.Issues, "No format != output.FormatTable branch for JSON/YAML handling")
	}

	if audit.BypassCount > 10 {
		audit.Issues = append(audit.Issues, "High number of direct stdout writes (consider using printer)")
	}

	// Compliant if has output flag and uses parse format
	audit.Compliant = audit.HasOutputFlag && audit.UsesParseFormat && len(audit.Issues) == 0
}

// calculateSummary computes aggregate statistics.
func (a *Analyzer) calculateSummary(audits []CommandAudit) AuditSummary {
	summary := AuditSummary{
		Total: len(audits),
	}

	for _, audit := range audits {
		if audit.Compliant {
			summary.Compliant++
		}
		summary.TotalIssues += len(audit.Issues)
	}

	if summary.Total > 0 {
		summary.Percentage = float64(summary.Compliant) / float64(summary.Total) * 100
	}

	return summary
}

// ValidateFile performs deep validation on a single file.
func (a *Analyzer) ValidateFile(filePath string) (*ValidationResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &ValidationResult{
		File:            filepath.Base(filePath),
		Command:         strings.TrimSuffix(filepath.Base(filePath), ".go"),
		FormatBranches:  make([]FormatBranch, 0),
		PrinterUsage:    make([]PrinterUsage, 0),
		Bypasses:        make([]BypassInfo, 0),
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	inFormatBranch := false
	braceDepth := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for output flag
		if outputFlagRe.MatchString(line) {
			result.OutputFlag = &FlagInfo{
				Line:      lineNum,
				FlagName:  "output",
				ShortFlag: "o",
			}
		}

		// Check for format check
		if formatCheckRe.MatchString(line) {
			result.FormatBranches = append(result.FormatBranches, FormatBranch{
				Line:      lineNum,
				Condition: strings.TrimSpace(line),
			})
			inFormatBranch = true
			braceDepth = 0
		}

		// Track braces to know when we exit format branch
		if inFormatBranch {
			braceDepth += strings.Count(line, "{") - strings.Count(line, "}")
			if braceDepth <= 0 {
				inFormatBranch = false
			}
		}

		// Check for NewPrinter usage
		if newPrinterRe.MatchString(line) {
			result.PrinterUsage = append(result.PrinterUsage, PrinterUsage{
				Line:     lineNum,
				Function: "output.NewPrinter",
			})
		}

		// Check for bypasses
		if fmtBypassRe.MatchString(line) {
			result.Bypasses = append(result.Bypasses, BypassInfo{
				Line:     lineNum,
				Call:     strings.TrimSpace(line),
				InFormat: inFormatBranch,
			})
		}
	}

	// Determine validity and generate recommendations
	a.determineValidity(result)

	return result, scanner.Err()
}

// determineValidity checks validation result and adds recommendations.
func (a *Analyzer) determineValidity(result *ValidationResult) {
	if result.OutputFlag == nil {
		result.Issues = append(result.Issues, "No --output/-o flag found")
		result.Recommendations = append(result.Recommendations,
			"Add: cmd.PersistentFlags().StringVarP(&outputFormat, \"output\", \"o\", \"table\", \"Output format (table|json|yaml)\")")
	}

	if len(result.FormatBranches) == 0 {
		result.Issues = append(result.Issues, "No format branch found for JSON/YAML handling")
		result.Recommendations = append(result.Recommendations,
			"Add: if format != output.FormatTable { printer := output.NewPrinter(os.Stdout, format); return printer.Print(data) }")
	}

	if len(result.PrinterUsage) == 0 && result.OutputFlag != nil {
		result.Issues = append(result.Issues, "Has output flag but no printer usage found")
		result.Recommendations = append(result.Recommendations,
			"Use output.NewPrinter() for structured output instead of fmt.Printf")
	}

	// Count bypasses outside format branches
	outsideBranch := 0
	for _, bypass := range result.Bypasses {
		if !bypass.InFormat {
			outsideBranch++
		}
	}

	if outsideBranch > 5 {
		result.Recommendations = append(result.Recommendations,
			"Consider reducing direct stdout writes for better agentic support")
	}

	result.Valid = len(result.Issues) == 0
}

// GetPatterns returns example patterns for compliant commands.
func GetPatterns() *Patterns {
	return &Patterns{
		OutputFlag: PatternExample{
			Name:        "Output Flag",
			Description: "Add --output/-o flag to command or parent group",
			Code: `var outputFormat string

func init() {
    cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table",
        "Output format (table|json|yaml)")
}`,
		},
		ParseFormat: PatternExample{
			Name:        "Parse Format",
			Description: "Parse the format string to a typed Format value",
			Code: `func runMyCommand(cmd *cobra.Command, args []string) error {
    format, err := output.ParseFormat(outputFormat)
    if err != nil {
        return err
    }
    // ...
}`,
		},
		FormatSwitch: PatternExample{
			Name:        "Format Switch",
			Description: "Handle non-table formats with the printer",
			Code: `if format != output.FormatTable {
    printer := output.NewPrinter(os.Stdout, format)
    return printer.Print(structuredData)
}
// Table format handling below...`,
		},
		TableOutput: PatternExample{
			Name:        "Table Output",
			Description: "Use output.Table for human-readable table format",
			Code: `table := output.NewTable("NAME", "STATUS", "VERSION")
for _, item := range items {
    table.AddRow(item.Name, item.Status, item.Version)
}
table.Render(os.Stdout)`,
		},
		StructuredData: PatternExample{
			Name:        "Structured Data",
			Description: "Define types with JSON tags for structured output",
			Code: `type MyResult struct {
    Name    string   ` + "`json:\"name\"`" + `
    Status  string   ` + "`json:\"status\"`" + `
    Items   []Item   ` + "`json:\"items\"`" + `
}`,
		},
	}
}
