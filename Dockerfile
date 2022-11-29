FROM golang:1.18

RUN mkdir /app
WORKDIR /app


#RUN go install github.com/go-sql-driver/mysql@latest

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod tidy


COPY ./ ./

RUN go build -o /main

EXPOSE 8080

CMD ["/main"]