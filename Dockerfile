FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o /myapp

CMD ["/myapp"]