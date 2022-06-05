# Dockerfile References: https://docs.docker.com/engine/reference/builder/
FROM golang:1.17.1 as builder
LABEL maintainer="TaiLeHuu"
RUN export PATH=$PATH:/go/bin

WORKDIR /app
COPY . .

RUN go mod tidy \
    && go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o main cmd/*


######## Start a new stage from scratch #######
FROM alpine:3.15.0
RUN apk update \
    && apk --no-cache add ca-certificates curl

WORKDIR /app
COPY --from=builder /app/main .

USER 1001
ENTRYPOINT ["./main"]
