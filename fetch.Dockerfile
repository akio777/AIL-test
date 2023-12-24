# Use the official Golang image as the base image
FROM golang:1.21.4

# Set the working directory inside the container
WORKDIR /build

# Copy the Go Modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code to the container
COPY . .

# Build the Go application
RUN go build -o fetch ./cmd/fetch/main.go

# Command to run the executable
CMD ["./fetch"]
