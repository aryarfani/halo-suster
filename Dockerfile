################################################### STAGE 1
FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app/

COPY . .

RUN go mod tidy

RUN go build -o main ./main.go

################################################### STAGE 2
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./
COPY --from=builder /app/.env ./

EXPOSE 8080

ENTRYPOINT ["./main"]