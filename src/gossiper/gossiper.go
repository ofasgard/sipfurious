package main

import "siplib"
import "fmt"

func main() {
	result,err := siplib.SIPOptionsUDP("192.168.0.20", 5060)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Enumerated server header: %s\n", result)
}


