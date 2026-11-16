FROM --platform=$BUILDPLATFORM node:26-slim AS frontend

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.27 AS backend

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY go go
COPY main.go .
COPY --from=frontend /app/dist dist
# scratch has no zoneinfo; timetzdata embeds it so the TZ env var still applies.
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=${TARGETVARIANT#v} \
    go build -trimpath -tags timetzdata -ldflags="-s -w" -o /mini-mail-dev

FROM scratch

COPY --from=backend /mini-mail-dev /usr/bin/mini-mail-dev

EXPOSE 1025/tcp
EXPOSE 1080/tcp

CMD ["/usr/bin/mini-mail-dev"]
