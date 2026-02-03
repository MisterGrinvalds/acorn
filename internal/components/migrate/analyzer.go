// Package migrate provides tools for analyzing and migrating shell components to Go.
package migrate

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/utils/component"
)

// FunctionType categorizes shell functions for migration decisions.
type FunctionType string

const (
	// FuncTypeAction performs an action (good migration candidate)
	FuncTypeAction FunctionType = "action"
	// FuncTypeWrapper wraps another command (keep as shell)
	FuncTypeWrapper FunctionType = "wrapper"
	// FuncTypeAlias simple alias-like function (keep as shell)
	FuncTypeAlias FunctionType = "alias"
	// FuncTypeEnv sets environment (keep as shell)
	FuncTypeEnv FunctionType = "env"
	// FuncTypeUnknown cannot be determined
	FuncTypeUnknown FunctionType = "unknown"
)

// ShellFunction represents a parsed shell function.
type ShellFunction struct {
	Name         string       `json:"name" yaml:"name"`
	File         string       `json:"file" yaml:"file"`
	LineNumber   int          `json:"line_number" yaml:"line_number"`
	Body         string       `json:"body" yaml:"body"`
	LineCount    int          `json:"line_count" yaml:"line_count"`
	Type         FunctionType `json:"type" yaml:"type"`
	Suggestion   string       `json:"suggestion" yaml:"suggestion"`
	Complexity   string       `json:"complexity" yaml:"complexity"` // low, medium, high
	UsesExternal []string     `json:"uses_external,omitempty" yaml:"uses_external,omitempty"`
}

// ComponentAnalysis contains migration analysis for a component.
type ComponentAnalysis struct {
	Component *component.Component `json:"component" yaml:"component"`
	Functions []ShellFunction      `json:"functions" yaml:"functions"`
	Aliases   []ShellAlias         `json:"aliases" yaml:"aliases"`
	EnvVars   []EnvVar             `json:"env_vars" yaml:"env_vars"`
	Summary   AnalysisSummary      `json:"summary" yaml:"summary"`
}

// ShellAlias represents a shell alias.
type ShellAlias struct {
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
	File    string `json:"file" yaml:"file"`
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
	File  string `json:"file" yaml:"file"`
}

// AnalysisSummary summarizes the migration analysis.
type AnalysisSummary struct {
	TotalFunctions    int    `json:"total_functions" yaml:"total_functions"`
	ActionFunctions   int    `json:"action_functions" yaml:"action_functions"`
	WrapperFunctions  int    `json:"wrapper_functions" yaml:"wrapper_functions"`
	TotalAliases      int    `json:"total_aliases" yaml:"total_aliases"`
	TotalEnvVars      int    `json:"total_env_vars" yaml:"total_env_vars"`
	MigrationScore    int    `json:"migration_score" yaml:"migration_score"` // 0-100
	RecommendedAction string `json:"recommended_action" yaml:"recommended_action"`
}

// Analyzer performs migration analysis on components.
type Analyzer struct {
	dotfilesRoot string
}

// NewAnalyzer creates a new Analyzer.
func NewAnalyzer(dotfilesRoot string) *Analyzer {
	return &Analyzer{dotfilesRoot: dotfilesRoot}
}

// AnalyzeComponent performs full analysis on a component.
func (a *Analyzer) AnalyzeComponent(comp *component.Component) (*ComponentAnalysis, error) {
	analysis := &ComponentAnalysis{
		Component: comp,
		Functions: []ShellFunction{},
		Aliases:   []ShellAlias{},
		EnvVars:   []EnvVar{},
	}

	// Parse functions.sh
	functionsPath := filepath.Join(comp.Path, "functions.sh")
	if funcs, err := a.parseFunctions(functionsPath); err == nil {
		analysis.Functions = funcs
	}

	// Parse aliases.sh
	aliasesPath := filepath.Join(comp.Path, "aliases.sh")
	if aliases, err := a.parseAliases(aliasesPath); err == nil {
		analysis.Aliases = aliases
	}

	// Parse env.sh
	envPath := filepath.Join(comp.Path, "env.sh")
	if envVars, err := a.parseEnvVars(envPath); err == nil {
		analysis.EnvVars = envVars
	}

	// Calculate summary
	analysis.Summary = a.calculateSummary(analysis)

	return analysis, nil
}

