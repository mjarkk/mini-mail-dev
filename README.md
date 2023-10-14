# WIP

![Screenshot](/screenshot.png?raw=true "Screnshot")

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
