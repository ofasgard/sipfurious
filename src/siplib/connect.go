package siplib

import "fmt"
import "net"

// Try to bind to UDP port 5060, since some implementations respond to :5060 regardless of your source port.
// If we can't nab it, use a random source port.

func ConnectUDP(target_host string, target_port int) (*net.UDPConn,error) {
	target := fmt.Sprintf("%s:%d", target_host, target_port)
	target_addr,err := net.ResolveUDPAddr("udp4", target)
	if err != nil {
		return nil,err
	}
	local_addr,err := net.ResolveUDPAddr("udp", ":5060")
	if err != nil {
		return nil,err
	}
	conn,err := net.DialUDP("udp4", local_addr, target_addr)
	if err == nil {
		return conn,nil
	}
	conn,err = net.DialUDP("udp4", nil, target_addr)
	return conn,err
}

// Send a SIPRequest along an open UDP connection and return the response (or an error).

func RequestUDP(conn *net.UDPConn, req SIPRequest) (string,error) {
	msg := req.Generate()
	_,err := fmt.Fprintf(conn, msg)
	if err != nil {
		return "",err
	}
	output := ""
	buf := make([]byte, 512)
	for {
		n,err := conn.Read(buf)
		if err != nil {
			return "",err
		}
		if n == 512 {
			output += string(buf)
		} else {
			output += string(buf[0:n])
			break
		}
	}
	return output,nil
}


