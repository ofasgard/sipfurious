package main

import "siplib"
import "fmt"
import "flag"
import "os"
import "strconv"
import "text/tabwriter"
import "io/ioutil"

func main() {
	//parse flags
	flag.Usage = usage
	timeout_ptr := flag.Int("timeout", 10, "")
	wordlist_ptr := flag.String("wordlist", "", "")
	throttle_ptr := flag.Int("throttle", 100, "")
	username_ptr := flag.String("user", "", "")
	flag.Parse()
	timeout := *timeout_ptr
	wordlist_path := *wordlist_ptr
	throttle := *throttle_ptr
	username := *username_ptr
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
	//validate wordlist
	wordlist := []string{}
	if wordlist_path != "" {
		fd,err := os.Open(wordlist_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open '%s' for reading as a wordlist file. Quitting!\n", wordlist_path)
			return
		}
		buf,err := ioutil.ReadAll(fd)
		fd.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open '%s' for reading as a wordlist file. Quitting!\n", wordlist_path)
			return
		}
		parts := splitlines(string(buf))
		for _,part := range parts {
			if part != "" { wordlist = append(wordlist, part) }
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
					extensions := default_extensions()
					if len(wordlist) > 0 {
						extensions = wordlist
					}
					war_udp(targets, port, timeout, throttle, extensions)
					return
				case "crack":
					if username == "" {
						fmt.Fprintf(os.Stderr, "Password cracking requires you provide a username/extension with the --user parameter. Quitting!\n")
						return
					}
					passwords := default_passwords()
					if len(wordlist) > 0 {
						passwords = wordlist
					}
					crack_udp(targets, port, timeout, throttle, username, passwords)
					return
				default:
					usage()
					return
			}
		case "tcp":
			switch method{
				case "map":
					map_tcp(targets, port, timeout)
					return
				case "war":
					extensions := default_extensions()
					if len(wordlist) > 0 {
						extensions = wordlist
					}
					war_tcp(targets, port, timeout, throttle, extensions)
					return
				case "crack":
					if username == "" {
						fmt.Fprintf(os.Stderr, "Password cracking requires you provide a username/extension with the --user parameter. Quitting!\n")
						return
					}
					passwords := default_passwords()
					if len(wordlist) > 0 {
						passwords = wordlist
					}
					crack_tcp(targets, port, timeout, throttle, username, passwords)
					return
				default:
					usage()
					return
			}
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
	fmt.Fprintf(os.Stderr, "'war': Wardialler that bruteforces extensions using the INVITE method.\n")
	fmt.Fprintf(os.Stderr, "'crack': Bruteforcer to crack SIP passwords for an extension.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Optional arguments:\n")
	w := new(tabwriter.Writer)
	w.Init(os.Stderr, 0, 8, 2, '\t', 0)
	fmt.Fprintf(w, "\t--timeout <sec>\tTimeout (in seconds) for each request. [DEFAULT: 10]\n")
	fmt.Fprintf(w, "\t--throttle <ms>\tDelay (in milliseconds) between each request when wardialing or password cracking. [DEFAULT: 100]\n")
	fmt.Fprintf(w, "\t--wordlist <file>\tSpecify a wordlist file to use for wardialing or password cracking.\n")
	fmt.Fprintf(w, "\t--user <user>\tSpecify a username to use; required for password cracking.\n")
	w.Flush()
	fmt.Fprintf(os.Stderr, "\n\nExample: %s map udp 192.168.0.20\n", os.Args[0])
}


func map_udp(targets []string, port int, timeout int) {
	res_targets := []string{}
	results := []string{}
	for _,target := range targets {
		fmt.Printf("Trying %s:%d...\n", target, port)
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


func map_tcp(targets []string, port int, timeout int) {
	res_targets := []string{}
	results := []string{}
	for _,target := range targets {
		fmt.Printf("Trying %s:%d...\n", target, port)
		result,err := siplib.MapTCP(target, port, timeout)
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

func war_udp(targets []string, port int, timeout int, throttle int, extensions []string) {
	res_targets := []string{}
	results := []map[string]string{}
	for _,target := range targets {
		fmt.Printf("Trying %s:%d...\n", target, port)
		result,err := siplib.WarInviteUDP(target, port, timeout, throttle, extensions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not wardial %s:%d (%s)\n", target, port, err.Error())
		} else if len(result) > 0 {
			res_targets = append(res_targets, target)
			results = append(results, result)
		}
	}
	fmt.Println("")
	if len(res_targets) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintf(w, "Target\tPort\tExtension\tResult\n")
		fmt.Fprintf(w, "\t\t\t\t\n")
		for index,_ := range res_targets {
			for extension,value := range results[index] {
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", res_targets[index], port, extension, value)
			}
		}
		w.Flush()
	} else {
		fmt.Println("No results found.")
	}
}

func war_tcp(targets []string, port int, timeout int, throttle int, extensions []string) {
	res_targets := []string{}
	results := []map[string]string{}
	for _,target := range targets {
		fmt.Printf("Trying %s:%d...\n", target, port)
		result,err := siplib.WarInviteTCP(target, port, timeout, throttle, extensions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not wardial %s:%d (%s)\n", target, port, err.Error())
		} else if len(result) > 0 {
			res_targets = append(res_targets, target)
			results = append(results, result)
		}
	}
	fmt.Println("")
	if len(res_targets) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintf(w, "Target\tPort\tExtension\tResult\n")
		fmt.Fprintf(w, "\t\t\t\t\n")
		for index,_ := range res_targets {
			for extension,value := range results[index] {
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", res_targets[index], port, extension, value)
			}
		}
		w.Flush()
	} else {
		fmt.Println("No results found.")
	}
}


func crack_udp(targets []string, port int, timeout int, throttle int, extension string, passwords []string) {
	res_targets := []string{}
	results := []string{}
	for _,target := range targets {
		fmt.Printf("Trying %s:%d...\n", target, port)
		result,err := siplib.BruteforceRegisterUDP(target, port, timeout, throttle, extension, passwords)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		} else if len(result) > 0 {
			res_targets = append(res_targets, target)
			results = append(results, result)
		}
	}
	fmt.Println("")
	if len(res_targets) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintf(w, "Target\tPort\tUser\tPassword\n")
		fmt.Fprintf(w, "\t\t\t\n")
		for index,_ := range res_targets {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", res_targets[index], port, extension, results[index])
		}
		w.Flush()
	} else {
		fmt.Println("No results found.")
	}
}

func crack_tcp(targets []string, port int, timeout int, throttle int, extension string, passwords []string) {
	res_targets := []string{}
	results := []string{}
	for _,target := range targets {
		fmt.Printf("Trying %s:%d...\n", target, port)
		result,err := siplib.BruteforceRegisterTCP(target, port, timeout, throttle, extension, passwords)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		} else if len(result) > 0 {
			res_targets = append(res_targets, target)
			results = append(results, result)
		}
	}
	fmt.Println("")
	if len(res_targets) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintf(w, "Target\tPort\tUser\tPassword\n")
		fmt.Fprintf(w, "\t\t\t\n")
		for index,_ := range res_targets {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", res_targets[index], port, extension, results[index])
		}
		w.Flush()
	} else {
		fmt.Println("No results found.")
	}
}

