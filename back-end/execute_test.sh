#!/bin/sh

readonly FILE_NAME="docker-compose.test.yml"

TARGET_DIR="."
case "$1" in
  "controller")
    TARGET_DIR="./controllers"
    ;;
  "service")
    TARGET_DIR="./services"
    ;;
  "repository")
    TARGET_DIR="./repositories"
    ;;
esac

cd `dirname $0`
docker-compose -f $FILE_NAME up -d
docker-compose -f $FILE_NAME exec -T test-back sh -c "go test -cover ${TARGET_DIR}/..."

docker-compose -f $FILE_NAME down