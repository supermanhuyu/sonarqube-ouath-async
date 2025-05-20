# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o sonarqube-ouath-async main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/sonarqube-ouath-async ./sonarqube-ouath-async
COPY --from=builder /app/config/ ./config/
COPY --from=builder /app/README.md ./README.md
EXPOSE 8080
ENTRYPOINT ["./sonarqube-ouath-async"]
