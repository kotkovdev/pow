FROM golang:1.21

WORKDIR /app

RUN mkdir /app/bin

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /app/bin ./...

EXPOSE 8080

RUN chmod +x /app/bin