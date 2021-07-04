FROM golang:1.16-alpine

ENV APP_ROOT /app
WORKDIR ${APP_ROOT}
COPY . ${APP_ROOT}
RUN go mod download
CMD ["go", "run", "main.go"]
