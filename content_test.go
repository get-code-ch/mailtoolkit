package mailtoolkit

import (
	"mime"
	"sort"
	"strings"
	"testing"
)

// crlf converts a test mail to CRLF line endings, as sent over SMTP.
func crlf(s string) string {
	return strings.ReplaceAll(s, "\n", "\r\n")
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestParseCRLFAndSpecialBoundary(t *testing.T) {
	// Boundary characters that broke the former regex based parser.
	raw := crlf(`From: a@example.com
To: b@example.com
Subject: boundary
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="=_a+b.c(d)?e"

--=_a+b.c(d)?e
Content-Type: text/plain; charset=utf-8

Hello
--=_a+b.c(d)?e
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

<p>caf=C3=A9</p>
--=_a+b.c(d)?e--
`)
	mail, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("%v Parse: %v", ballotX, err)
	}
	if got := sortedKeys(mail.Contents); strings.Join(got, ",") != "0,1" {
		t.Fatalf("%v contents = %v, want [0 1]", ballotX, got)
	}
	if got := string(mail.Contents["0"].Data); got != "Hello" {
		t.Errorf("%v text part = %q, want %q", ballotX, got, "Hello")
	}
	html, err := mail.Contents["1"].Decode()
	if err != nil {
		t.Fatalf("%v Decode: %v", ballotX, err)
	}
	if got := string(html); got != "<p>café</p>" {
		t.Errorf("%v html part = %q, want %q", ballotX, got, "<p>café</p>")
	}
	if got := mail.Contents["0"].ContentInfo.Type.Parameters["charset"]; got != "utf-8" {
		t.Errorf("%v charset = %q, want utf-8", ballotX, got)
	}
}

func TestParseNestedAttachments(t *testing.T) {
	encodedName := mime.QEncoding.Encode("utf-8", "résumé.pdf")
	raw := `From: a@example.com
To: b@example.com
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=outer

--outer
Content-Type: multipart/related; boundary=related

--related
Content-Type: multipart/alternative; boundary=alt

--alt
Content-Type: text/plain

plain
--alt
Content-Type: text/html

<img src="cid:img1@x">
--alt--
--related
Content-Type: image/png
Content-ID: <img1@x>
Content-Transfer-Encoding: base64

aGVsbG8gd29ybGQ=
--related--
--outer
Content-Type: text/plain; name=a.txt
Content-Disposition: attachment; filename=a.txt
Content-Transfer-Encoding: base64

aGVsbG8g
d29ybGQ=
--outer
Content-Type: text/plain
Content-Disposition: attachment; filename="a.txt"

second
--outer
Content-Type: application/pdf; name="` + encodedName + `"
Content-Disposition: attachment

%PDF
--outer
Content-Type: text/plain
Content-Disposition: attachment; filename*=UTF-8''caf%C3%A9.txt

rfc2231
--outer--
`
	mail, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("%v Parse: %v", ballotX, err)
	}

	wantContents := "0,1,img1@x"
	if got := strings.Join(sortedKeys(mail.Contents), ","); got != wantContents {
		t.Errorf("%v contents = %s, want %s", ballotX, got, wantContents)
	}
	wantAttachments := "a (2).txt,a.txt,café.txt,résumé.pdf"
	if got := strings.Join(sortedKeys(mail.Attachments), ","); got != wantAttachments {
		t.Fatalf("%v attachments = %s, want %s", ballotX, got, wantAttachments)
	}

	for key, want := range map[string]string{"a.txt": "hello world", "a (2).txt": "second"} {
		data, err := mail.Attachments[key].Decode()
		if err != nil {
			t.Fatalf("%v Decode %s: %v", ballotX, key, err)
		}
		if string(data) != want {
			t.Errorf("%v attachment %s = %q, want %q", ballotX, key, data, want)
		}
	}
	image, err := mail.Contents["img1@x"].Decode()
	if err != nil || string(image) != "hello world" {
		t.Errorf("%v inline image = %q (%v), want %q", ballotX, image, err, "hello world")
	}
}

func TestParseEncodedHeaders(t *testing.T) {
	subject := "Café ✓ résumé"
	raw := `From: =?utf-8?q?Ren=C3=A9?= <rene@example.com>
To: "Alice" <alice@example.com>, bob@example.com
Cc: undisclosed-recipients:;
Subject: ` + mime.BEncoding.Encode("utf-8", subject) + `
Date: Fri, 19 Oct 2018 08:13:36 +0200

body
`
	mail, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("%v Parse: %v", ballotX, err)
	}
	h := mail.Header
	if h.Subject != subject {
		t.Errorf("%v Subject = %q, want %q", ballotX, h.Subject, subject)
	}
	if h.To != "alice@example.com, bob@example.com" {
		t.Errorf("%v To = %q", ballotX, h.To)
	}
	if len(h.ToList) != 2 || h.ToList[0].Name != "Alice" {
		t.Errorf("%v ToList = %v", ballotX, h.ToList)
	}
	if len(h.FromList) != 1 || h.FromList[0].Name != "René" {
		t.Errorf("%v FromList = %v", ballotX, h.FromList)
	}
	if h.Cc != "" {
		t.Errorf("%v Cc = %q, want empty", ballotX, h.Cc)
	}
	if h.Time.IsZero() || h.Time.Year() != 2018 {
		t.Errorf("%v Time = %v", ballotX, h.Time)
	}
	if h.IsMime || h.ContentInfo.Type.Type != "text" || h.ContentInfo.Type.Subtype != "plain" {
		t.Errorf("%v non MIME mail should default to text/plain, got %+v", ballotX, h.ContentInfo.Type)
	}
}

// Malformed mails must never panic; the returned error tells whether the
// mail could be fully read.
func TestParseMalformed(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"empty", "", true},
		{"header only, no blank line", "From: a@example.com\nSubject: x", false},
		{"content type without slash", "Content-Type: text\n\nbody", false},
		{"address field without address", "To: nobody here\n\nbody", false},
		{"multipart without boundary", "Content-Type: multipart/mixed\n\nbody", true},
		{"truncated multipart", "Content-Type: multipart/mixed; boundary=b\n\n--b\nContent-Type: text/plain\n\ncut", true},
		{"boundary never found", "Content-Type: multipart/mixed; boundary=b\n\nno parts", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.raw))
			if (err != nil) != tt.wantErr {
				t.Errorf("%v Parse error = %v, wantErr %v", ballotX, err, tt.wantErr)
			}
		})
	}
}

func TestParseTruncatedKeepsReadParts(t *testing.T) {
	raw := "Content-Type: multipart/mixed; boundary=b\n\n--b\nContent-Type: text/plain\n\nfirst\n--b\nContent-Type: text/plain\n\ncut"
	mail, err := Parse([]byte(raw))
	if err == nil {
		t.Fatalf("%v expected an error for a truncated mail", ballotX)
	}
	if got := string(mail.Contents["0"].Data); got != "first" {
		t.Errorf("%v first part = %q, want %q", ballotX, got, "first")
	}
}
