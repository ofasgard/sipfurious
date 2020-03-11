package main

import "siplib"
import "fmt"
import "net"

func main() {
	//test test this is just atest
	request := siplib.SIPGenTest()
	fmt.Println(request)
	webserver_addr,err := net.ResolveUDPAddr("udp4", "192.168.0.20:5060")
	if err != nil {
		fmt.Println(err)
		return
	}
	local_addr, err := net.ResolveUDPAddr("udp4", ":5060")
	if err != nil {
		fmt.Println(err)
		return
	}
	webserver_conn,err := net.DialUDP("udp4", local_addr, webserver_addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := []byte(request)
	_,err = webserver_conn.Write(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp := make([]byte, 2048)
	_,err = webserver_conn.Read(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resp))
}

