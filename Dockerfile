# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

COPY . .

# Copy the local main.go file to the container's working directory
COPY main.go .

# Build the Go application
RUN go build -o main .

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
