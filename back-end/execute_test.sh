#!/bin/sh

readonly FILE_NAME="docker-compose.test.yml"

cd `dirname $0`
docker-compose -f $FILE_NAME up -d
docker-compose -f $FILE_NAME exec -T test-back sh -c "go test -cover ./..."

# ローカルでの実行時は、最後にコンテナを破棄する
if [ "$1" == "local" ]; then
  docker-compose -f $FILE_NAME down
fi