package siplib

import "net"

// Functions to handle responses that need an ACK/BYE to be generated and sent, used by all modules.
// To avoid a race condition, we should only handle ACK/BYE for the Call-ID we currently care about.
// Superfluous responses from old Call-IDs can be safely ignored and left to time out.

func UDPHandleACK(conn *net.UDPConn, req SIPRequest, res SIPResponse) error {
	ack := req
	ack.Method = "ACK"
	ack.URI = GenerateURI(ack.Host, "")
	ack.Headers["Cseq"] = "1 ACK"
	ack.Headers["To"] = res.Headers["To"]
	err := SendUDP(conn, ack)
	return err
}

func UDPHandleBYE(conn *net.UDPConn, req SIPRequest, res SIPResponse) error {
	bye := req
	bye.Method = "BYE"
	bye.URI = GenerateURI(bye.Host, "")
	bye.Headers["Cseq"] = "2 BYE"
	err := SendUDP(conn, bye)
	return err
}

func TCPHandleACK(conn *net.TCPConn, req SIPRequest, res SIPResponse) error {
	ack := req
	ack.Method = "ACK"
	ack.URI = GenerateURI(ack.Host, "")
	ack.Headers["Cseq"] = "1 ACK"
	ack.Headers["To"] = res.Headers["To"]
	err := SendTCP(conn, ack)
	return err
}

func TCPHandleBYE(conn *net.TCPConn, req SIPRequest, res SIPResponse) error {
	bye := req
	bye.Method = "BYE"
	bye.URI = GenerateURI(bye.Host, "")
	bye.Headers["Cseq"] = "2 BYE"
	err := SendTCP(conn, bye)
	return err
}
