FROM golang:1.23 

WORKDIR /app

# Copy the Go module files first
COPY go.mod go.sum ./

# Initialize the module and download dependencies
RUN go mod tidy
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
WORKDIR /app/cmd/job
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Change back to root directory
WORKDIR /app

EXPOSE 8080

CMD ["./main"]