package main

import "siplib"
import "fmt"

func main() {
	//generate the request
	z := siplib.SIPRequest{}
	z.Init("UDP", "192.168.0.20", "OPTIONS", 2000)
	z.DefaultHeaders()
	z.SetContactHeaders("blahblahblah", 666)
	//print to screen
	request := z.Generate()
	fmt.Println(request)
	//connect to server
	conn,err := siplib.ConnectUDP("192.168.0.20", 5060)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	//make request
	resp,err := siplib.RequestUDP(conn, z)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

