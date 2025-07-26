FROM golang:1.24

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o tweet-api ./cmd/tweet-api

EXPOSE 8082
CMD ["./tweet-api"]