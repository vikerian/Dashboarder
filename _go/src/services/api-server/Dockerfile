# services/api-server/Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./services/api-server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
