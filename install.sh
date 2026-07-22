#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tmuxer .


sudo mv tmuxer /usr/local/bin
