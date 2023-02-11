FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .

WORKDIR /app/cmd/place
RUN go mod download
RUN go build -o /app/main

FROM alpine:3.14

WORKDIR /app
COPY --from=builder /app/main /app/main
RUN chmod +x /app/main

CMD ["/app/main"]