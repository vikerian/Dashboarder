# services/web-server/Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy templates and static files
COPY ./services/web-server/templates ./services/web-server/templates
COPY ./services/web-server/static ./services/web-server/static

RUN go build -o web ./services/web-server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/web .
COPY --from=builder /app/services/web-server/templates ./services/web-server/templates
COPY --from=builder /app/services/web-server/static ./services/web-server/static

EXPOSE 8090

CMD ["./web"]
