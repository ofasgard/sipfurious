# SIPfurious

A SIP scanner written in Golang, inspired by SIPvicious. It can enumerate SIP servers over UDP, TCP and TLS. This tool is currently in development, and has not been widely tested across many different platforms. It is likely to be riddled with bugs. 

![example usage](https://user-images.githubusercontent.com/19550999/76960818-23da6880-6914-11ea-89d2-b7f2347e3e5d.png)

Currently, this tool has been tested against:

- Asterisk PBX 1.6.2.11
- Asterisk PBX 16.6.2
- Asterisk PBX 13.10.0 (throttle increased to 500ms)

## Building

SIPfurious has no dependencies besides Go itself. To build the program, just do:

```shell
$ git clone https://github.com/ofasgard/sipfurious
$ cd sipfurious
$ ./build.sh
$ bin/sipfurious --help
```

## Usage

```
Usage: bin/sipfurious <map|war|crack> <udp|tcp|tls> <target> [port]

'map': Scanner that uses OPTIONS to attempt to retrieve the SIP Server header.
'war': Wardialler that bruteforces extensions using the INVITE method.
'crack': Bruteforcer to crack SIP passwords for an extension.

Optional arguments:
	--timeout <sec>		Timeout (in seconds) for each request. [DEFAULT: 10]
	--throttle <ms>		Delay (in milliseconds) between each request when wardialing or password cracking. [DEFAULT: 100]
	--wordlist <file>	Specify a wordlist file to use for wardialing or password cracking.
	--user <user>		Specify a username to use; required for password cracking.


Example: bin/sipfurious map udp 192.168.0.20
```

## TODO

Currently implemented:

- 'map' over UDP
- 'war' over UDP
- 'crack' over UDP
- 'map' over TCP (*untested*)
- 'war' over TCP (*untested*)

Coming soon:

- 'crack' over TCP
- 'map' over TLS
- 'war' over TLS
- 'crack' over TLS
- Implement some kind of warning when a lot of 403 forbidden is detected (need higher throttle).
- Much, much more testing!


