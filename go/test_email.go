package src

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// TestSendingMail test sending some mock emails
func TestSendingMail() {
	time.Sleep(time.Second)

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

	email := mail.NewMSG()
	email.SetFrom("Example sender <sender@example.org>").
		AddTo("recipient@example.net").
		SetSubject("Test sending mail").
		AddAlternative(mail.TextPlain, "Hello world!\nThis is test\n")

	if email.Error != nil {
		log.Fatalln(email.Error)
	}

	err = email.Send(smtpClient)
	if err != nil {
		log.Fatalln("Failed to send mail", err)
	}

	for i := 0; i < 2; i++ {
		email := mail.NewMSG()
		email.SetFrom(fmt.Sprintf("Example sender <sender-%d@example.org>", i)).
			AddTo(fmt.Sprintf("recipient-%d@example.net", i)).
			SetSubject(fmt.Sprintf("Test sending mail #%d", i)).
			SetBody(mail.TextHTML, fmt.Sprintf("<h1>Hello world!</h1><p>This is test #%d</p>", i)).
			AddAlternative(mail.TextPlain, fmt.Sprintf("Hello world!\nThis is test #%d\n", i)).
			Attach(&imageFile)

		if email.Error != nil {
			log.Fatalln(email.Error)
		}

		err = email.Send(smtpClient)
		if err != nil {
			log.Fatalln("Failed to send mail", err)
		}
	}

	smtpClient.Close()
}
