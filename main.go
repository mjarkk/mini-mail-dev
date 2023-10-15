package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	. "github.com/mjarkk/mini-mail-dev/go"
	"github.com/spf13/pflag"
)

//go:embed dist
var dist embed.FS

func getenv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func main() {
	ifEmptyNotRequired := ", if empty no credentials required"

	defaultSMTPAddr := "localhost:1025"
	defaultHTTPAddr := "localhost:1080"
	defaultSMPTDomain := "localhost"

	argSMTPAddr := pflag.String("smtp", defaultSMTPAddr, "SMTP server address")
	argHTTPAddr := pflag.String("http", defaultHTTPAddr, "HTTP server address")
	argSMTPDomain := pflag.String("smtp-domain", defaultSMPTDomain, "SMTP server domain")
	argSMTPUsername := pflag.String("smtp-incoming-user", "", "SMTP server username"+ifEmptyNotRequired)
	argSMTPPassword := pflag.String("smtp-incoming-pass", "", "SMTP server password"+ifEmptyNotRequired)
	argHTTPUsername := pflag.String("http-user", "", "HTTP server username"+ifEmptyNotRequired)
	argHTTPPassword := pflag.String("http-pass", "", "HTTP server address"+ifEmptyNotRequired)
	argDisableWeb := pflag.Bool("disable-web", false, "Disable the web interface")
	pflag.Parse()

	envSMTPAddr := getenv("SMTP_ADDR")
	envHTTPAddr := getenv("HTTP_ADDR")
	envSMTPDomain := getenv("SMTP_DOMAIN")
	envSMTPUsername := getenv("SMTP_INCOMING_USER")
	envSMTPPassword := getenv("SMTP_INCOMING_PASS")
	envHTTPUsername := getenv("HTTP_USER")
	envHTTPPassword := getenv("HTTP_PASS")
	envDisableWeb := strings.ToLower(getenv("DISABLE_WEB")) == "true"

	if *argSMTPAddr == defaultHTTPAddr && envSMTPAddr != "" {
		*argSMTPAddr = envSMTPAddr
	}
	if *argHTTPAddr == defaultHTTPAddr && envHTTPAddr != "" {
		*argHTTPAddr = envHTTPAddr
	}
	if *argSMTPDomain == defaultSMPTDomain && envSMTPDomain != "" {
		*argSMTPDomain = envSMTPDomain
	}
	if *argSMTPUsername == "" && envSMTPUsername != "" {
		*argSMTPUsername = envSMTPUsername
	}
	if *argSMTPPassword == "" && envSMTPPassword != "" {
		*argSMTPPassword = envSMTPPassword
	}
	if *argHTTPUsername == "" && envHTTPUsername != "" {
		*argHTTPUsername = envHTTPUsername
	}
	if *argHTTPPassword == "" && envHTTPPassword != "" {
		*argHTTPPassword = envHTTPPassword
	}
	if !*argDisableWeb && envDisableWeb {
		*argDisableWeb = envDisableWeb
	}

	if *argDisableWeb {
		fmt.Println("Not starting webserver as --disable-web is set")
	} else {
		go StartWebserver(dist, StartWebserverOptions{
			Addr:              *argHTTPAddr,
			BasicAuthUsername: *argHTTPUsername,
			BasicAuthPassword: *argHTTPPassword,
		})
	}
	StartEmailServer(StartEmailServerOptions{
		Addr:     *argSMTPAddr,
		Domain:   *argSMTPDomain,
		Username: *argSMTPUsername,
		Password: *argSMTPPassword,
	})
}
