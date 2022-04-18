# Build stage
FROM golang:1.17.8-alpine3.15 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o main main.go

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate.linux-amd64 ./migrate
# COPY app.env .
# COPY start.sh .
# COPY wait-for.sh .
# COPY db/migration ./migration
    
EXPOSE 8080
CMD [ "/app/main" ]
# ENTRYPOINT [ "/app/start.sh" ]