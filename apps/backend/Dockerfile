FROM golang:1.25.3-alpine AS builder 

WORKDIR /app

COPY go.mod go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

FROM alpine:3.22.2 AS runner

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/config.yml ./config.yml
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./main"]
