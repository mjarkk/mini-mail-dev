package main

import . "github.com/mjarkk/mini-mail-server/go"

func main() {
	go StartWebserver()
	go TestSendingMail()
	StartEmailServer()
}
