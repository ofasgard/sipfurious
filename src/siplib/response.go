package siplib

import "strings"

// Struct used to keep track of a SIP response.

type SIPResponse struct {
	Status string
	Headers map[string]string
	Body string
}

// Parse a raw response string; return a SIPResponse object.

func ParseResponse(resp string) SIPResponse {
	output := SIPResponse{}
	parts := strings.Split(resp, "\r\n\r\n")
	header_part := parts[0]
	output.Body = parts[1]
	output.Headers = make(map[string]string)
	headers := strings.Split(header_part, "\r\n")
	output.Status = headers[0]
	for _,header := range headers[1:] {
		header_parts := strings.Split(header, ": ")
		output.Headers[header_parts[0]] = header_parts[1]
	}
	return output
}
