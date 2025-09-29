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

  clone <url> <directory> [--select <path>...]
    Clone a repository from a URL into a new directory.
    Use --select to clone only specific files or folders.
    Example: bruv clone localhost:9418/my-repo my-cloned-repo
    Example: bruv clone localhost:9418/my-repo my-cloned-repo --select src/ docs/

  push <remote> [<branch>] [--select <path>...]
    Push commits to a remote repository.
    Use --select to push only specific files or folders.
    Example: bruv push localhost:9418/my-repo main
    Example: bruv push localhost:9418/my-repo main --select src/ docs/

  pull <remote> [<branch>] [--select <path>...]
    Pull commits from a remote repository.
    Use --select to pull only specific files or folders.
    Example: bruv pull localhost:9418/my-repo main
    Example: bruv pull localhost:9418/my-repo main --select src/ docs/

  merge <source-branch> <target-branch> [--select <path>...]
    Submit a merge request from source branch to target branch.
    Use --select to merge only specific files or folders.
    Requires owner approval for merging to master/main.
    Example: bruv merge feature-branch main
    Example: bruv merge feature-branch main --select src/ docs/

  approve <merge-request-id>
    Approve and merge a pending merge request (owner only).
    Example: bruv approve 1234567890

  list-requests
    List all pending and approved merge requests.
    Example: bruv list-requests

  help
    Show this help message.

Notes:
  - .bruvignore files work like .gitignore files to exclude files from version control
  - Remote URLs use format: host:port/repository-path
`)
}