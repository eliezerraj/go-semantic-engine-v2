# docker build -t go-semantic-engine-v2 .
# docker run -dit --name go-semantic-engine-v2 -p 7004:7004 go-semantic-engine-v2

FROM golang:1.25 AS builder

RUN apt-get update && apt-get install bash && apt-get install -y --no-install-recommends ca-certificates

WORKDIR /app
COPY . .
RUN go mod tidy

WORKDIR /app
RUN go build -o go-semantic-engine-v2 -ldflags '-linkmode external -w -extldflags "-static"'

FROM alpine

WORKDIR /app
COPY --from=builder /app/go-semantic-engine-v2 .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/go-semantic-engine-v2"]