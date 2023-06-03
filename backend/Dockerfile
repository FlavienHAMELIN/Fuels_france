FROM golang:1.20.4

WORKDIR /

COPY * ./

RUN go build

CMD ["./backend"]

EXPOSE 8080