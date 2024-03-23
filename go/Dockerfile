FROM golang:1.22 AS builder

LABEL maintainer="dmaloney@medicalone.com.au"
LABEL version="1.0"

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum main.go ./
RUN go mod download && \
    go mod verify

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/smdr . && \
    chmod +x /app/smdr


FROM alpine:3

WORKDIR /app

COPY --from=builder /app/smdr /app/smdr

# Run
CMD ["/app/smdr"]