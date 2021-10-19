#!/bin/bash
env GOOS=linux GOARCH=arm64 go build -o "bme280-device-service" -ldflags "-w -s"

if [ ! -d bin ];then
    mkdir bin
fi
mv -f bme280-device-service bin