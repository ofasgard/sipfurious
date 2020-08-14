package siplib

import "net"

func UDPHandleACK(conn *net.UDPConn, req SIPRequest, res SIPResponse) error {
	ack := req
	ack.Method = "ACK"
	ack.URI = GenerateURI(ack.Host, "")
	ack.Headers["Cseq"] = "1 ACK"
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
