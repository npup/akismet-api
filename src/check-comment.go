package akismet

import (
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
func (c *Client) CheckComment(comment *Comment) (*CheckCommentResult, *AkismetError) {
	//endpoint := fmt.Sprintf("%s/comment-check", c.keyedBaseURL)
	//endpoint := fmt.Sprintf("%s/comment-check", baseURL)
	resp, err := c.httpClient.PostForm(ApiEndpoints.CheckComment, c.commentValues(comment))
	if err != nil {
		// things went totally wrong
		return nil, NewAkismetError(err, nil, "")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// reading body went totally wrong
		return nil, AkismetErrorFromResponse(err, resp)
	}

	bodyStr := strings.TrimSpace(string(body))
	if bodyStr != "true" && bodyStr != "false" {
		// unexpected response body
		err := fmt.Errorf("unexpected response body (not `true` or `false`) [http status:%s]", resp.Status)
		return nil, AkismetErrorFromResponse(err, resp)
	}

	result := &CheckCommentResult{
		IsSpam:      strings.TrimSpace(string(body)) == "true",
		Discard:     resp.Header.Get(AkismetHeaders.ProTip) == "discard",
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
