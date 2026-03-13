package akismet

import (
	"fmt"
	"io"
	"net/url"
	"strings"
)

// Calls the Akismet verify-key endpoint and returns an error if the key is not valid.
func (c *Client) verifyKey() error {
	endpoint := fmt.Sprintf("%s/verify-key", akismetBaseURL)
	resp, err := c.httpClient.PostForm(endpoint, url.Values{
		"key":  {c.apiKey},
		"blog": {c.blogURL},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(body)) != "valid" {
		return fmt.Errorf("akismet: invalid API key")
	}

	return nil
}
