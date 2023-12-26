# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required third-party dependencies
RUN go get -d -v ./...

# Build the Go application
RUN go install -v ./...

# Expose port 8080 to the outside world
EXPOSE 8080


# Run the application when the container starts
CMD ["app"]
