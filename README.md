# `mini-mail-dev`

> **This project is work in progress!**

![Screenshot](/screenshot.jpg?raw=true "Screnshot")

A alternative take on a development testing email server.

## Why

Something something re-inventing the wheel :^)

But beside the fun of re-inventing the wheel i do have some reasons for making this.
I used mail dev before this and i have some annoyences with it, it's an awesome project and has improved my live but after using it for many years i have also noticed it's not without it's flaws so here is `mini-mail-dev` to fix my problems with it and my take on a mail development server.

Here are things this project aims improve / do different compared to maildev.

- Attachments not hidden behind a button
- By default a fixed max amound of emails and when exceeded old emails are removed. _(So you can safely deploy this on a server without much storage or cpu as the server does not have to store a shitload of emails when sending a shitlaod of emails on a staging server hehe)_
- Diffrent server side language. _Mail dev has a tendency to break when running it on a server and sending more than the "average" amound of mails (my experiance is at around 5000 mails). By applying the above and using a language that is more made for this kinds of things i can deploy this on every potato in the world._

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

**Docker WIP**

_WIP_

```bash
# docker build -t mini-mail-dev https://github.com/mjarkk/mini-mail-dev.git#main
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

Envourment variables

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

Now you can start mini-mail-dev using

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

_The javascript file are hot module reloaded only when going to [localhost:**3001**](http://localhost:3001), the go files are not reloaded on changes made to them_
