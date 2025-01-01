FROM golang:1.23.4 AS build-stage 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bookstore .

FROM alpine:latest  
WORKDIR /root/
COPY --from=build-stage /bookstore .
CMD ["./bookstore"]