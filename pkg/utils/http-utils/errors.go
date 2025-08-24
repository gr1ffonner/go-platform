package httputils

import "strings"

type HTTPErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Message     string `json:"errorMessage,omitempty"`
	Description string `json:"error_description,omitempty"`
}

func (e HTTPErrorResponse) String() string {
	var res strings.Builder
	if len(e.Error) > 0 {
		res.WriteString(e.Error)
	}
	if len(e.Message) > 0 {
		if res.Len() > 0 {
			res.WriteString(": ")
		}
		res.WriteString(e.Message)
	}
	if len(e.Description) > 0 {
		if res.Len() > 0 {
			res.WriteString(": ")
		}
		res.WriteString(e.Description)
	}
	return res.String()
}
