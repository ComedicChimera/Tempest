FROM golang:latest

RUN go get github.com/gorilla/mux github.com/jinzhu/gorm github.com/dgrijalva/jwt-go golang.org/x/crypto/bcrypt github.com/lib/pq

ADD . /go/src/github.com/tempest-server

WORKDIR /go/src/github.com/tempest-server

RUN go build -o main main.go

EXPOSE 5000
