# Use the Go base image
FROM golang:1.23.2

# Set the working directory
WORKDIR /app

# Copy the Go files
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Expose the port for the Go service
EXPOSE 4000

# Run the Go app
CMD ["./main"]
