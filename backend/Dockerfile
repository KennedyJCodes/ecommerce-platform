FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /build/app .

EXPOSE 8080

CMD ["./app"]
