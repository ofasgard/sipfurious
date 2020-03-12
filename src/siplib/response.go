package siplib

import "strings"
import "strconv"

// Struct used to keep track of a SIP response.

type SIPResponse struct {
	Status string
	StatusCode int
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
	if len(output.Status) > 12 {
		code,err := strconv.Atoi(output.Status[8:11])
		if err == nil {
			output.StatusCode = code
		}
	}
	for _,header := range headers[1:] {
		header_parts := strings.Split(header, ": ")
		output.Headers[header_parts[0]] = header_parts[1]
	}
	return output
}

// Constants used to represent SIP response status codes.

const statuslen = 12
const PROXYAUTHREQ = "SIP/2.0 407 "
const AUTHREQ = "SIP/2.0 401 "
const OKAY = "SIP/2.0 200 "
const NOTFOUND = "SIP/2.0 404 "
const INVALIDPASS = "SIP/2.0 403 "
const TRYING = "SIP/2.0 100 "
const RINGING = "SIP/2.0 180 "
const NOTALLOWED = "SIP/2.0 405 "
const UNAVAILABLE = "SIP/2.0 480 "
const DECLINED = "SIP/2.0 603 "
const INEXISTENTTRANSACTION = "SIP/2.0 481"
const BADREQUEST = "SIP/2.0 400 "
const SERVICEUNAVAILABLE = "SIP/2.0 503 "
