FROM golang:1.23 AS builder

# COPY go.mod go.sum ./

# RUN go mod tidy

COPY . .

RUN go build -o /main .

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /ppml

COPY . .

COPY --from=builder /main .

CMD ["./main"]