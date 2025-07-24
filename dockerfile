# ---- Build Stage ----
FROM golang:1.23.3-alpine AS builder 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# ---- Run Stage ----
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app ./app
EXPOSE 8320
CMD ["./app"]