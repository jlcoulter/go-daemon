#!/usr/bin/env bash
# setup.sh — rename module and clean up template markers
# Usage: ./setup.sh mydaemon github.com/you/mydaemon

set -euo pipefail

if [ $# -lt 2 ]; then
  echo "Usage: $0 <daemon-name> <module-path>"
  echo "Example: $0 sshbastion github.com/jlcoulter/sshbastion"
  exit 1
fi

NAME="$1"
MODULE="$2"
OLD_MODULE="github.com/jlcoulter/go-daemon-template"
OLD_NAME="go-daemon-template"

# Replace module path in all Go files
find . -name "*.go" -exec sed -i "s|${OLD_MODULE}|${MODULE}|g" {} +

# Replace in go.mod
sed -i "s|${OLD_MODULE}|${MODULE}|g" go.mod

# Replace in Makefile
sed -i "s|${OLD_NAME}|${NAME}|g" Makefile

# Replace in Dockerfile
sed -i "s|${OLD_NAME}|${NAME}|g" Dockerfile

# Replace in README
sed -i "s|${OLD_NAME}|${NAME}|g" README.md

# Replace in .goreleaser.yml
sed -i "s|${OLD_NAME}|${NAME}|g" .goreleaser.yml

# Replace in CI
sed -i "s|${OLD_NAME}|${NAME}|g" .github/workflows/ci.yml

# Replace in systemd unit
sed -i "s|${OLD_NAME}|${NAME}|g" contrib/${OLD_NAME}.service
mv contrib/${OLD_NAME}.service contrib/${NAME}.service

# Remove setup script
rm -- "$0"

echo "Template configured: name=${NAME}, module=${MODULE}"
echo "Next steps:"
echo "  1. Run: go mod tidy"
echo "  2. Run: make test"
echo "  3. Replace internal/daemon/daemon.go with your protocol logic"