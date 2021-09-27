FROM golang:1.17.1-alpine as builder


ARG UID=1001
ARG GID=1001
# Set the Current Working Directory inside the container
WORKDIR /app


# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
 go build -a -installsuffix cgo -o ./cmd/main ./cmd


######### Start a new stage from scratch #######
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/cmd/main /app/bin/main

WORKDIR /app/bin

## Command to run the executable
CMD ["./main"]