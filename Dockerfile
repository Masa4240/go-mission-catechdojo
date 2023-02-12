FROM golang:1.18

RUN mkdir /app
WORKDIR /app


# RUN go install github.com/go-sql-driver/mysql@latest

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy
RUN go mod download


COPY ./ ./

#RUN go build -o app main.go

EXPOSE 8080

#CMD ["/app"]