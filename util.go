package akismet

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func getDebugHelp(resp *http.Response) string {
	return resp.Header.Get(akismetHeaders.DebugHelp)
}

func getAlert(resp *http.Response) *Alert {
	code := -1
	if codeStr := resp.Header.Get(akismetHeaders.AlertCode); codeStr != "" {
		fmt.Sscanf(codeStr, "%d", &code)
	}
	message := resp.Header.Get(akismetHeaders.AlertMsg)
	if code == -1 || message == "" {
		return nil
	}
	return newAlert(code, message)
}

func getUrlWithParameters(endpoint string, params url.Values) url.URL {
	url, _ := url.Parse(endpoint)
	url.RawQuery = params.Encode()
	return *url
}

func (c *Client) doPost(ctx context.Context, endpoint string, values url.Values) (*http.Response, *AkismetError) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, newAkismetError(err, nil, "")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newAkismetError(err, nil, "")
	}
	return resp, nil
}

func (c *Client) doGet(ctx context.Context, u url.URL) (*http.Response, *AkismetError) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, newAkismetError(err, nil, "")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newAkismetError(err, nil, "")
	}
	return resp, nil
}
