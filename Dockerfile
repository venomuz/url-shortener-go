FROM golang:1.18rc1-alpine3.15
WORKDIR /go/src/github.com/url-shortener-go
COPY . .
#RUN go mod init
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main