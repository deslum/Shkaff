#!/bin/bash
export GOPATH=$PWD
export GOARCH=amd64
export GOOS=linux
export GOBIN=$(go env GOPATH)/bin
go get github.com/takama/daemon
go get github.com/streadway/amqp
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
go get github.com/syndtr/goleveldb/leveldb
go install src/shkaff/shkaff.go
