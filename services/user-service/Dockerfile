FROM golang:1.23.2 as builder
WORKDIR /go/src/user-service
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/worker/main.go

FROM alpine:3.20.3 as release
WORKDIR /app
COPY --from=builder /go/src/user-service/main .
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/main

EXPOSE 3001
ENTRYPOINT ["/usr/bin/dumb-init", "--"]