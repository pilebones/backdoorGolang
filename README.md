# backdoorGolang
Backdoor with Golang (Cross-Plateform)

_/!\ Work in progress, not a stable release /!\_

##Main goal

A fork of my own project named : "pilebones/backdoorBash" (see: https://github.com/pilebones/backdoorBash) but instead of using Bash as programming language (Unix-like only) this new one will work on Windows too by using a Golang API (cross-plateform) developed from scratch (as much as possible).

## Requirements

- Golang SDK : Compiler and tools for the Go programming language from Google (see: https://golang.org/doc/install)

From Arch Linux :
```bash
(sudo) pacman -S community/go
```

From Debian :
```bash
(sudo) apt-get install golang-go
```

## Installation

```bash
cd $GOPATH
go get github.com/pilebones/backdoorGolang
./bin/backdoorGolang --help

```

## Usage

```bash
./bin/backdoorGolang --help
Usage of ./bin/backdoorGolang:
  -d, --debug         Enable mode debug
  -h, --host string   Set hostname to use (default "localhost")
  -l, --listen        Enable listen mode (server socket mode)
  -p, --port int      Set port number to use (default 9876)
  -v, --verbose       Enable mode verbose
  -V, --version       Display version number
```

## Server-mode

```bash
./bin/backdoorGolang -h localhost -p 1234
```

Notice : Server is multi-user capable (one server for X client)

## Client-mode

/!\ Not implemented yet, use netcat meanwhile !

```bash
netcat localhost 1234
```

