FROM golang:1.24

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o follow-api ./cmd/follow-api

EXPOSE 8083
CMD ["./follow-api"]