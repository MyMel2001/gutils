# Sneed Coreutils

A collection of basic Unix-like utilities written in Go, inspired by BusyBox. Each utility is a separate binary for simplicity and modularity.

## Utilities Included
- highway (minimal shell)
- ls
- cat
- echo
- cp
- mv
- rm
- printf
- mkdir
- grep
- head
- alias
- tail
- wc

## Project Structure

```
sneed-coreutils/
├── cmd/
│   ├── highway/
│   │   └── main.go
│   ├── ls/
│   │   └── main.go
│   └── ... (other utilities)
├── go.mod
├── .gitignore
└── README.md
```

## Building a Utility

To build a utility, run:

```
go build -o bin/ls ./cmd/ls
```

Replace `ls` with the utility you want to build. All binaries will be placed in the `bin/` directory.

## Running

After building, run a utility like this:

```
./bin/ls
```

## License
MIT 