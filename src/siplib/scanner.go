package siplib

import "fmt"
import "time"
import "strings"
import "math/rand"

// High-level function to do an OPTIONS check over UDP.

func MapUDP(target string, port int, timeout int) (string,error) {
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "OPTIONS", "100")
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
	//receive and parse responses until we get one that matches or time out
	call_id := req.Headers["Call-ID"]
	for {
		resp,err := RecvUDP(conn)
		if err != nil {
			return "",err
		}
		parsed := ParseResponse(resp)
		if val,ok := parsed.Headers["Call-ID"]; ok {
			if val == call_id {
				//check and return the server header
				if val,ok := parsed.Headers["Server"]; ok {
					return val,nil
				} else {
					return "[NONE]",nil
				}
			}
		}
	}
}

// High-level function to do an INVITE wardial over UDP.

func WarInviteUDP(target string, port int, timeout int, throttle bool, extensions []string) (map[string]string,error) {
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
	for index,extension := range extensions {
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
		if (index % 100 == 0) && throttle {
			time.Sleep(1 * time.Second) //throttle to prevent flooding
		}
	}
	return output,nil
}

func InviteUDP(target string, port int, timeout int, extension string) (SIPResponse,error) {
	output := SIPResponse{}
	//connect to server
	conn,err := ConnectUDP(target, port)
	if err != nil {
		fmt.Println(err)
		return output,err
	}
	defer conn.Close()
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "INVITE", extension)
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 5060)
	recipient_uri := GenerateURI(target, "INVITE", extension)
	req.SetRecipients(extension, recipient_uri, extension, recipient_uri)
	//make the request
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)
	err = SendUDP(conn, req)
	if err != nil {
		return output,err
	}
	//receive and parse responses until we get one that matches or time out
	call_id := req.Headers["Call-ID"]
	for {
		resp,err := RecvUDP(conn)
		if err != nil {
			return output,err
		}
		parsed := ParseResponse(resp)
		if val,ok := parsed.Headers["Call-ID"]; ok {
			//check if an ACK is needed
			if (parsed.StatusCode >= 200) && (parsed.StatusCode < 699) {
				ack := req
				ack.Method = "ACK"
				ack.URI = GenerateURI(ack.Host, ack.Method, "")
				ack.Headers["Cseq"] = "1 ACK"
				SendUDP(conn, ack)
			}
			//check if a BYE is needed
			if (parsed.StatusCode == 200) {
				bye := req
				bye.Method = "BYE"
				bye.URI = GenerateURI(bye.Host, bye.Method, "")
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

