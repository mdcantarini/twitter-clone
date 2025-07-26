FROM golang:1.24

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o user-api ./cmd/user-api

EXPOSE 8081
CMD ["./user-api"]