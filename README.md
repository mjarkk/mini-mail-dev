# `mini-mail-dev`

![Screenshot](/screenshot.jpg?raw=true "Screnshot")

A alternative take on a development testing email server.

## Why

Something something reinventing the wheel :^)

But besides the fun of reinventing the wheel I do have some reasons for making this.
I used maildev before this and I have some annoyances with it. It's an awesome project and has improved my life, but after using it for many years I have also noticed it's not without its flaws. So here is `mini-mail-dev`, my take on a mail development server and with fixes for my maildev problems.

Here are things this project aims to improve / do different compared to maildev:

- Attachments not hidden behind a button
- By default a fixed max amount of emails and when exceeded old emails are removed. _(So you can safely deploy this on a server without much storage or cpu as the server does not have to store a shitload of emails when sending a shitload of emails on a staging server hehe)_
- Different server side language. _Mail dev has a tendency to break when running it on a server and sending more than the "average" amound of mails (my experience is at around 5.000 emails). By applying the above and using a language that is better suited for these kinds of things I can deploy this on every potato in the world._

## Install

**Go Install**

_(you might not have the latest version if you have installed it earlier)_

```bash
go install github.com/mjarkk/mini-mail-dev@latest
mini-mail-dev
```

**Git clone & Go install**

```bash
git clone https://github.com/mjarkk/mini-mail-dev
cd mini-mail-dev
go install
mini-mail-dev
```

**Docker**

```bash
docker build -t mini-mail-dev https://github.com/mjarkk/mini-mail-dev.git#main
docker run -it --rm -p 1080:1080 -p 1025:1025 mini-mail-dev
```

**Git clone & Docker**

```bash
git clone https://github.com/mjarkk/mini-mail-dev
cd mini-mail-dev
docker build -t mini-mail-dev .
docker run -it --rm -p 1080:1080 -p 1025:1025 mini-mail-dev
```

## Usage

```sh
mini-mail-dev
```

Options:

```
Usage of mini-mail-dev
    --disable-web                 Disable the web interface
    --http string                 HTTP server address (default "localhost:1080")
    --http-pass string            HTTP server address, if empty no credentials required
    --http-user string            HTTP server username, if empty no credentials required
    --max-emails uint16           The max amount of emails to keep (default 200)
    --smtp string                 SMTP server address (default "localhost:1025")
    --smtp-domain string          SMTP server domain (default "localhost")
    --smtp-incoming-pass string   SMTP server password, if empty no credentials required
    --smtp-incoming-user string   SMTP server username, if empty no credentials required
```

Environment variables

_These are used the same way as the arguments_

```
DISABLE_WEB
HTTP_ADDR
HTTP_PASS
HTTP_USER
MAX_EMAILS
SMTP_ADDR
SMTP_DOMAIN
SMTP_INCOMING_PASS
SMTP_INCOMING_USER
```

## Build & Run

Requirements:

- Golang v1.21.2 +
- Nodejs

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

## Dev

```bash
go run .
# In another terminal tab run:
npm run dev

# For sending a few testing mails to the server run
go run test/test_email.go
```

_The javascript files are hot module reloaded only when going to [localhost:**3001**](http://localhost:3001), the go files are not reloaded on changes made_
