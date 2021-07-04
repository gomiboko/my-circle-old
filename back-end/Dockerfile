FROM golang:1.16-alpine

# Goのホットリロードツール(Air)インストール
WORKDIR /air
RUN apk update && apk add --no-cache curl
RUN set -o pipefail && curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
RUN ./bin/air init

ENV APP_ROOT /app
WORKDIR ${APP_ROOT}
COPY . ${APP_ROOT}
RUN go mod download
CMD ["/air/bin/air"]
