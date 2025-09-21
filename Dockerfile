FROM golang:1.24

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build
