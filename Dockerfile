FROM golang:1.15.3-alpine3.12

RUN mkdir /app

RUN mkdir /.contract

COPY .contract/contract.json /.contract/

ADD . /app

WORKDIR /app

ENV CGO_ENABLED=0

RUN go build -o main .

CMD ["/app/main"]
