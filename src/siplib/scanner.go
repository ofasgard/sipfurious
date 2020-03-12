package siplib

import "time"

// High-level function to do an OPTIONS scan over UDP.

func SIPOptionsUDP(target string, port int, timeout int) (string,error) {
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
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)
	resp,err := RequestUDP(conn, req)
	if err != nil {
		return "",err
	}
	//parse the response
	parsed := ParseResponse(resp)
	if val, ok := parsed.Headers["Server"]; ok {
		return val,nil
	} else {
		return "[NONE]",nil
	}
}
