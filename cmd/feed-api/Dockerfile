FROM golang:1.24

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o feed-api ./cmd/feed-api

EXPOSE 8084
CMD ["./feed-api"]