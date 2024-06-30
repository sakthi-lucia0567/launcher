FROM golang:1.22-alphine

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]