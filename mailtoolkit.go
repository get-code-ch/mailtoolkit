package mailtoolkit

import (
	"bytes"
	"fmt"
	"net/mail"
)

// Parse parses a raw email.
//
// A malformed mail never makes Parse panic. When the body cannot be fully
// read (truncated multipart, missing boundary...), Parse returns everything
// it could extract together with a non-nil error.
func Parse(buffer []byte) (Mail, error) {
	var m Mail

	msg, err := mail.ReadMessage(bytes.NewReader(buffer))
	if err != nil {
		return m, fmt.Errorf("mailtoolkit: reading header: %w", err)
	}

	m.Header = newHeader(msg.Header)
	m.Contents, m.Attachments, err = parseBody(msg.Body, m.Header.ContentInfo)
	return m, err
}

// ParseHeader parses only the header of a raw email.
func ParseHeader(buffer []byte) (Header, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(buffer))
	if err != nil {
		return Header{}, fmt.Errorf("mailtoolkit: reading header: %w", err)
	}
	return newHeader(msg.Header), nil
}
