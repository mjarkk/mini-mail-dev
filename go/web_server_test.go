package src

import (
	"strings"
	"testing"
	"time"

	"github.com/jhillyerd/enmime/v2"
	"github.com/oklog/ulid/v2"
)

func convertTestEmail(t *testing.T, raw, realFrom string) Email {
	t.Helper()
	env, err := enmime.ReadEnvelope(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	return ConvertEmail(ulid.Make(), []byte(raw), env, time.Now(), realFrom, "")
}

func TestSearchEmails(t *testing.T) {
	// A bounce: no Sender header and an empty MAIL FROM, so Sender stays nil.
	bounce := convertTestEmail(t, "Subject: Undeliverable\r\nFrom: mailer-daemon@example.com\r\nTo: me@example.com\r\n\r\nYour message could not be delivered.\r\n", "")
	if bounce.Sender != nil {
		t.Fatal("bounce has a sender")
	}
	welcome := convertTestEmail(t, "Subject: Welcome\r\nTo: me@example.com\r\n\r\nHello!\r\n", "shop@example.com")
	list := []Email{bounce, welcome}

	cases := []struct {
		search string
		want   []string
	}{
		{"", []string{"Welcome", "Undeliverable"}},
		{"  COULD NOT ", []string{"Undeliverable"}},
		{"shop@", []string{"Welcome"}},
		{"no such text", nil},
	}
	for _, c := range cases {
		var got []string
		for _, email := range searchEmails(list, c.search) {
			got = append(got, email.Subject)
		}
		if strings.Join(got, ",") != strings.Join(c.want, ",") {
			t.Errorf("searchEmails(%q) = %q, want %q", c.search, got, c.want)
		}
	}
}

func TestNotifyNewEmailWithoutReader(t *testing.T) {
	ch := subscribeToNewEmails()
	defer unsubscribeFromNewEmails(ch)

	notifyNewEmail()
	notifyNewEmail()

	<-ch
	select {
	case <-ch:
		t.Error("two notifications pending, want them merged into one")
	default:
	}
}
