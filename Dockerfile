FROM golang:1.18

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go test ./markets ./discord -v -count=1
RUN go build -v -o /usr/local/bin/monitor .
