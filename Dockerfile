FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /review-assigner "./cmd"

FROM alpine:latest

WORKDIR /app
COPY --from=builder /review-assigner .

EXPOSE 8080

CMD ["./review-assigner"]