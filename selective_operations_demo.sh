#!/bin/bash

echo "=== Bruv Selective Operations Demo ==="
echo ""

echo "This script demonstrates the new selective operations features in bruv:"
echo ""

echo "1. Selective Clone"
echo "   Usage: bruv clone <url> <directory> --select <path>..."
echo "   Example: bruv clone localhost:9418/my-repo my-cloned-repo --select src/ docs/"
echo "   Implementation: cmd/bruv/main.go (cmdClone function)"
echo ""

echo "2. Selective Pull"
echo "   Usage: bruv pull <remote> [<branch>] --select <path>..."
echo "   Example: bruv pull localhost:9418/my-repo main --select src/ docs/"
echo "   Implementation: cmd/bruv/pull.go (cmdPull function)"
echo ""

echo "3. Selective Push"
echo "   Usage: bruv push <remote> [<branch>] --select <path>..."
echo "   Example: bruv push localhost:9418/my-repo main --select src/ docs/"
echo "   Implementation: cmd/bruv/push.go (cmdPush function)"
echo ""

echo "4. Selective Merge"
echo "   Usage: bruv merge <source-branch> <target-branch> --select <path>..."
echo "   Example: bruv merge feature-branch main --select src/ docs/"
echo "   Implementation: cmd/bruv/merge.go (cmdMerge function)"
echo ""

echo "These selective operations allow you to:"
echo "- Clone only specific files or folders from a repository"
echo "- Pull only specific files or folders from a remote repository"
echo "- Push only specific files or folders to a remote repository"
echo "- Merge only specific files or folders between branches"
echo ""

echo "Implementation details:"
echo "- All commands parse --select arguments to identify paths"
echo "- Server-side support handles selective requests"
echo "- Client-side creates selective packfiles with only specified objects"
echo "- Working directory is filtered to include only specified paths"
echo ""

echo "Note: The current implementation shows selective operations as simplified examples."
echo "In a full implementation, these operations would:"
echo "- Parse commit trees to identify objects related to specified paths"
echo "- Create packfiles containing only those objects and their dependencies"
echo "- Filter working directories to include only specified paths"
echo ""