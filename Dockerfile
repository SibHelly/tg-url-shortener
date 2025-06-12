FROM golang:alpine3.22

RUN go version

ENV GOPATH=/
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN go build -ldflags="-w -s" -o tg_shroter ./cmd/main.go

RUN chmod +x ./tg_shroter

CMD ["./tg_shroter"]