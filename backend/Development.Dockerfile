FROM golang:stretch

RUN mkdir /backend/
WORKDIR /backend/

COPY . .

RUN apt-get update && apt-get install -y inotify-tools

ENV GO111MODULE=on
ENV GOBIN=$GOPATH/bin
RUN CGO_ENABLED=0 GOOS=linux go install ./...

CMD ["/bin/bash", "-c", "./watch.sh binscrape api"]
