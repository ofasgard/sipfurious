package siplib

import "errors"

// High-level function to do an OPTIONS scan over UDP.

func SIPOptionsUDP(target string, port int) (string,error) {
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "OPTIONS", 100)
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 666)
	//connect to server
	conn,err := ConnectUDP(target, port)
	if err != nil {
		return "",err
	}
	defer conn.Close()
	//make the request
	resp,err := RequestUDP(conn, req)
	if err != nil {
		return "",err
	}
	//parse the response
	parsed := ParseResponse(resp)
	if val, ok := parsed.Headers["Server"]; ok {
		return val,nil
	} else {
		return "",errors.New("No server header returned by SIP server.")
	}
}
