FROM golang:alpine AS builder

WORKDIR /build
COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine

WORKDIR /build
COPY --from=builder /build/main /build/main

CMD ["./main"]
