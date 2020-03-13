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
	timeout_ptr := flag.Int("timeout", 3, "")
	start_ext_ptr := flag.Int("ext-start", 0, "")
	end_ext_ptr := flag.Int("ext-end", 2000, "")
	flag.Parse()
	timeout := *timeout_ptr
	start_ext := *start_ext_ptr
	end_ext := *end_ext_ptr
	//validate args
	if flag.NArg() < 3 {
		usage()
		return
	}
	method := flag.Arg(0)
	protocol := flag.Arg(1)
	targets := parse_target(flag.Arg(2))
	//port argument is optional
	var port int = 5060
	var err error
	if flag.NArg() > 3 {
		port,err = strconv.Atoi(flag.Arg(3))
		if (err != nil) || (port < 1) {
			usage()
			return
		}
	}
	//defer to the correct function based on arguments
	switch protocol {
		case "udp":
			switch method{
				case "map":
					map_udp(targets, port, timeout)
					return
				case "war":
					war_udp(targets, port, timeout, start_ext, end_ext) //todo - make extensions configurable
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
	fmt.Fprintf(os.Stderr, "Usage: %s <map|war|crack> <udp|tcp|tls> <target> [port]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "'map': Scanner that uses OPTIONS to attempt to retrieve the SIP Server header.\n")
	fmt.Fprintf(os.Stderr, "'war': Wardialler that bruteforces extensions using various SIP methods.\n")
	fmt.Fprintf(os.Stderr, "'crack': Bruteforcer to crack SIP passwords for an extension.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Optional arguments:\n")
	w := new(tabwriter.Writer)
	w.Init(os.Stderr, 0, 8, 2, '\t', 0)
	fmt.Fprintf(w, "\t--timeout <sec>\tTimeout (in seconds) for each request. [DEFAULT: 3]\n")
	fmt.Fprintf(w, "\t--ext-start <num>\tFirst extension to try (wardialling). [DEFAULT: 0]\n")
	fmt.Fprintf(w, "\t--ext-end <num>\tLast extension to try (wardialling). [DEFAULT: 2000]\n")
	w.Flush()
	fmt.Fprintf(os.Stderr, "\n\nExample: %s map udp 192.168.0.20\n", os.Args[0])
}


func map_udp(targets []string, port int, timeout int) {
	res_targets := []string{}
	results := []string{}
	for _,target := range targets {
		result,err := siplib.MapUDP(target, port, timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not map %s:%d (%s)\n", target, port, err.Error())
		} else {
			res_targets = append(res_targets, target)
			results = append(results, result)
		}
	}
	fmt.Println("")
	if len(res_targets) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintf(w, "Target\tPort\tServer Header\n")
		fmt.Fprintf(w, "\t\t\t\n")
		for index,_ := range res_targets {
			fmt.Fprintf(w, "%s\t%d\t%s\n", res_targets[index], port, results[index])
		}
		w.Flush()
	} else {
		fmt.Println("No results found.")
	}
}

func war_udp(targets []string, port int, timeout int, min_ext int, max_ext int) {
	res_targets := []string{}
	results := []map[int]string{}
	extensions := []int{}
	for i := min_ext; i <= max_ext; i++ {
		extensions = append(extensions, i)
	}
	for _,target := range targets {
		result,err := siplib.WarInviteUDP(target, port, timeout, extensions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not wardial %s:%d (%s)\n", target, port, err.Error())
		}
		if len(result) > 0 {
			res_targets = append(res_targets, target)
			results = append(results, result)
		}
	}
	if len(res_targets) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintf(w, "Target\tPort\tExtension\tResult\n")
		fmt.Fprintf(w, "\t\t\t\t\n")
		for index,_ := range res_targets {
			for extension,value := range results[index] {
				fmt.Fprintf(w, "%s\t%d\t%d\t%s\n", res_targets[index], port, extension, value)
			}
		}
		w.Flush()
	} else {
		fmt.Println("No results found.")
	}
}
