FROM golang:1.15

WORKDIR /go/src/wiki

RUN go get github.com/pilu/fresh
COPY cmd/wiki/ .

CMD ["fresh"]
