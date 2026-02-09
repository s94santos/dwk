FROM golang:1.25.5-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN mkdir /app/cache

COPY . .

RUN chmod +x ./scripts/random.sh

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/project

CMD ["/app/project"]