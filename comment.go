package akismet

import "time"

// Represents a comment submission to be checked or reported.
// All fields are optional, but the more you provide the more accurate
// Akismet's evaluation will be. UserIP and UserAgent are especially valuable.
type Comment struct {
	Content            string      // the text of the comment being submitted
	Type               CommentType // type of content being submitted, see CommentType constants
	IsTest             bool        // marks the request as a test, does not affect Akismet's training data
	UserIP             string      // IP address of the commenter's request
	UserAgent          string      // value of the HTTP User-Agent header from the commenter's request; do not confuse with the user agent of your Akismet library
	Author             string      // display name submitted by the comment author
	AuthorEmail        string      // email address provided by the comment author; may be self-reported or verified depending on your authentication setup
	AuthorURL          string      // personal website or homepage URL manually entered by the comment author; do not send automatically generated URLs such as a profile URL on your site
	Referrer           string      // value of the HTTP Referer header from the commenter's request, i.e. the URL they visited before arriving on your site
	Permalink          string      // URL of the page with the post that the comment was submitted on
	DateGMT            time.Time   // UTC timestamp of comment creation
	PostModifiedGMT    time.Time   // UTC timestamp of the publication time for the post, page or thread on which the comment was posted
	BlogLang           string      // language(s) of the site or app, not the comment itself; ISO 639-1 codes, comma-separated e.g. "en-GB, en, fr-CA, sv-SE"
	BlogCharset        string      // character encoding used by the blog, e.g. "UTF-8" or "ISO-8859-1"; helps Akismet correctly interpret submitted content
	UserRole           string      // role of the user submitting the comment; set to "administrator" to bypass spam checking entirely
	RecheckReason      string      // reason for rechecking previously submitted content, e.g. "user updated content", "incorrect previous classification"
	CommentParent      string      // ID of the parent comment in a threaded comment system
	HoneypotFieldName  string      // name of the honeypot form field
	HoneypotFieldValue string      // value of the honeypot form field; must be submitted together with HoneypotFieldName
	PostContextTags    []string    // tags or categories applied to the post the comment was submitted on; do not use values supplied by the commenter
}

// Creates a new Comment with the type and content set.
func NewComment(content string, commentType CommentType) *Comment {
	return &Comment{
		Type:    commentType,
		Content: content,
	}
}

func (c *Comment) WithUserIP(userIP string) *Comment {
	c.UserIP = userIP
	return c
}

func (c *Comment) WithUserAgent(userAgent string) *Comment {
	c.UserAgent = userAgent
	return c
}

func (c *Comment) WithReferrer(referrer string) *Comment {
	c.Referrer = referrer
	return c
}

func (c *Comment) WithPermalink(permalink string) *Comment {
	c.Permalink = permalink
	return c
}

func (c *Comment) WithAuthor(author string) *Comment {
	c.Author = author
	return c
}

func (c *Comment) WithAuthorEmail(authorEmail string) *Comment {
	c.AuthorEmail = authorEmail
	return c
}

func (c *Comment) WithAuthorURL(authorURL string) *Comment {
	c.AuthorURL = authorURL
	return c
}

func (c *Comment) WithDateGMT(dateGMT time.Time) *Comment {
	c.DateGMT = dateGMT
	return c
}

func (c *Comment) WithPostModifiedGMT(postModifiedGMT time.Time) *Comment {
	c.PostModifiedGMT = postModifiedGMT
	return c
}

func (c *Comment) WithBlogLang(blogLang string) *Comment {
	c.BlogLang = blogLang
	return c
}

func (c *Comment) WithBlogCharset(blogCharset string) *Comment {
	c.BlogCharset = blogCharset
	return c
}

func (c *Comment) WithUserRole(userRole string) *Comment {
	c.UserRole = userRole
	return c
}

func (c *Comment) WithIsTest(isTest bool) *Comment {
	c.IsTest = isTest
	return c
}

func (c *Comment) WithRecheckReason(recheckReason string) *Comment {
	c.RecheckReason = recheckReason
	return c
}

func (c *Comment) WithCommentParent(parentID string) *Comment {
	c.CommentParent = parentID
	return c
}

func (c *Comment) WithHoneypot(fieldName string, fieldValue string) *Comment {
	c.HoneypotFieldName = fieldName
	c.HoneypotFieldValue = fieldValue
	return c
}

func (c *Comment) WithPostContextTags(context []string) *Comment {
	c.PostContextTags = context
	return c
}

// Known comment types accepted by the Akismet API.
type CommentType string

const (
	CommentTypeComment     CommentType = "comment"
	CommentTypeForumPost   CommentType = "forum-post"
	CommentTypeBlogPost    CommentType = "blog-post"
	CommentTypeContactForm CommentType = "contact-form"
	CommentTypeSignup      CommentType = "signup"
	CommentTypeMessage     CommentType = "message"
	CommentTypeReply       CommentType = "reply"
	CommentTypeTweet       CommentType = "tweet"
	CommentTypePingback    CommentType = "pingback"
	CommentTypeTrackback   CommentType = "trackback"
)
