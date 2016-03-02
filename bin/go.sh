#!/bin/sh

# cd go src, Run GOOS=linux GOARCH=amd64 ./make.bash --no-clean; you can build a go package or binary for this architecture using GOOS=linux GOARCH=amd64 go build. Y

x=`date +%F_%T`
y=`git rev-parse HEAD`

OUT=st_linux_amd64
#GOOS=linux go build -ldflags "-X main.date=$x -X main.rev=$y" -o $OUT
GOOS=linux  GOARCH=amd64 go build -o $OUT

scp $OUT sharetrace@beijing6.appdao.com:~/sharetrace
