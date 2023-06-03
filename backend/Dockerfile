FROM golang:1.20.4-alpine
WORKDIR /
COPY /backend .
RUN go build -o main .
EXPOSE 8181
CMD ["/main"]
