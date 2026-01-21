FROM golang:1.25.5-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /log-output

CMD ["/log-output"]