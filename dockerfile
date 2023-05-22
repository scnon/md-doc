FROM golang:1.20 AS builder

WORKDIR /app
COPY . .

RUN make build

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/ /app/bin/

ENTRYPOINT ["/app/bin/md-doc", "server"]