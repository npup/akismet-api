package akismet

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// The result of a comment check.
type CheckCommentResult struct {
	IsSpam       bool
	Discard      bool           // true if Akismet considers the comment blatant spam that can be discarded without saving
	RecheckAfter *time.Duration // if set, Akismet requests a recheck after this duration; resubmit with RecheckReason set to "recheck"
	AkismetGUID  string         // unique identifier for this request, useful when contacting Akismet support
}

// String returns a pretty-printed JSON representation of the CheckResult.
func (r *CheckCommentResult) String() string {
	out, _ := json.MarshalIndent(r, "", "  ")
	return string(out)
}

// Checks whether a comment is spam. Returns an error if the request
// fails or Akismet reports a problem with the submitted fields.
func (c *Client) CheckComment(ctx context.Context, comment *Comment) (*CheckCommentResult, *AkismetError) {
	// POST comment data
	postBody := c.commentValues(comment)
	resp, akismetErr := c.doPost(ctx, c.endpoints.CheckComment, postBody)
	if akismetErr != nil {
		// performing request went wrong
		return nil, akismetErr
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// reading body went wrong
		return nil, AkismetErrorFromResponse(err, resp)
	}

	bodyStr := strings.TrimSpace(string(responseBody))
	if bodyStr != BODY_SPAM_RESPONSE && bodyStr != BODY_HAM_RESPONSE {
		// unexcpected response body
		err := fmt.Errorf("unexpected response body (not `%s` or `%s`) [http status:%s]", BODY_SPAM_RESPONSE, BODY_HAM_RESPONSE, resp.Status)
		return nil, AkismetErrorFromResponse(err, resp)
	}

	result := &CheckCommentResult{
		IsSpam:      bodyStr == BODY_SPAM_RESPONSE,
		Discard:     resp.Header.Get(AkismetHeaders.ProTip) == HEADER_PROTIP_DISCARD_RESPONSE,
		AkismetGUID: resp.Header.Get(AkismetHeaders.GUID),
	}

	if s := resp.Header.Get(AkismetHeaders.RecheckAfter); s != "" {
		if secs, err := strconv.Atoi(s); err == nil {
			d := time.Duration(secs) * time.Second
			result.RecheckAfter = &d
		}
	}

	return result, nil
}
