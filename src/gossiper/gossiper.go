package main

import "siplib"
import "fmt"

func main() {
	result,err := siplib.SIPOptionsUDP("192.168.0.20", 5060, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Enumerated server header: %s\n", result)
}


