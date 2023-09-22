#!/bin/bash

# Compile for Linux
env GOOS=linux GOARCH=amd64 go build -o bin/kubelogger-linux

# Compile for macOS
env GOOS=darwin GOARCH=amd64 go build -o bin/kubelogger-mac-intel

# compile for macOS M1
env GOOS=darwin GOARCH=arm64 go build -o bin/kubelogger-mac-m1

# Compile for Windows
env GOOS=windows GOARCH=amd64 go build -o bin/kubelogger-windows.exe
