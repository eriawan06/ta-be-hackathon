FROM golang:1.17

#ENV GO111MODULE=on

#ENV GOPATH /go
#ENV PKG_NAME=be-sagara-hackathon
#ENV PKG_PATH=$GOPATH/src/$PKG_NAME

WORKDIR /go/src/be-sagara-hackathon/
#COPY . /go/src/be-sagara-hackathon/
COPY . /go/src/be-sagara-hackathon/

RUN go mod vendor
RUN go build .

CMD ["sh", "-c", "main"]
