#!/bin/sh

swag init -d cmd/api/,internal/app/api/controllers/system/,internal/app/api/controllers/switches/,internal/core/switches/models/,internal/core/common/
go build -C ./cmd/api/ -v -o ../../main -ldflags "-X main.compileDate=`date +%Y-%m-%dT%T.%9N%:z`"