// parseFunctions extracts functions from a shell file.
func (a *Analyzer) parseFunctions(path string) ([]ShellFunction, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var functions []ShellFunction
	var currentFunc *ShellFunction
	var bodyLines []string
	braceCount := 0

	// Regex patterns
	funcStartPattern := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*\(\)\s*\{?\s*$`)
	funcStartPattern2 := regexp.MustCompile(`^function\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\{?\s*$`)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Skip comments and empty lines when not in function
		if currentFunc == nil {
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
				continue
			}
		}

		// Check for function start
		if currentFunc == nil {
			var funcName string
			if matches := funcStartPattern.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				funcName = matches[1]
			} else if matches := funcStartPattern2.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				funcName = matches[1]
			}

			if funcName != "" {
				currentFunc = &ShellFunction{
					Name:       funcName,
					File:       filepath.Base(path),
					LineNumber: lineNum,
				}
				bodyLines = []string{}
				braceCount = strings.Count(line, "{") - strings.Count(line, "}")
				continue
			}
		}

		// Inside function
		if currentFunc != nil {
			bodyLines = append(bodyLines, line)
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")

			// Function ends when braces balance
			if braceCount <= 0 {
				currentFunc.Body = strings.Join(bodyLines, "\n")
				currentFunc.LineCount = len(bodyLines)
				currentFunc.Type = a.classifyFunction(currentFunc)
				currentFunc.Complexity = a.assessComplexity(currentFunc)
				currentFunc.Suggestion = a.generateSuggestion(currentFunc)
				currentFunc.UsesExternal = a.findExternalCommands(currentFunc.Body)
				functions = append(functions, *currentFunc)
				currentFunc = nil
			}
		}
	}

	return functions, scanner.Err()
}

// parseAliases extracts aliases from a shell file.
func (a *Analyzer) parseAliases(path string) ([]ShellAlias, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var aliases []ShellAlias
	aliasPattern := regexp.MustCompile(`^alias\s+([a-zA-Z_][a-zA-Z0-9_-]*)=['"]?(.+?)['"]?\s*$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := aliasPattern.FindStringSubmatch(line); len(matches) > 2 {
			aliases = append(aliases, ShellAlias{
				Name:    matches[1],
				Command: strings.Trim(matches[2], "'\""),
				File:    filepath.Base(path),
			})
		}
	}

	return aliases, scanner.Err()
}

// parseEnvVars extracts environment variables from a shell file.
func (a *Analyzer) parseEnvVars(path string) ([]EnvVar, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var envVars []EnvVar
	exportPattern := regexp.MustCompile(`^export\s+([A-Z_][A-Z0-9_]*)=(.*)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := exportPattern.FindStringSubmatch(line); len(matches) > 2 {
			envVars = append(envVars, EnvVar{
				Name:  matches[1],
				Value: strings.Trim(matches[2], "'\""),
				File:  filepath.Base(path),
			})
		}
	}

	return envVars, scanner.Err()
}

// classifyFunction determines the function type.
func (a *Analyzer) classifyFunction(f *ShellFunction) FunctionType {
	body := strings.ToLower(f.Body)

	// Check for environment setting
	if strings.Contains(body, "export ") {
		return FuncTypeEnv
	}

	// Check if it's a simple wrapper (single command passthrough)
	lines := strings.Split(strings.TrimSpace(f.Body), "\n")
	nonCommentLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") && trimmed != "}" {
			nonCommentLines++
		}
	}

	if nonCommentLines <= 2 && strings.Contains(body, `"$@"`) {
		return FuncTypeWrapper
	}

	// Check for action indicators
	actionIndicators := []string{
		"mkdir", "rm ", "cp ", "mv ", "curl", "wget",
		"git clone", "git init", "git add", "git commit",
		"echo \"", "printf", "cat >", "cat <<",
		"if [", "for ", "while ", "case ",
	}

	for _, indicator := range actionIndicators {
		if strings.Contains(body, indicator) {
			return FuncTypeAction
		}
	}

	// Check for simple alias-like functions
	if nonCommentLines == 1 {
		return FuncTypeAlias
	}

	// Default to action if it has significant logic
	if f.LineCount > 5 {
		return FuncTypeAction
	}

	return FuncTypeWrapper
}

// assessComplexity assesses function complexity.
func (a *Analyzer) assessComplexity(f *ShellFunction) string {
	score := 0

	// Line count
	if f.LineCount > 50 {
		score += 3
	} else if f.LineCount > 20 {
		score += 2
	} else if f.LineCount > 10 {
		score += 1
	}

	// Control structures
	body := f.Body
	score += strings.Count(body, "if [")
	score += strings.Count(body, "for ")
	score += strings.Count(body, "while ")
	score += strings.Count(body, "case ")

	// External command calls
	score += len(a.findExternalCommands(body)) / 3

	if score >= 5 {
		return "high"
	} else if score >= 2 {
		return "medium"
	}
	return "low"
}

// generateSuggestion creates a migration suggestion.
func (a *Analyzer) generateSuggestion(f *ShellFunction) string {
	switch f.Type {
	case FuncTypeAction:
		if f.Complexity == "high" {
			return fmt.Sprintf("Migrate to Go: acorn %s %s (complex, benefit from Go error handling)", f.Name, strings.ReplaceAll(f.Name, "_", " "))
		}
		return fmt.Sprintf("Migrate to Go: acorn %s", strings.ReplaceAll(f.Name, "_", " "))
	case FuncTypeWrapper:
		return "Keep as shell wrapper - simple passthrough"
	case FuncTypeAlias:
		return "Keep as shell alias - too simple to migrate"
	case FuncTypeEnv:
		return "Keep as shell - sets environment state"
	default:
		return "Review manually"
	}
}

// findExternalCommands finds external command calls in function body.
func (a *Analyzer) findExternalCommands(body string) []string {
	commonCommands := []string{
		"git", "curl", "wget", "docker", "kubectl", "helm",
		"npm", "yarn", "pnpm", "go", "python", "pip", "uv",
		"brew", "apt", "yum", "make", "cargo",
	}

	var found []string
	seen := make(map[string]bool)

	for _, cmd := range commonCommands {
		if strings.Contains(body, cmd+" ") || strings.Contains(body, cmd+"\n") {
			if !seen[cmd] {
				found = append(found, cmd)
				seen[cmd] = true
			}
		}
	}

	return found
}

// calculateSummary generates the analysis summary.
func (a *Analyzer) calculateSummary(analysis *ComponentAnalysis) AnalysisSummary {
	summary := AnalysisSummary{
		TotalFunctions: len(analysis.Functions),
		TotalAliases:   len(analysis.Aliases),
		TotalEnvVars:   len(analysis.EnvVars),
	}

	for _, f := range analysis.Functions {
		switch f.Type {
		case FuncTypeAction:
			summary.ActionFunctions++
		case FuncTypeWrapper:
			summary.WrapperFunctions++
		}
	}

	// Calculate migration score (0-100)
	// Higher score = more worth migrating
	if summary.TotalFunctions > 0 {
		actionRatio := float64(summary.ActionFunctions) / float64(summary.TotalFunctions)
		summary.MigrationScore = int(actionRatio * 100)
	}

	// Generate recommendation
	switch {
	case summary.MigrationScore >= 70:
		summary.RecommendedAction = "High priority - many action functions benefit from Go migration"
	case summary.MigrationScore >= 40:
		summary.RecommendedAction = "Medium priority - selective migration of action functions"
	case summary.MigrationScore >= 20:
		summary.RecommendedAction = "Low priority - mostly wrappers, keep in shell"
	default:
		summary.RecommendedAction = "Skip - component is primarily shell configuration"
	}

	return summary
}
