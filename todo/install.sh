#!/bin/bash
set -e

echo "Building todo..."
go build -o todo

echo "Installing to /usr/local/bin..."
sudo cp todo /usr/local/bin/

echo "Done! You can now use 'todo' from anywhere."
