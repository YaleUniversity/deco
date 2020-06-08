#!/bin/bash
rm -rf bin/*
go build -o bin/deco-native github.com/YaleUniversity/deco
chmod 755 bin/deco-native
export VERSION=`./bin/deco-native version -s`
for GOOS in darwin linux; do
  for GOARCH in 386 amd64; do
    echo "Building deco-v${VERSION}-${GOOS}-${GOARCH}"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o bin/deco-v${VERSION}-${GOOS}-${GOARCH} github.com/YaleUniversity/deco
  done
done
