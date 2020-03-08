#!/bin/bash

mkdir file2src_builds -p

cp README.md file2src_builds
env GOOS=linux GOARCH=amd64 go build -o file2src_builds/linux_64bit/file2src main.go
env GOOS=linux GOARCH=arm go build -o file2src_builds/linux_arm/file2src main.go
env GOOS=darwin GOARCH=amd64 go build -o file2src_builds/mac_64bit/file2src main.go
env GOOS=windows GOARCH=386 go build -o file2src_builds/win/file2src.exe main.go

rm file2src_builds.zip
zip -r file2src_builds.zip file2src_builds 