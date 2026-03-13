package akismet

// Reports a comment that was not caught by CheckComment (false negative).
func (c *Client) ReportSpam(comment *Comment) *AkismetError {
	return c.report(ApiEndpoints.SubmitSpam, comment)
}

// Reports a comment that was wrongly flagged as spam by CheckComment (false positive).
func (c *Client) ReportHam(comment *Comment) *AkismetError {
	return c.report(ApiEndpoints.SubmitHam, comment)
}

// Shared implementation for ReportSpam and ReportHam.
func (c *Client) report(endpoint string, comment *Comment) *AkismetError {

	resp, err := c.httpClient.PostForm(endpoint, c.commentValues(comment))
	if err != nil {
		return NewAkismetError(err, nil, "")
	}
	defer resp.Body.Close()

	akismetErr := AkismetErrorFromResponse(nil, resp)
	if akismetErr != nil {
		return akismetErr
	}

	return nil
}
