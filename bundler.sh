#!/bin/bash
VERSION=v0.8.0

echo "Bundling xim $VERSION"

( GOOS=darwin GOARCH=amd64 go build -o bin/macos-darwin/amd64/xim && tar -czf bin/macos-darwin/amd64/xim-$VERSION-amd64.tar.gz -C bin/macos-darwin/amd64/ $(ls  bin/macos-darwin/amd64/))
( GOOS=darwin GOARCH=arm64 go build -o bin/macos-darwin/arm64/xim && tar -czf bin/macos-darwin/arm64/xim-$VERSION-arm64.tar.gz -C bin/macos-darwin/arm64/ $(ls  bin/macos-darwin/arm64/))
( GOOS=linux GOARCH=386 go build -o bin/linux/386/xim && tar -czf bin/linux/386/xim-$VERSION-386.tar.gz -C bin/linux/386/ $(ls  bin/linux/386/))
( GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/xim && tar -czf bin/linux/amd64/xim-$VERSION-amd64.tar.gz -C bin/linux/amd64/ $(ls  bin/linux/amd64/))
( GOOS=linux GOARCH=arm go build -o bin/linux/arm/xim && tar -czf bin/linux/arm/xim-$VERSION-arm.tar.gz -C bin/linux/arm/ $(ls  bin/linux/arm/))
( GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/xim && tar -czf bin/linux/arm64/xim-$VERSION-arm64.tar.gz -C bin/linux/arm64/ $(ls  bin/linux/arm64/))
( GOOS=freebsd GOARCH=arm go build -o bin/freebsd/arm/xim && tar -czf bin/freebsd/arm/xim-$VERSION-arm.tar.gz -C bin/freebsd/arm/ $(ls  bin/freebsd/arm/))
( GOOS=freebsd GOARCH=amd64 go build -o bin/freebsd/amd64/xim && tar -czf bin/freebsd/amd64/xim-$VERSION-amd64.tar.gz -C bin/freebsd/amd64/ $(ls  bin/freebsd/amd64/))
( GOOS=freebsd GOARCH=386 go build -o bin/freebsd/386/xim && tar -czf bin/freebsd/386/xim-$VERSION-386.tar.gz -C bin/freebsd/386/ $(ls  bin/freebsd/386/))
( GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/xim.exe && zip -r -j bin/windows/amd64/xim-$VERSION-x64.zip bin/windows/amd64/xim.exe)
( GOOS=windows GOARCH=arm64 go build -o bin/windows/arm64/xim.exe && zip -r -j bin/windows/arm64/xim-$VERSION-x64.zip bin/windows/arm64/xim.exe)
( GOOS=windows GOARCH=386 go build -o bin/windows/386/xim.exe && zip -r -j bin/windows/386/xim-$VERSION-x86.zip bin/windows/386/xim.exe)


echo "Fin"