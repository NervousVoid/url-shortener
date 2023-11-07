FROM golang:1.21.0-alpine3.18

WORKDIR /urlshortener

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o urlshortener ./cmd/url-shortener

CMD ["./urlshortener"]