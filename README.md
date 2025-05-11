# Gutils

_Pronounced like "goo-tills" - this is the repo for both the "utils" part of gutils as well as the build system for a MVP NodeMixaholic/Linux system._

A collection of basic Unix-like utilities written in Go, inspired by BusyBox and GNU CoreUtils. Each utility is a separate binary for simplicity and modularity.

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

***... and many, many more!***

## Project Structure

```
gutils/
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

## Using dosu (Minimal su/sudo Clone)

`dosu` allows you to run commands as root after password authentication using a SHA-256 hash stored in `/etc/dosu_passwd`.

### 1. Generate a Password Hash

You can generate a SHA-256 hash of your password using the following Go code or a command:

```
echo -n 'yourpassword' | sha256sum | awk '{print $1}'
```

### 2. Create /etc/dosu_passwd

Create the file `/etc/dosu_passwd` as root and paste the hash (no spaces or newlines):

```
sudo sh -c 'echo "<your_sha256_hash>" > /etc/dosu_passwd'
sudo chmod 600 /etc/dosu_passwd
```

### 3. Build dosu

```
make dosu
```

### 4. Set dosu as setuid root

This is required for dosu to escalate privileges:

```
sudo chown root:root bin/dosu
sudo chmod u+s bin/dosu
```

### 5. Run dosu

```
./bin/dosu whoami
```

You will be prompted for the password. If correct, the command will run as root.

## Dependencies

- Go 1.22 or newer
- The `golang.org/x/term` package (for dosu password prompt)

### Install Go

Follow instructions at https://go.dev/doc/install or use your package manager:

```
sudo dnf install golang  # Fedora
sudo apt install golang  # Debian/Ubuntu
```

### Install Go Module Dependencies

In the project root, run:

```
go mod tidy
```

This will fetch `golang.org/x/term` and any other required modules.

## License
MIT 

## Example: Using tar with gzip compression

Create a compressed archive:
```
./bin/tar -c -f archive.tar.gz file1 dir2
```

Extract a compressed archive:
```
./bin/tar -x -f archive.tar.gz
``` 
