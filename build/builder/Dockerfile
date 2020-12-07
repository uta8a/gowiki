FROM golang:1.15

WORKDIR /go/src/wiki

COPY cmd/wiki/ .
RUN go get github.com/pilu/fresh
CMD ["fresh"]
EXPOSE 9000