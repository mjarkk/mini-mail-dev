# `mini-mail-dev`

![Screenshot](/screenshot.jpg?raw=true "Screnshot")

A alternative take on a development / testing email server.

Features:

- Get started with a single command
- Simple UI
- Fast
- Can handle a vast amounts of emails

## Run

**Docker**

```bash
docker run -it --rm -p 1080:1080 -p 1025:1025 ghcr.io/mjarkk/mini-mail-dev:latest
```

**Docker build**

```bash
git clone https://github.com/mjarkk/mini-mail-dev
cd mini-mail-dev
docker build -t mini-mail-dev .
docker run -it --rm -p 1080:1080 -p 1025:1025 mini-mail-dev
```

**Go Install**

_(you might not have the latest version if you have installed it earlier)_

```bash
go install github.com/mjarkk/mini-mail-dev@latest
mini-mail-dev
```

**Git clone & Go build**

```bash
git clone https://github.com/mjarkk/mini-mail-dev
cd mini-mail-dev
go build
./mini-mail-dev
```

## Usage

```sh
mini-mail-dev [--flags]
```

Options:

| name                   | Flag                   | Envourment variable  | Default        |
| ---------------------- | ---------------------- | -------------------- | -------------- |
| Disable web            | `--disable-web`        | `DISABLE_WEB`        | `false`        |
| HTTP address           | `--http`               | `HTTP_ADDR`          | `0.0.0.0:1080` |
| HTTP password          | `--http-pass`          | `HTTP_PASS`          |                |
| HTTP username          | `--http-user`          | `HTTP_USER`          |                |
| Max emails             | `--max-emails`         | `MAX_EMAILS`         | `200`          |
| SMTP address           | `--smtp`               | `SMTP_ADDR`          | `0.0.0.0:1025` |
| SMTP domain            | `--smtp-domain`        | `SMTP_DOMAIN`        | `localhost`    |
| SMTP incoming password | `--smtp-incoming-pass` | `SMTP_INCOMING_PASS` |                |
| SMTP incoming username | `--smtp-incoming-user` | `SMTP_INCOMING_USER` |                |
| SMTP TLS               | `--smtp-tls`           | `SMTP_TLS`           | `false`        |

_Note that authentication is only enabled if you provided the relevant envourment variables_

## Why

Something something reinventing the wheel :^)

But besides the fun of reinventing the wheel I do have some reasons for making this.
I used maildev before this and I have some annoyances with it. It's an awesome project and has improved my life, but after using it for many years I have also noticed it's not without its flaws. So here is `mini-mail-dev`, my take on a mail development server and with fixes for my maildev problems.

Here are things this project aims to improve / do different compared to maildev:

- Attachments not hidden behind a button
- By default a fixed max amount of emails and when exceeded old emails are removed. _(So you can safely deploy this on a server without much storage or cpu as the server does not have to store a shitload of emails when sending a shitload of emails on a staging server hehe)_
- Different server side language. _Mail dev has a tendency to break when running it on a server and sending more than the "average" amound of mails (my experience is at around 5.000 emails). By applying the above and using a language that is better suited for these kinds of things I can deploy this on every potato in the world._

## Development:

Requirements:

- Golang v1.21 +
- Nodejs

### Dev

```bash
go run .
# In another terminal tab run:
npm run dev

# For sending a few testing mails to the server run
go run test/test_email.go
```

_The javascript files are hot module reloaded only when going to [localhost:**3001**](http://localhost:3001), the go files are not reloaded on changes made_

### Create and run a release build

```bash
# Build frontend
npm i
npm run build

# Build backend
go get
go build
```

Now you can start `mini-mail-dev` using

```bash
./mini-mail-dev
```
