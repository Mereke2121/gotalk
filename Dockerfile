FROM golang:latest AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o gotalk ./cmd/main.go

CMD ["./gotalk"]