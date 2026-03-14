FROM golang:1.24.2 as builder

WORKDIR /app

# Copy the GO module files
COPY go.mod .
COPY go.sum .

# Download the Go module
RUN go mod download

COPY main.go .

RUN go build

RUN mkdir cache
EXPOSE 8000

CMD ["./nx-cache-server"]
