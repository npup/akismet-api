package akismet

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ApiEndpoints = struct {
	VerifyKey    string
	CheckComment string
	SubmitSpam   string
	SubmitHam    string
	KeySites     string
	UsageLimit   string
}

// Holds the API credentials and constructed URLs for making requests.
type Client struct {
	apiKey     string
	blogURL    string
	httpClient *http.Client
	endpoints  ApiEndpoints
}

// Creates a new Akismet client and verifies the API key before returning.
// Returns an error if the key is invalid or the verification request fails.
func NewClient(ctx context.Context, apiKey, blogURL string) (*Client, *AkismetError) {
	akismetBaseURL := "https://rest.akismet.com"
	return newClientWithApiBaseURL(ctx, apiKey, blogURL, akismetBaseURL)
}

func newClientWithApiBaseURL(ctx context.Context, apiKey, blogURL string, apiBaseURL string) (*Client, *AkismetError) {
	endpoints := ApiEndpoints{
		VerifyKey:    fmt.Sprintf("%s/%s/verify-key", apiBaseURL, "1.1"),
		CheckComment: fmt.Sprintf("%s/%s/comment-check", apiBaseURL, "1.1"),
		SubmitSpam:   fmt.Sprintf("%s/%s/submit-spam", apiBaseURL, "1.1"),
		SubmitHam:    fmt.Sprintf("%s/%s/submit-ham", apiBaseURL, "1.1"),
		KeySites:     fmt.Sprintf("%s/%s/key-sites", apiBaseURL, "1.2"),
		UsageLimit:   fmt.Sprintf("%s/%s/usage-limit", apiBaseURL, "1.2"),
	}
	client := &Client{apiKey: apiKey, blogURL: blogURL, httpClient: &http.Client{}, endpoints: endpoints}
	if err := client.verifyKey(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

// Converts a Comment into form values for the Akismet API, omitting empty fields.
func (c *Client) commentValues(comment *Comment) url.Values {
	v := url.Values{
		"api_key":    {c.apiKey},
		"blog":       {c.blogURL},
		"user_ip":    {comment.UserIP},
		"user_agent": {comment.UserAgent},
	}

	if comment.Referrer != "" {
		v.Set("referrer", comment.Referrer)
	}
	if comment.Permalink != "" {
		v.Set("permalink", comment.Permalink)
	}
	if comment.Type != "" {
		v.Set("comment_type", string(comment.Type))
	}
	if comment.Author != "" {
		v.Set("comment_author", comment.Author)
	}
	if comment.AuthorEmail != "" {
		v.Set("comment_author_email", comment.AuthorEmail)
	}
	if comment.AuthorURL != "" {
		v.Set("comment_author_url", comment.AuthorURL)
	}
	if comment.Content != "" {
		v.Set("comment_content", comment.Content)
	}
	if !comment.DateGMT.IsZero() {
		v.Set("comment_date_gmt", comment.DateGMT.UTC().Format(time.RFC3339))
	}
	if !comment.PostModifiedGMT.IsZero() {
		v.Set("comment_post_modified_gmt", comment.PostModifiedGMT.UTC().Format(time.RFC3339))
	}
	if comment.BlogLang != "" {
		v.Set("blog_lang", comment.BlogLang)
	}
	if comment.BlogCharset != "" {
		v.Set("blog_charset", comment.BlogCharset)
	}
	if comment.UserRole != "" {
		v.Set("user_role", comment.UserRole)
	}
	if comment.IsTest {
		v.Set("is_test", "1")
	}
	if comment.CommentParent != "" {
		v.Set("comment_parent", comment.CommentParent)
	}
	if comment.RecheckReason != "" {
		v.Set("recheck_reason", comment.RecheckReason)
	}
	if comment.HoneypotFieldName != "" {
		v.Set("honeypot_field_name", comment.HoneypotFieldName)
		v.Set(comment.HoneypotFieldName, comment.HoneypotFieldValue)
	}
	for _, ctx := range comment.PostContextTags {
		v.Add("comment_context[]", ctx)
	}
	return v
}
