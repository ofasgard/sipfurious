# SIPfurious

A SIP scanner written in Golang, inspired by SIPvicious. It can enumerate SIP servers over UDP, TCP and TLS. This tool is currently in beta, and is not yet ready for public consumption. Use it at your own risk!

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

Coming soon:

- 'crack' over UDP
- 'map' over TCP
- 'war' over TCP
- 'crack' over TCP
- 'map' over TLS
- 'war' over TLS
- 'crack' over TLS
- Much, much more testing!


