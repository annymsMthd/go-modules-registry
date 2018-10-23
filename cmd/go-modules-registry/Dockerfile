FROM golang:1.11.0

WORKDIR /app

COPY . .

RUN go build -o go-modules-registry ./cmd/go-modules-registry

FROM ubuntu:bionic

RUN mkdir /app
WORKDIR /app

COPY --from=0 /app/go-modules-registry /app/

RUN mkdir /tmp/storage

CMD ["/app/go-modules-registry"]
