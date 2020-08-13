package siplib

import "strings"
import "strconv"
import "errors"

// Struct used to keep track of a SIP response.

type SIPResponse struct {
	Status string
	StatusCode int
	Headers map[string]string
	Body string
}

// Parse a raw response string; return a SIPResponse object.

func ParseResponse(resp string) (SIPResponse,error) {
	output := SIPResponse{}
	output.Headers = make(map[string]string)
	//split up the headers from the body
	parts := strings.Split(resp, "\r\n\r\n")
	if len(parts) < 2 {
		return output,errors.New("Received a SIP response we couldn't parse")
	}
	header_part := parts[0]
	output.Body = parts[1]
	//extract the status line
	headers := strings.Split(header_part, "\r\n")
	if len(headers) < 2 {
		return output,errors.New("Received a SIP response we couldn't parse")
	}
	output.Status = headers[0]
	//attempt to extract a response code from the status line
	status_parts := strings.Split(output.Status, " ")
	if len(status_parts) < 2 {
		return output,errors.New("Received a SIP response we couldn't parse")
	}
	code,err := strconv.Atoi(status_parts[1])
	if err != nil {
		return output,err
	}
	output.StatusCode = code
	//parse each line and identify any valid headers
	for _,header := range headers[1:] {
		header_parts := strings.Split(header, ": ")
		if len(header_parts) > 1 {
			output.Headers[header_parts[0]] = strings.Join(header_parts[1:], ": ")
		}
	}
	return output,nil
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
