#!/bin/sh

goose status
echo "goose up started"
goose up
echo "goose up ended"
goose status

swag init -d cmd/api/,internal/app/api/controllers/system/,internal/app/api/controllers/switches/,internal/core/switches/models/,internal/core/common/
CompileDaemon --exclude-dir="docs" --build="./bin/build.sh" --command="./main" --color
