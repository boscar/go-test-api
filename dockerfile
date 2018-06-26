FROM golang:latest

RUN mkdir /go-test-api
ADD . /go-test-api/
WORKDIR /go-test-api

RUN go get github.com/gorilla/handlers
RUN go get github.com/tkanos/gonfig
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/gorilla/context
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get github.com/boscar/go-test-api

RUN go build 

RUN go build -o main .
CMD ["/app/main"]