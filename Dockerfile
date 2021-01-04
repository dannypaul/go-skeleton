FROM golang:1.15.6

# Create a working directory
WORKDIR /usr/src/app

# Copy source code
COPY [".", "."]

# Build app binary
RUN go build -o ./app-name ./cmd/app-name

ENTRYPOINT [ "./app-name" ]
