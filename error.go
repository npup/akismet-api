package akismet

import (
	"encoding/json"
	"net/http"
)

type AkismetError struct {
	Err       error
	DebugHelp string
	Alert     *Alert
}

// Buils an AkismetError from an http.Response, extracting
// alert and debug help info from the headers if available.
// return nil if none of that info was present
func akismetErrorFromResponse(err error, resp *http.Response) *AkismetError {
	// alert message from akismet headers
	alert := getAlert(resp)
	// debughelp really not available unless reponse is invalid
	// (not true|false) but picking it up here anyway
	debugHelp := getDebugHelp(resp)
	if alert == nil && debugHelp == "" && err == nil {
		return nil
	}
	return newAkismetError(err, alert, debugHelp)
}

func newAkismetError(err error, alert *Alert, debugHelp string) *AkismetError {
	return &AkismetError{
		Err:       err,
		DebugHelp: debugHelp,
		Alert:     alert,
	}
}

func (e *AkismetError) Error() string {
	if e.Err == nil && e.DebugHelp == "" && e.Alert == nil {
		return "akismet: unknown error"
	}
	tmp := struct {
		Err       string `json:"err,omitempty"`
		DebugHelp string `json:"debugHelp,omitempty"`
		Alert     *Alert `json:"alert,omitempty"`
	}{
		DebugHelp: e.DebugHelp,
		Alert:     e.Alert,
	}
	if e.Err != nil {
		tmp.Err = e.Err.Error()
	}
	out, _ := json.Marshal(tmp)
	return string(out)
}

func (e *AkismetError) Unwrap() error {
	return e.Err
}
