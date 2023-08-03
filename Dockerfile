FROM golang:1.19.11-alpine

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /TwootService

EXPOSE 10000

CMD [ "/TwootService" ]
