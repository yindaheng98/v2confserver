FROM golang:1.16-alpine AS builder
WORKDIR /usr/src/app
COPY . .
RUN go mod download && go mod verify && \
    go build -v -o /v2confserver cmd/v2confserver

FROM alpine

COPY --from=builder /v2confserver /usr/local/bin/v2confserver

ENTRYPOINT [ "v2confserver" ]