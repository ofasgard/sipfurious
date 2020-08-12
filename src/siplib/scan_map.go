package siplib

import "time"

// High-level function to do an OPTIONS check over UDP.

func MapUDP(target string, port int, timeout int) (string,error) {
	res,err := OptionsUDP(target, port, timeout)
	if err != nil {
		return "",err
	}
	if val,ok := res.Headers["Server"]; ok {
		return val,nil
	}
	if val,ok := res.Headers["User-Agent"]; ok {
		return val,nil
	}
	return "[NONE]",nil
}

func OptionsUDP(target string, port int, timeout int) (SIPResponse, error) {
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "OPTIONS", "100")
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 5060)
	//connect to server
	conn,err := ConnectUDP(target, port)
	if err != nil {
		return SIPResponse{},err
	}
	defer conn.Close()
	//make the request
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)
	err = SendUDP(conn, req)
	if err != nil {
		return SIPResponse{},err
	}
	//receive and parse responses until we get one that matches or time out
	call_id := req.Headers["Call-ID"]
	for {
		resp,err := RecvUDP(conn)
		if err != nil {
			return SIPResponse{},err
		}
		parsed := ParseResponse(resp)
		if val,ok := parsed.Headers["Call-ID"]; ok {
			if val == call_id {
				return parsed,nil
			}
		}
	}
}
