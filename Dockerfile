FROM node:20-slim AS frontend

COPY . .

RUN npm install && npm run build

FROM golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum .

RUN go mod download

COPY go go
COPY main.go .
COPY --from=frontend dist dist

RUN go build -o /usr/bin/mini-mail-dev && \
    rm -rf go main.go go.mod go.sum dist

EXPOSE 1025/tcp
EXPOSE 1080/tcp

CMD ["/usr/bin/mini-mail-dev"]
