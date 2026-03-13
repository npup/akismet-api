package akismet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// The result of a usage-limit request.
type UsageLimitResult struct {
	// Monthly API call allowance for your plan; "none" if the key has no limit.
	Limit string
	// Number of API calls made since the start of the current month.
	Usage int
	// Percentage of the monthly limit consumed so far this month.
	Percentage string
	// True if Akismet is currently throttling requests due to consistent overuse.
	Throttled bool
}

// Returns the API usage limit and current-month usage for the configured key.
func (c *Client) GetUsageLimit() (*UsageLimitResult, *AkismetError) {
	url := getUrlWithParameters(ApiEndpoints.UsageLimit, url.Values{"api_key": {c.apiKey}})
	resp, err := c.httpClient.Get(url.String())
	if err != nil {
		return nil, NewAkismetError(err, nil, "")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewAkismetError(err, getAlert(resp), getDebugHelp(resp))
	}

	if strings.TrimSpace(string(body)) == "invalid" {
		return nil, NewAkismetError(fmt.Errorf("akismet: invalid API key"), getAlert(resp), getDebugHelp(resp))
	}

	var raw struct {
		Limit      string `json:"limit"`
		Usage      int    `json:"usage"`
		Percentage string `json:"percentage"`
		Throttled  bool   `json:"throttled"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, NewAkismetError(fmt.Errorf("akismet: unexpected response: %w", err), getAlert(resp), getDebugHelp(resp))
	}

	return &UsageLimitResult{
		Limit:      raw.Limit,
		Usage:      raw.Usage,
		Percentage: raw.Percentage,
		Throttled:  raw.Throttled,
	}, nil
}
