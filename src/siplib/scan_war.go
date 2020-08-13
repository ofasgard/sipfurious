package siplib

import "fmt"
import "strings"
import "time"
import "math/rand"

// High-level function to do an INVITE wardial over UDP.

func WarInviteUDP(target string, port int, timeout int, throttle int, extensions []string) (map[string]string,error) {
	output := make(map[string]string)
	//check a random extension to get "bad" result
	rand.Seed(time.Now().UnixNano())
	bad_ext := fmt.Sprintf("%d", rand.Intn(999999))
	bad_res,err := InviteUDP(target, port, timeout, bad_ext)
	if err != nil {
		return output,err
	}
	bad_status := bad_res.Status
	//now we can begin bruteforcing
	for _,extension := range extensions {
		res,err := InviteUDP(target, port, timeout, extension)
		if err != nil {
			return output,err
		}
		if res.Status != bad_status {
			if strings.HasPrefix(res.Status, AUTHREQ) || strings.HasPrefix(res.Status, PROXYAUTHREQ) || strings.HasPrefix(res.Status, INVALIDPASS) {
				output[extension] = "AUTHREQ"
			} else if strings.HasPrefix(res.Status, OKAY) {
				output[extension] = "NOAUTH"
			} else {
				output[extension] = "WEIRD"
			}
		}
		time.Sleep(time.Duration(throttle) * time.Millisecond) //throttle to prevent flooding
	}
	return output,nil
}

func InviteUDP(target string, port int, timeout int, extension string) (SIPResponse,error) {
	//connect to server
	conn,err := ConnectUDP(target, port)
	if err != nil {
		return SIPResponse{},err
	}
	defer conn.Close()
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "INVITE", extension)
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 5060)
	recipient_uri := GenerateURI(target, extension)
	req.SetRecipients(extension, recipient_uri, extension, recipient_uri)
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
		parsed,err := ParseResponse(resp)
		if err != nil {
			return SIPResponse{},err
		}
		if val,ok := parsed.Headers["Call-ID"]; ok {
			//check if an ACK is needed
			if (parsed.StatusCode >= 200) && (parsed.StatusCode < 699) {
				ack := req
				ack.Method = "ACK"
				ack.URI = GenerateURI(ack.Host, "")
				ack.Headers["Cseq"] = "1 ACK"
				SendUDP(conn, ack)
			}
			//check if a BYE is needed
			if (parsed.StatusCode == 200) {
				bye := req
				bye.Method = "BYE"
				bye.URI = GenerateURI(bye.Host, "")
				bye.Headers["Cseq"] = "2 BYE"
				SendUDP(conn, bye)
			}
			if val == call_id {
				//return the SIPResponse
				return parsed,nil
			}
		}
	}
}

// High-level function to do an INVITE wardial over TCP.

func WarInviteTCP(target string, port int, timeout int, throttle int, extensions []string) (map[string]string,error) {
	output := make(map[string]string)
	//check a random extension to get "bad" result
	rand.Seed(time.Now().UnixNano())
	bad_ext := fmt.Sprintf("%d", rand.Intn(999999))
	bad_res,err := InviteTCP(target, port, timeout, bad_ext)
	if err != nil {
		return output,err
	}
	bad_status := bad_res.Status
	//now we can begin bruteforcing
	for _,extension := range extensions {
		res,err := InviteTCP(target, port, timeout, extension)
		if err != nil {
			return output,err
		}
		if res.Status != bad_status {
			if strings.HasPrefix(res.Status, AUTHREQ) || strings.HasPrefix(res.Status, PROXYAUTHREQ) || strings.HasPrefix(res.Status, INVALIDPASS) {
				output[extension] = "AUTHREQ"
			} else if strings.HasPrefix(res.Status, OKAY) {
				output[extension] = "NOAUTH"
			} else {
				output[extension] = "WEIRD"
			}
		}
		time.Sleep(time.Duration(throttle) * time.Millisecond) //throttle to prevent flooding
	}
	return output,nil
}

func InviteTCP(target string, port int, timeout int, extension string) (SIPResponse,error) {
	//connect to server
	conn,err := ConnectTCP(target, port)
	if err != nil {
		return SIPResponse{},err
	}
	defer conn.Close()
	//generate the request
	req := SIPRequest{}
	req.Init("TCP", target, "INVITE", extension)
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 5060)
	recipient_uri := GenerateURI(target, extension)
	req.SetRecipients(extension, recipient_uri, extension, recipient_uri)
	//make the request
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)
	err = SendTCP(conn, req)
	if err != nil {
		return SIPResponse{},err
	}
	//receive and parse responses until we get one that matches or time out
	call_id := req.Headers["Call-ID"]
	for {
		resp,err := RecvTCP(conn)
		if err != nil {
			return SIPResponse{},err
		}
		parsed,err := ParseResponse(resp)
		if err != nil {
			return SIPResponse{},err
		}
		if val,ok := parsed.Headers["Call-ID"]; ok {
			//check if an ACK is needed
			if (parsed.StatusCode >= 200) && (parsed.StatusCode < 699) {
				ack := req
				ack.Method = "ACK"
				ack.URI = GenerateURI(ack.Host, "")
				ack.Headers["Cseq"] = "1 ACK"
				SendTCP(conn, ack)
			}
			//check if a BYE is needed
			if (parsed.StatusCode == 200) {
				bye := req
				bye.Method = "BYE"
				bye.URI = GenerateURI(bye.Host, "")
				bye.Headers["Cseq"] = "2 BYE"
				SendTCP(conn, bye)
			}
			if val == call_id {
				//return the SIPResponse
				return parsed,nil
			}
		}
	}
}


