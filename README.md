# mailtoolkit
[![CI](https://github.com/get-code-ch/mailtoolkit/actions/workflows/ci.yml/badge.svg)](https://github.com/get-code-ch/mailtoolkit/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/get-code-ch/mailtoolkit.svg)](https://pkg.go.dev/github.com/get-code-ch/mailtoolkit)

GO mail toolkit to parse email content to an object.

Built on the standard library only (`net/mail`, `mime`, `mime/multipart`), requires Go 1.23 or later.

```go
buffer, err := os.ReadFile("mail.eml")
if err != nil {
	log.Fatal(err)
}
mail, err := mailtoolkit.Parse(buffer)
if err != nil {
	// The mail is malformed or truncated, mail still holds what could be read.
	log.Print(err)
}
fmt.Println(mail.Header.From, mail.Header.Subject)
for name, attachment := range mail.Attachments {
	data, err := attachment.Decode() // removes base64 / quoted-printable
	...
}
```

## Breaking changes since the 2018 version

- `Parse` and `ParseHeader` return an error instead of panicking on malformed mails.
- `ParseContents` and the exported regular expressions / `MIMEContentTypes` variables are gone.
- `Header.To`, `Cc` and `Bcc` hold every address (comma separated), not only the first one. The new `ToList`, `CcList`, ... fields keep the display names.
- `Subject`, `Content-Description` and attachment names are decoded (RFC 2047 / RFC 2231).
- Part data no longer includes the line break preceding the next boundary.
