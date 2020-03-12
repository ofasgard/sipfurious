package siplib

import "fmt"

// Struct used to keep track of a SIP request.

type SIPRequest struct {
	Proto string
	Host string
	Method string
	Extension int
	URI string
	PreHeaders map[string]string
	Headers map[string]string
	PostHeaders map[string]string
	Body string
}

// Just set the request up with basic info.

func (r *SIPRequest) Init(proto string, host string, method string, extension int) {
	r.Proto = proto
	r.Host = host
	r.Method = method
	r.Extension = extension
	r.URI = GenerateURI(host, method, extension)
	r.PreHeaders = make(map[string]string)
	r.Headers = make(map[string]string)
	r.PostHeaders = make(map[string]string)
}

// Initialise some default headers suitable for an OPTIONS request - these won't be right for every request!

func (r *SIPRequest) DefaultHeaders() {
	r.Headers["Accept"] = "application/sdp"
	r.Headers["User-Agent"] = "gossiper-scanner"
	r.Headers["Max-Forwards"] = "70"
	r.Headers["Cseq"] = fmt.Sprintf("1 %s", r.Method)
	r.Headers["Content-Length"] = fmt.Sprintf("%d", len(r.Body))
	r.Headers["Call-ID"] = random_number_string(24)
	r.Headers["From"] = fmt.Sprintf("\"gossiper\"<sip:100@1.1.1.1>;tag=%s", random_number_string(46))
	r.Headers["To"] = "\"gossiper\"<sip:100@1.1.1.1>"
}

// Set the Via and Contact headers with your local host/port (or the one you want to spoof).

func (r *SIPRequest) SetContactHeaders(srchost string, srcport int) {
	branch_id := random_number_string(10)
	r.PreHeaders["Via"] = fmt.Sprintf("SIP/2.0/%s %s:%d;branch=z9hG4bK-%s;rport", r.Proto, srchost, srcport, branch_id)
	r.Headers["Contact"] = fmt.Sprintf("sip:%d@%s:%d", r.Extension, srchost, srcport)
}

// Set the From and To headers manually with non-default values.

func (r *SIPRequest) SetRecipients(from_name string, from_uri string, to_name string, to_uri string) {
	r.Headers["From"] = fmt.Sprintf("\"%s\"<%s>;tag=%s", from_name, from_uri, random_number_string(46))
	r.Headers["To"] = fmt.Sprintf("\"%s\"<%s>", to_name, to_uri)
}

// Include a requestbody and recalculate Content-Length header.

func (r *SIPRequest) SetBody (body string) {
	r.Body = body
	r.Headers["Content-Length"] = fmt.Sprintf("%d", len(r.Body))
}

// Generate the request.

func (r SIPRequest) Generate() string {
	output := fmt.Sprintf("%s %s SIP/2.0\r\n", r.Method, r.URI)
	for name,value := range r.PreHeaders {
		output += fmt.Sprintf("%s: %s\r\n", name, value)
	}
	for name,value := range r.Headers {
		output += fmt.Sprintf("%s: %s\r\n", name, value)
	}
	for name,value := range r.PostHeaders {
		output += fmt.Sprintf("%s: %s\r\n", name, value)
	}
	output += "\r\n"
	output += r.Body
	return output
}

// Helper function to generate a SIP URI from host, method and extension.

func GenerateURI(host string, method string, extension int) string {
	if (extension == 0) || (method == "REGISTER") {
		return fmt.Sprintf("sip:%s", host)
	} else {
		return fmt.Sprintf("sip:%d@%s", extension, host)
	}
}

