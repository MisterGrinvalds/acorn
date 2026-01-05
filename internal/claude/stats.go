package claude

import (
	"fmt"
	"sort"
)

// Stats represents the stats-cache.json structure.
type Stats struct {
	Version          int                   `json:"version" yaml:"version"`
	LastComputedDate string                `json:"lastComputedDate" yaml:"lastComputedDate"`
	DailyActivity    []DailyActivity       `json:"dailyActivity" yaml:"dailyActivity"`
	DailyModelTokens []DailyModelTokens    `json:"dailyModelTokens" yaml:"dailyModelTokens"`
	ModelUsage       map[string]ModelUsage `json:"modelUsage" yaml:"modelUsage"`
	TotalSessions    int                   `json:"totalSessions" yaml:"totalSessions"`
	TotalMessages    int                   `json:"totalMessages" yaml:"totalMessages"`
	LongestSession   *LongestSession       `json:"longestSession,omitempty" yaml:"longestSession,omitempty"`
	FirstSessionDate string                `json:"firstSessionDate,omitempty" yaml:"firstSessionDate,omitempty"`
	HourCounts       map[string]int        `json:"hourCounts,omitempty" yaml:"hourCounts,omitempty"`
}

// DailyActivity represents activity for a single day.
type DailyActivity struct {
	Date          string `json:"date" yaml:"date"`
	MessageCount  int    `json:"messageCount" yaml:"messageCount"`
	SessionCount  int    `json:"sessionCount" yaml:"sessionCount"`
	ToolCallCount int    `json:"toolCallCount" yaml:"toolCallCount"`
}

// DailyModelTokens represents token usage by model for a single day.
type DailyModelTokens struct {
	Date          string         `json:"date" yaml:"date"`
	TokensByModel map[string]int `json:"tokensByModel" yaml:"tokensByModel"`
}

// ModelUsage represents token usage for a single model.
type ModelUsage struct {
	InputTokens              int     `json:"inputTokens" yaml:"inputTokens"`
	OutputTokens             int     `json:"outputTokens" yaml:"outputTokens"`
	CacheReadInputTokens     int     `json:"cacheReadInputTokens" yaml:"cacheReadInputTokens"`
	CacheCreationInputTokens int     `json:"cacheCreationInputTokens" yaml:"cacheCreationInputTokens"`
	WebSearchRequests        int     `json:"webSearchRequests,omitempty" yaml:"webSearchRequests,omitempty"`
	CostUSD                  float64 `json:"costUSD,omitempty" yaml:"costUSD,omitempty"`
}

// LongestSession represents the longest session info.
type LongestSession struct {
	SessionID    string `json:"sessionId" yaml:"sessionId"`
	Duration     int64  `json:"duration" yaml:"duration"`
	MessageCount int    `json:"messageCount" yaml:"messageCount"`
	Timestamp    string `json:"timestamp" yaml:"timestamp"`
}

// StatsSummary is a simplified view of stats for display.
type StatsSummary struct {
	TotalSessions  int             `json:"total_sessions" yaml:"total_sessions"`
	TotalMessages  int             `json:"total_messages" yaml:"total_messages"`
	LastComputed   string          `json:"last_computed" yaml:"last_computed"`
	ModelBreakdown []ModelSummary  `json:"model_breakdown" yaml:"model_breakdown"`
	RecentActivity []DailyActivity `json:"recent_activity" yaml:"recent_activity"`
}

// ModelSummary is a simplified view of model usage.
type ModelSummary struct {
	Model        string `json:"model" yaml:"model"`
	InputTokens  int    `json:"input_tokens" yaml:"input_tokens"`
	OutputTokens int    `json:"output_tokens" yaml:"output_tokens"`
	CacheTokens  int    `json:"cache_tokens" yaml:"cache_tokens"`
	TotalTokens  int    `json:"total_tokens" yaml:"total_tokens"`
}

// TokenUsage is a view of token usage by model.
type TokenUsage struct {
	Models []ModelTokenDetail `json:"models" yaml:"models"`
	Total  TokenTotals        `json:"total" yaml:"total"`
}

// ModelTokenDetail shows detailed token info for a model.
type ModelTokenDetail struct {
	Model                    string `json:"model" yaml:"model"`
	InputTokens              int    `json:"input_tokens" yaml:"input_tokens"`
	OutputTokens             int    `json:"output_tokens" yaml:"output_tokens"`
	CacheReadInputTokens     int    `json:"cache_read_tokens" yaml:"cache_read_tokens"`
	CacheCreationInputTokens int    `json:"cache_creation_tokens" yaml:"cache_creation_tokens"`
}

