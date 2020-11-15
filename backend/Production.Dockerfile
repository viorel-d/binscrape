# stage 1
FROM golang:alpine as builder

RUN mkdir /backend/
WORKDIR /backend/

COPY . .

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build ./...

# stage 2
FROM alpine:latest

RUN mkdir /binscrape/
WORKDIR /binscrape/

COPY --from=builder /backend/binscrape .
COPY --from=builder /backend/api .

CMD ["./run.sh"]
