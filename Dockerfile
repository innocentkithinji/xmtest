FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

copy . .

RUN go build -o main .

EXPOSE 8080

ENTRYPOINT ["./main"]