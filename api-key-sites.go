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

// Column by which key-sites results are sorted.
type KeySitesOrder string

const (
	KeySitesOrderTotal          KeySitesOrder = "total"
	KeySitesOrderSpam           KeySitesOrder = "spam"
	KeySitesOrderHam            KeySitesOrder = "ham"
	KeySitesOrderMissedSpam     KeySitesOrder = "missed_spam"
	KeySitesOrderFalsePositives KeySitesOrder = "false_positives"
)

// Optional parameters for GetKeySites.
type KeySitesParams struct {
	// Month to report on, in YYYY-MM format. Defaults to the current month.
	Month string
	// Filter results to sites whose URL contains this string.
	Filter string
	// Column to sort by. Defaults to KeySitesOrderTotal.
	Order KeySitesOrder
	// Maximum number of results. Defaults to 500.
	Limit int
	// Number of results to skip. Defaults to 0.
	Offset int
}

// Usage statistics for a single site.
type KeySiteEntry struct {
	Site           string
	APICalls       int
	Spam           int
	Ham            int
	MissedSpam     int
	FalsePositives int
	IsRevoked      bool
}

// The result of a key-sites request.
type KeySitesResult struct {
	Sites  []KeySiteEntry
	Limit  int
	Offset int
	Total  int
}

// Lists the sites using your API key along with their usage statistics.
// Pass nil for params to use all defaults.
func (c *Client) GetKeySites(ctx context.Context, params *KeySitesParams) (*KeySitesResult, *AkismetError) {
	q := url.Values{"api_key": {c.apiKey}}
	if params != nil {
		if params.Month != "" {
			q.Set("month", params.Month)
		}
		if params.Filter != "" {
			q.Set("filter", params.Filter)
		}
		if params.Order != "" {
			q.Set("order", string(params.Order))
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
	}

	u := getUrlWithParameters(c.endpoints.KeySites, q)
	resp, akismetErr := c.doGet(ctx, u)
	if akismetErr != nil {
		return nil, akismetErr
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, akismetErrorFromResponse(err, resp)
	}

	if strings.TrimSpace(string(body)) == bodyInvalidMessage {
		return nil, akismetErrorFromResponse(fmt.Errorf("akismet: invalid API key"), resp)
	}

	// The response is a flat JSON object where most keys are site URLs mapping
	// to stat objects, with "limit", "offset", and "total" as metadata keys.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, akismetErrorFromResponse(fmt.Errorf("akismet: unexpected response: %w", err), resp)
	}

	result := &KeySitesResult{}
	json.Unmarshal(raw["limit"], &result.Limit)
	json.Unmarshal(raw["offset"], &result.Offset)
	json.Unmarshal(raw["total"], &result.Total)

	for key, val := range raw {
		if key == "limit" || key == "offset" || key == "total" {
			continue
		}
		var entry struct {
			APICalls       int  `json:"api_calls"`
			Spam           int  `json:"spam"`
			Ham            int  `json:"ham"`
			MissedSpam     int  `json:"missed_spam"`
			FalsePositives int  `json:"false_positives"`
			IsRevoked      bool `json:"is_revoked"`
		}
		if err := json.Unmarshal(val, &entry); err != nil {
			continue
		}
		result.Sites = append(result.Sites, KeySiteEntry{
			Site:           key,
			APICalls:       entry.APICalls,
			Spam:           entry.Spam,
			Ham:            entry.Ham,
			MissedSpam:     entry.MissedSpam,
			FalsePositives: entry.FalsePositives,
			IsRevoked:      entry.IsRevoked,
		})
	}

	return result, nil
}
