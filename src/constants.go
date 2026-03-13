package akismet

import "fmt"

// Akismet REST API base URLs.
const (
	akismetBaseURL = "https://rest.akismet.com"
)

var ApiEndpoints = struct {
	VerifyKey    string
	CheckComment string
	SubmitSpam   string
	SubmitHam    string
	KeySites     string
	UsageLimit   string
}{
	VerifyKey:    fmt.Sprintf("%s/%s/verify-key", akismetBaseURL, "1.1"),
	CheckComment: fmt.Sprintf("%s/%s/comment-check", akismetBaseURL, "1.1"),
	SubmitSpam:   fmt.Sprintf("%s/%s/submit-spam", akismetBaseURL, "1.1"),
	SubmitHam:    fmt.Sprintf("%s/%s/submit-ham", akismetBaseURL, "1.1"),
	KeySites:     fmt.Sprintf("%s/%s/key-sites", akismetBaseURL, "1.2"),
	UsageLimit:   fmt.Sprintf("%s/%s/usage-limit", akismetBaseURL, "1.2"),
}

// AkismetHeaders contains the Akismet response header names.
var AkismetHeaders = struct {
	DebugHelp    string
	ProTip       string
	RecheckAfter string
	GUID         string
	AlertCode    string
	AlertMsg     string
}{
	DebugHelp:    "X-Akismet-debug-help",
	ProTip:       "X-Akismet-pro-tip",
	RecheckAfter: "X-Akismet-recheck-after",
	GUID:         "X-Akismet-guid",
	AlertCode:    "X-Akismet-alert-code",
	AlertMsg:     "X-Akismet-alert-msg",
}

var AlertDescriptionsByCode = map[int]string{
	10001: "Your site is using an expired Yahoo! Small Business API key.",
	10003: "You must upgrade your Personal subscription to continue using Akismet.",
	10005: "Your Akismet API key may be in use by someone else.",
	10006: "Your subscription has been suspended due to improper use.",
	10009: "Your subscription has been suspended due to overuse.",
	10010: "Your subscription has been suspended due to inappropriate use.",
	10011: "Your subscription needs to be upgraded due to high usage.",
	10402: "Your API key was suspended for non-payment.",
	10403: "The owner of your API key has revoked your site's access to the key.",
	10404: "Your site was not found in the list of sites allowed to use the API key you used.",
	30001: "Your Personal subscription needs to be upgraded based on your usage.",
}
