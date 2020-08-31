# Start from golang base image
FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -v ./cmd/game-keys-store

FROM alpine
WORKDIR /app

COPY --from=builder /app .

CMD ["./game-keys-store"]


#RUN mkdir /app
#WORKDIR /app
#
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#COPY . .
#RUN go build -v ./cmd/game-keys-store
#
#CMD ["./apiserver"]