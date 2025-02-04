# Build binary file
FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
COPY . .

RUN go build -o main main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]