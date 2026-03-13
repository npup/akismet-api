package akismet

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// The result of a usage-limit request.
type UsageLimitResult struct {
	// Monthly API call allowance for your plan. Nil if the key has no limit.
	Limit *int
	// Number of API calls made since the start of the current month.
	Usage int
	// Percentage of the monthly limit consumed so far this month.
	Percentage float64
	// True if Akismet is currently throttling requests due to consistent overuse.
	Throttled bool
}

// Returns the API usage limit and current-month usage for the configured key.
func (c *Client) GetUsageLimit(ctx context.Context) (*UsageLimitResult, *AkismetError) {
	u := getUrlWithParameters(c.endpoints.UsageLimit, url.Values{"api_key": {c.apiKey}})
	resp, akismetErr := c.doGet(ctx, u)
	if akismetErr != nil {
		return nil, akismetErr
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, AkismetErrorFromResponse(err, resp)
	}

	if strings.TrimSpace(string(body)) == BODY_INVALID_MESSAGE {
		return nil, AkismetErrorFromResponse(fmt.Errorf("akismet: invalid API key"), resp)
	}

	var raw struct {
		Limit      string `json:"limit"`
		Usage      int    `json:"usage"`
		Percentage string `json:"percentage"`
		Throttled  bool   `json:"throttled"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, AkismetErrorFromResponse(fmt.Errorf("akismet: unexpected response: %w", err), resp)
	}

	result := &UsageLimitResult{
		Usage:     raw.Usage,
		Throttled: raw.Throttled,
	}

	if raw.Limit != PROP_USAGE_LIMIT_NO_LIMIT {
		if n, err := strconv.Atoi(raw.Limit); err == nil {
			result.Limit = &n
		}
	}

	if p, err := strconv.ParseFloat(raw.Percentage, 64); err == nil {
		result.Percentage = p
	}

	return result, nil
}
