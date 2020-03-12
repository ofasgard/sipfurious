package siplib

import "fmt"
import "time"
import "strings"

// High-level function to do an OPTIONS check over UDP.

func MapUDP(target string, port int, timeout int) (string,error) {
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "OPTIONS", 100)
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 5060)
	//connect to server
	conn,err := ConnectUDP(target, port)
	if err != nil {
		return "",err
	}
	defer conn.Close()
	//make the request
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)
	err = SendUDP(conn, req)
	if err != nil {
		return "",err
	}
	resp,err := RecvUDP(conn)
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

// High-level function to do an INVITE wardial over UDP.

func WarInviteUDP(target string, port int, timeout int, extensions []int) (map[int]string,error) {
	output := make(map[int]string)
	for _,extension := range extensions {
		//connect to server
		conn,err := ConnectUDP(target, port)
		if err != nil {
			fmt.Println(err)
			return output,err
		}
		//generate the request
		req := SIPRequest{}
		req.Init("UDP", target, "INVITE", extension)
		req.DefaultHeaders()
		req.SetContactHeaders("1.1.1.1", 5060)
		recipient_name := fmt.Sprintf("%d", extension)
		recipient_uri := GenerateURI(target, "INVITE", extension)
		req.SetRecipients(recipient_name, recipient_uri, recipient_name, recipient_uri)
		//make the request
		deadline := time.Now().Add(time.Duration(timeout) * time.Second)
		conn.SetDeadline(deadline)
		err = SendUDP(conn, req)
		if err != nil {
			return output,err
		}
		resp,err := RecvUDP(conn)
		if err != nil {
			return output,err
		}
		//parse the response
		parsed := ParseResponse(resp)
		if strings.HasPrefix(parsed.Status, AUTHREQ) || strings.HasPrefix(parsed.Status, PROXYAUTHREQ) || strings.HasPrefix(parsed.Status, INVALIDPASS) {
			output[extension] = "AUTHREQ"
		}
		if strings.HasPrefix(parsed.Status, OKAY) {
			output[extension] = "NOAUTH"
		}
		//handle ACK and BYE requests to gracefully terminate the connection
		for {
			if (parsed.StatusCode >= 200) && (parsed.StatusCode < 699) {
				//we need to send an ACK
				req.Method = "ACK"
				req.URI = GenerateURI(req.Host, req.Method, req.Extension)
				req.Headers["Cseq"] = "1 ACK"
				SendUDP(conn, req)
			}
			if (parsed.StatusCode == 200) {	
				//we need to send a BYE
				req.Method = "BYE"
				req.URI = GenerateURI(req.Host, req.Method, req.Extension)
				req.Headers["Cseq"] = "2 BYE"
				SendUDP(conn, req)
				RecvUDP(conn)
				break
			}
			resp,err = RecvUDP(conn)
			if err != nil {break}
			if len(resp) == 0 {break}
			parsed = ParseResponse(resp)
		}
		conn.Close()
	}
	return output,nil
}
