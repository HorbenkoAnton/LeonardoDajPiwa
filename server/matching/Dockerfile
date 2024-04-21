FROM golang:latest

WORKDIR /data/app/matching

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the app
COPY backend/ ./
RUN go build -o main .

CMD ["./main"]