package akismet

import (
	"fmt"

	"net/http"
	"net/url"
)

func debugHeaders(resp *http.Response) {
	for k, v := range resp.Header {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func getDebugHelp(resp *http.Response) string {
	return resp.Header.Get(AkismetHeaders.DebugHelp)
}

func getAlert(resp *http.Response) *Alert {
	code := -1
	if codeStr := resp.Header.Get(AkismetHeaders.AlertCode); codeStr != "" {
		fmt.Sscanf(codeStr, "%d", &code)
	}
	message := resp.Header.Get(AkismetHeaders.AlertMsg)
	if code == -1 || message == "" {
		return nil
	}
	return NewAlert(code, message)
}

func getUrlWithParameters(endpoint string, params url.Values) url.URL {
	url, _ := url.Parse(endpoint)
	url.RawQuery = params.Encode()
	return *url
}
