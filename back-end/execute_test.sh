#!/bin/sh

readonly FILE_NAME="docker-compose.test.yml"

cd `dirname $0`
docker-compose -f $FILE_NAME up -d
docker-compose -f $FILE_NAME exec -T test-back sh -c "go test -cover ./..."

docker-compose -f $FILE_NAME down