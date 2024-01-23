#!/bin/bash

killall printer-bot
CGO_ENABLED=0 go build -o httpd/printer-bot

httpd/printer-bot &
