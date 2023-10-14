package main

import (
	"embed"

	. "github.com/mjarkk/mini-mail-dev/go"
)

//go:embed dist
var dist embed.FS

func main() {
	go StartWebserver(dist)
	StartEmailServer()
}
