package akismet

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

// text descriptions pertaining to alert codes
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

// Constant messages that the API responds with - instead of
// proper http status codes in some cases :-(
const BODY_REPORT_SUCCESS_MESSAGE = "Thanks for making the web a better place."
const BODY_INVALID_MESSAGE = "invalid"
const BODY_VALID_MESSAGE = "valid"
const BODY_SPAM_RESPONSE = "true"
const BODY_HAM_RESPONSE = "false"
const HEADER_PROTIP_DISCARD_RESPONSE = "discard"
const PROP_USAGE_LIMIT_NO_LIMIT = "none"
