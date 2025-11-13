FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o orchestrator .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/orchestrator .
EXPOSE 8081
CMD ["./orchestrator"]
