package siplib

// Testing Generation of a UDP OPTIONS request.

func SIPGenTest() string {
	z := SIPRequest{}
	z.Init("UDP", "192.168.0.20", "OPTIONS", 2000)
	z.DefaultHeaders()
	z.SetContactHeaders("blahblahblah", 666)
	return z.Generate()
}
