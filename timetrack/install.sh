#!/bin/bash
set -e

cd "$(dirname "$0")"
go build -o timetrack
sudo cp timetrack /usr/local/bin/
echo "Installed timetrack to /usr/local/bin/"
