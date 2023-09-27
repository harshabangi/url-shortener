FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY resources/ /app/config/
RUN go build -o url-shortener
EXPOSE 8080
CMD ["./url-shortener"]