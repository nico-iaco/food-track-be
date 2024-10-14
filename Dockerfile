FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /food-track-be

EXPOSE 8080

FROM alpine

COPY --from=builder /food-track-be .

CMD [ "/food-track-be" ]