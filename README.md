# SIPfurious

A SIP scanner written in Golang, inspired by SIPvicious. It can enumerate SIP servers over UDP, TCP and TLS. This tool is currently in beta, and is not yet ready for public consumption. Use it at your own risk!

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

## TODO

Currently implemented:

- 'map' over UDP
- 'war' over UDP
- 'crack' over UDP

Coming soon:

- 'map' over TCP
- 'war' over TCP
- 'crack' over TCP
- 'map' over TLS
- 'war' over TLS
- 'crack' over TLS
- Implement some kind of warning when a lot of 403 forbidden is detected (need higher throttle).
- Much, much more testing!


