FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema

RUN go build -o main server.go

EXPOSE 8080

CMD ["./main"]