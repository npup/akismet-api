package akismet

import "fmt"

type Alert struct {
	Code        int
	Message     string
	Description string
}

func NewAlert(code int, message string) *Alert {
	descr, ok := AlertDescriptionsByCode[code]
	if !ok {
		descr = fmt.Sprintf("No descr for alert code %d", code)
	}
	return &Alert{
		Code:        code,
		Message:     message,
		Description: descr,
	}
}
