FROM golang:1.16.5

# Create a working directory
WORKDIR /usr/src/app

# Download dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# Copy source code
COPY [".", "."]

# Build app binary
RUN go build -o ./app-name ./cmd/app-name

ENTRYPOINT [ "./app-name" ]