// TokenTotals shows total token counts.
type TokenTotals struct {
	Input         int `json:"input" yaml:"input"`
	Output        int `json:"output" yaml:"output"`
	CacheRead     int `json:"cache_read" yaml:"cache_read"`
	CacheCreation int `json:"cache_creation" yaml:"cache_creation"`
}

// DailyUsage is a view of daily token usage.
type DailyUsage struct {
	Days []DailyTokenSummary `json:"days" yaml:"days"`
}

// DailyTokenSummary shows token usage for a single day.
type DailyTokenSummary struct {
	Date   string             `json:"date" yaml:"date"`
	Models []ModelDailyTokens `json:"models" yaml:"models"`
	Total  int                `json:"total" yaml:"total"`
}

// ModelDailyTokens shows tokens for a model on a specific day.
type ModelDailyTokens struct {
	Model  string `json:"model" yaml:"model"`
	Tokens int    `json:"tokens" yaml:"tokens"`
}

// GetStats reads and returns the full stats.
func (h *Helper) GetStats() (*Stats, error) {
	if !h.FileExists(h.paths.StatsCache) {
		return nil, fmt.Errorf("no stats file found at %s", h.paths.StatsCache)
	}

	var stats Stats
	if err := h.ReadJSONFile(h.paths.StatsCache, &stats); err != nil {
		return nil, fmt.Errorf("failed to read stats: %w", err)
	}

	return &stats, nil
}

// GetStatsSummary returns a summary of stats for display.
func (h *Helper) GetStatsSummary() (*StatsSummary, error) {
	stats, err := h.GetStats()
	if err != nil {
		return nil, err
	}

	summary := &StatsSummary{
		TotalSessions: stats.TotalSessions,
		TotalMessages: stats.TotalMessages,
		LastComputed:  stats.LastComputedDate,
	}

	// Get model breakdown
	for model, usage := range stats.ModelUsage {
		summary.ModelBreakdown = append(summary.ModelBreakdown, ModelSummary{
			Model:        model,
			InputTokens:  usage.InputTokens,
			OutputTokens: usage.OutputTokens,
			CacheTokens:  usage.CacheReadInputTokens,
			TotalTokens:  usage.InputTokens + usage.OutputTokens,
		})
	}

	// Sort by total tokens descending
	sort.Slice(summary.ModelBreakdown, func(i, j int) bool {
		return summary.ModelBreakdown[i].TotalTokens > summary.ModelBreakdown[j].TotalTokens
	})

	// Get last 7 days of activity
	activityLen := len(stats.DailyActivity)
	start := 0
	if activityLen > 7 {
		start = activityLen - 7
	}
	summary.RecentActivity = stats.DailyActivity[start:]

	return summary, nil
}

// GetTokenUsage returns token usage breakdown by model.
func (h *Helper) GetTokenUsage() (*TokenUsage, error) {
	stats, err := h.GetStats()
	if err != nil {
		return nil, err
	}

	usage := &TokenUsage{}

	for model, mu := range stats.ModelUsage {
		usage.Models = append(usage.Models, ModelTokenDetail{
			Model:                    model,
			InputTokens:              mu.InputTokens,
			OutputTokens:             mu.OutputTokens,
			CacheReadInputTokens:     mu.CacheReadInputTokens,
			CacheCreationInputTokens: mu.CacheCreationInputTokens,
		})
		usage.Total.Input += mu.InputTokens
		usage.Total.Output += mu.OutputTokens
		usage.Total.CacheRead += mu.CacheReadInputTokens
		usage.Total.CacheCreation += mu.CacheCreationInputTokens
	}

	// Sort by output tokens descending
	sort.Slice(usage.Models, func(i, j int) bool {
		return usage.Models[i].OutputTokens > usage.Models[j].OutputTokens
	})

	return usage, nil
}

// GetDailyUsage returns daily token usage for the last N days.
func (h *Helper) GetDailyUsage(days int) (*DailyUsage, error) {
	stats, err := h.GetStats()
	if err != nil {
		return nil, err
	}

	usage := &DailyUsage{}

	// Get last N days
	tokensLen := len(stats.DailyModelTokens)
	start := 0
	if tokensLen > days {
		start = tokensLen - days
	}

	for _, day := range stats.DailyModelTokens[start:] {
		summary := DailyTokenSummary{
			Date: day.Date,
		}
		for model, tokens := range day.TokensByModel {
			summary.Models = append(summary.Models, ModelDailyTokens{
				Model:  model,
				Tokens: tokens,
			})
			summary.Total += tokens
		}
		// Sort models by tokens descending
		sort.Slice(summary.Models, func(i, j int) bool {
			return summary.Models[i].Tokens > summary.Models[j].Tokens
		})
		usage.Days = append(usage.Days, summary)
	}

	return usage, nil
}
