FROM golang:latest

WORKDIR /data/app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the app
COPY . .
RUN go build -o main .

CMD ["./main"]