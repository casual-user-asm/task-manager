FROM golang:1.23 AS builder

WORKDIR /usr/src/app

# Copy source files
COPY . .

# Download dependencies
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /main

CMD ["/main"]