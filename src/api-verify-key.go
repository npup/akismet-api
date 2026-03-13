package akismet

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// Calls the Akismet verify-key endpoint and returns an error if the key is not valid.
func (c *Client) verifyKey(ctx context.Context) *AkismetError {
	postBody := url.Values{
		"key":  {c.apiKey},
		"blog": {c.blogURL},
	}
	resp, akismetErr := c.doPost(ctx, c.endpoints.VerifyKey, postBody)
	if akismetErr != nil {
		return akismetErr
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return AkismetErrorFromResponse(err, resp)
	}

	if strings.TrimSpace(string(responseBody)) != BODY_VALID_MESSAGE {
		return AkismetErrorFromResponse(fmt.Errorf("akismet: invalid API key"), resp)
	}

	return nil
}
