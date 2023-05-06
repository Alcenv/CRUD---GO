FROM golang:1.19

WORKDIR /app
COPY go.mod go.sum ./
COPY main.go ./main.go
COPY plantillas ./plantillas

RUN go mod download

EXPOSE 3000