// Package tools provides tool detection and management capabilities.
package tools

// ToolStatus represents the status of a single tool.
type ToolStatus struct {
	Name      string `json:"name" yaml:"name"`
	Installed bool   `json:"installed" yaml:"installed"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
	Path      string `json:"path,omitempty" yaml:"path,omitempty"`
	Category  string `json:"category" yaml:"category"`
}

// ToolCategory groups tools by purpose.
type ToolCategory struct {
	Name  string       `json:"name" yaml:"name"`
	Tools []ToolStatus `json:"tools" yaml:"tools"`
}

// StatusResult contains the full tools status report.
type StatusResult struct {
	Categories []ToolCategory `json:"categories" yaml:"categories"`
	Summary    StatusSummary  `json:"summary" yaml:"summary"`
}

// StatusSummary provides totals.
type StatusSummary struct {
	Total     int `json:"total" yaml:"total"`
	Installed int `json:"installed" yaml:"installed"`
	Missing   int `json:"missing" yaml:"missing"`
}

// UpdateResult contains the result of an update operation.
type UpdateResult struct {
	Tool    string `json:"tool" yaml:"tool"`
	Success bool   `json:"success" yaml:"success"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}
