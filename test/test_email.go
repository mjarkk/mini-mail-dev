package main

import (
	"crypto/tls"
	_ "embed"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

//go:embed client-compatibility.html
var clientCompatibilityHTML string

func main() {
	client := mail.NewSMTPClient()

	client.Host = "localhost"
	client.Port = 1025
	client.Username = "username"
	client.Password = "password"
	client.Encryption = mail.EncryptionNone
	client.KeepAlive = true
	client.ConnectTimeout = 10 * time.Second
	client.SendTimeout = 10 * time.Second
	client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpClient, err := client.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	defer smtpClient.Close()

	image, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAAaADAAQAAAABAAAAAQAAAAD5Ip3+AAAADUlEQVQIHWP4v5ThPwAG7wKkSFotfwAAAABJRU5ErkJggg==")
	if err != nil {
		log.Fatalln(err)
	}

	imageFile := mail.File{
		FilePath: "./example.png",
		Name:     "example.png",
		MimeType: "image/png",
		Data:     image,
		Inline:   false,
	}

	email := mail.NewMSG().
		SetFrom("Example sender <sender@example.org>").
		AddTo("recipient@example.net").
		SetSubject("Test sending mail").
		AddAlternative(mail.TextPlain, "Hello world!\nThis is test\n")

	mustSendEmail(smtpClient, email)

	for i := 0; i < 2; i++ {
		email := mail.NewMSG().
			SetFrom(fmt.Sprintf("Example sender <sender-%d@example.org>", i)).
			AddTo(fmt.Sprintf("recipient-%d@example.net", i)).
			SetSubject(fmt.Sprintf("Test sending mail #%d", i)).
			SetBody(mail.TextHTML, fmt.Sprintf("<h1>Hello world!</h1><p>This is test #%d</p>", i)).
			AddAlternative(mail.TextPlain, fmt.Sprintf("Hello world!\nThis is test #%d\n", i)).
			Attach(&imageFile)

		mustSendEmail(smtpClient, email)
	}

	screenshotBytes, err := os.ReadFile("./screenshot.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	screenshot := mail.File{
		FilePath: "./screenshot.jpg",
		Name:     "screenshot.jpg",
		MimeType: "image/jpeg",
		Data:     screenshotBytes,
		Inline:   false,
	}

	email = mail.NewMSG().
		SetFrom("Example sender <attachment-sender@example.org>").
		AddTo("attachment-recipient@example.net").
		SetSubject("Test sending attachment").
		SetBody(mail.TextHTML, "<h1>Attachment test!</h1>").
		Attach(&imageFile).
		Attach(&screenshot)

	mustSendEmail(smtpClient, email)

	email = mail.NewMSG().
		SetFrom("Example sender <compatibility@example.org>").
		AddTo("recipient@example.net").
		SetSubject("Email client compatibility test").
		SetBody(mail.TextHTML, clientCompatibilityHTML).
		AddAlternative(mail.TextPlain, "Switch between email clients to compare how they render this message.\n")

	mustSendEmail(smtpClient, email)

	// Gmail clips messages whose HTML exceeds 102 KB.
	paragraph := "<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>\n"
	email = mail.NewMSG().
		SetFrom("Example sender <long@example.org>").
		AddTo("recipient@example.net").
		SetSubject("Very long email").
		SetBody(mail.TextHTML, "<h1>Very long email</h1>\n"+strings.Repeat(paragraph, 1000)+"<p><b>The end.</b></p>")

	mustSendEmail(smtpClient, email)

	// nr := uint(0)
	// for {
	// 	time.Sleep(time.Millisecond * 10)
	// 	email := mail.NewMSG().
	// 		SetFrom(fmt.Sprintf("Stress test %d <stress-%d@example.org>", nr, nr)).
	// 		AddTo(fmt.Sprintf("stress-recipient-%d@example.net", nr)).
	// 		SetSubject("Test sending mail").
	// 		AddAlternative(mail.TextPlain, fmt.Sprintf("Hello world!\nThis is a stress test #%d\n", nr))
	// 	mightSendEmail(smtpClient, email)
	// 	nr++
	// }

	smtpClient.Close()
}

func mustSendEmail(smtpClient *mail.SMTPClient, email *mail.Email) {
	if email.Error != nil {
		log.Fatalln(email.Error)
	}

	err := email.Send(smtpClient)
	if err != nil {
		log.Fatalln("Failed to send mail", err)
	}
}

func mightSendEmail(smtpClient *mail.SMTPClient, email *mail.Email) {
	if email.Error != nil {
		log.Fatalln(email.Error)
	}

	err := email.Send(smtpClient)
	if err != nil {
		fmt.Println("WARN: Failed to send mail", err)
	}
}
