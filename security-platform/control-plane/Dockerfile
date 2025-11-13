FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY api/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o control-plane .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/control-plane .
EXPOSE 8080
CMD ["./control-plane"]
