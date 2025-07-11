package main

import "fmt"

func cmdHelp() {
	fmt.Println(`bruv: a simple git-like tool

Usage: bruv <command> [arguments]

Commands:
  init [<directory>]
    Create an empty bruv repository in the specified directory, or the current one if not provided.
    Example: bruv init my-new-repo

  hash-object -w <file>
    Compute object ID and optionally creates a blob from a file.
    Example: bruv hash-object -w README.md

  cat-file [-t | -p] <object>
    Provide content, type, or size information for repository objects.
    -t: show object type
    -p: pretty-print object content
    Example: bruv cat-file -p a1b2c3d4

  add <file>...
    Add file contents to the index.
    Example: bruv add main.go utils.go

  commit -m <message>
    Record changes to the repository.
    Example: bruv commit -m "Initial commit"

  serve
    Start a bruv server, listening for connections on port 9418.
    Example: bruv serve

  clone <url> <directory>
    Clone a repository from a URL into a new directory.
    Example: bruv clone localhost:9418/my-repo my-cloned-repo

  help
    Show this help message.
`)
} 