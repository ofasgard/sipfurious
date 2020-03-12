package main

import "siplib"
import "fmt"
import "flag"
import "os"
import "strconv"
import "text/tabwriter"

func main() {
	//parse flags
	flag.Usage = usage
	timeout_ptr := flag.Int("timeout", 10, "")
	flag.Parse()
	if flag.NArg() != 4 {
		usage()
		return
	}
	method := flag.Arg(0)
	protocol := flag.Arg(1)
	target := flag.Arg(2)
	timeout := *timeout_ptr
	//validate flags
	port,err := strconv.Atoi(flag.Arg(3))
	if (err != nil) || (port < 1) {
		usage()
		return
	}
	//defer to the correct function based on arguments
	switch protocol {
		case "udp":
			switch method{
				case "map":
					map_udp(target, port, timeout)
				case "war":
					fmt.Fprintf(os.Stderr, "Wardialing is not yet implemented.\n")
					return
				case "crack":
					fmt.Fprintf(os.Stderr, "Cracking is not yet implemented.\n")
					return
				default:
					usage()
					return
			}
		case "tcp":
			fmt.Fprintf(os.Stderr, "TCP is not yet implemented.\n")
			return
		case "tls":
			fmt.Fprintf(os.Stderr, "TLS is not yet implemented.\n")
			return
		default:
			usage()
			return
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <map|war|crack> <udp|tcp|tls> <target> <port>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "'map': Scanner that uses OPTIONS to attempt to retrieve the SIP Server header.\n")
	fmt.Fprintf(os.Stderr, "'war': Wardialler that bruteforces extensions using various SIP methods.\n")
	fmt.Fprintf(os.Stderr, "'crack': Bruteforcer to crack SIP passwords for an extension.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Optional arguments:\n")
	w := new(tabwriter.Writer)
	w.Init(os.Stderr, 0, 8, 2, '\t', 0)
	fmt.Fprintf(w, "\t--timeout <sec>\tTimeout (in seconds) for each request. [DEFAULT: 10]\n")
	w.Flush()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\n\nExample: %s map udp 192.168.0.20\n", os.Args[0])
}


func map_udp(target string, port int, timeout int) {
	result,err := siplib.SIPOptionsUDP(target, port, timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not map %s:%d over UDP:\n\t%s\n", target, port, err)
		return
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintf(w, "Target\tPort\tServer Header\n")
	fmt.Fprintf(w, "\t\t\t\n")
	fmt.Fprintf(w, "%s\t%d\t%s\n", target, port, result)
	w.Flush()
}


