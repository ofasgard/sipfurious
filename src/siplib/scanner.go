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

func WarInviteUDP(target string, port int, timeout int, extensions []string) (map[string]string,error) {
	output := make(map[string]string)
	//check a random extension to get "bad" result
	rand.Seed(time.Now().UnixNano())
	bad_ext := fmt.Sprintf("%d", 20000 + rand.Intn(1000))
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
	resp,err := RecvUDP(conn)
	if err != nil {
		return output,err
	}
	//parse the response
	parsed := ParseResponse(resp)
	output = parsed
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
	return output,nil
}

