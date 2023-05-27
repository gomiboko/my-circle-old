#!/bin/sh

readonly FILE_NAME="docker-compose.test.yml"

cd `dirname $0`
docker compose -f $FILE_NAME up
docker compose -f $FILE_NAME down
