package akismet

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// Reports a comment that was not caught by CheckComment (false negative).
func (c *Client) ReportSpam(ctx context.Context, comment *Comment) *AkismetError {
	return c.report(ctx, c.endpoints.SubmitSpam, comment)
}

// Reports a comment that was wrongly flagged as spam by CheckComment (false positive).
func (c *Client) ReportHam(ctx context.Context, comment *Comment) *AkismetError {
	return c.report(ctx, c.endpoints.SubmitHam, comment)
}

// Shared implementation for ReportSpam and ReportHam.
func (c *Client) report(ctx context.Context, endpoint string, comment *Comment) *AkismetError {
	resp, akismetErr := c.doPost(ctx, endpoint, c.commentValues(comment))
	if akismetErr != nil {
		return akismetErr
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return akismetErrorFromResponse(err, resp)
	}

	if strings.TrimSpace(string(body)) != bodyReportSuccessMessage {
		err := fmt.Errorf("akismet: unexpected response: %s", strings.TrimSpace(string(body)))
		return akismetErrorFromResponse(err, resp)
	}

	return nil
}
