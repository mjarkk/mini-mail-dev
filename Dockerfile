ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

FROM --platform=linux/amd64 node:21-slim AS frontend

COPY . /app
WORKDIR /app

RUN npm i && npm run build

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum .

RUN go mod download

COPY go go
COPY main.go .
COPY --from=frontend /app/dist dist

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /usr/bin/mini-mail-dev && \
    rm -rf go main.go go.mod go.sum dist

EXPOSE 1025/tcp
EXPOSE 1080/tcp

CMD ["/usr/bin/mini-mail-dev"]
