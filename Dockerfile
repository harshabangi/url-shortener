# Use the official Golang image as the base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code into the container
COPY . .

# Build the Golang application
RUN go build -o url-shortener

# Expose the port that the application will listen on
EXPOSE 8080

# Command to run the application
CMD ["./url-shortener"]