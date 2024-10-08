FROM golang:latest

RUN apt update && apt install -y git

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go get github.com/githubnemo/CompileDaemon

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/githubnemo/CompileDaemon
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ENV PATH="${PATH}:$HOME/go/bin"
ENV PATH="${PATH}:/usr/local/go/bin"
ENV TZ="Asia/Tbilisi"

COPY . .

RUN git config --global --add safe.directory /app

CMD ["./bin/start.sh"]
