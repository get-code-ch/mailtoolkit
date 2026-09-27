// Package mailtoolkit parses a raw RFC 5322 / MIME email into a Mail object.
package mailtoolkit

import (
	"net/mail"
	"time"
)

// Mail is a parsed email.
type Mail struct {
	Header Header
	// Contents are the displayable parts of the mail, keyed by Content-ID
	// (without angle brackets) or, when the part has no Content-ID, by its
	// position ("0", "1", ...).
	Contents map[string]Content
	// Attachments are keyed by their decoded file name, made unique when
	// several attachments share the same name.
	Attachments map[string]Attachment
}

// Header holds the top-level header fields of a mail.
type Header struct {
	IsMime bool
	// From, To, Cc and Bcc are comma separated lists of e-mail addresses.
	// When the mail has no To field, To falls back to Delivered-To.
	From string
	To   string
	Cc   string
	Bcc  string
	// FromList, ToList, CcList and BccList keep the display names.
	FromList []*mail.Address
	ToList   []*mail.Address
	CcList   []*mail.Address
	BccList  []*mail.Address
	// Date is the raw Date field, Time its parsed value (zero if unparsable).
	Date    string
	Time    time.Time
	Subject string
	// Elements contains every header field, keyed by lowercase name.
	// Repeated fields are joined with "\n".
	Elements    map[string]string
	ContentInfo ContentInfo
}

// Content is a displayable part of a mail. Data is still transfer-encoded,
// use Decode to get the raw bytes.
type Content struct {
	ContentInfo ContentInfo
	Data        []byte
}

// Attachment is an attached file. Data is still transfer-encoded, use Decode
// to get the raw bytes.
type Attachment struct {
	ContentInfo ContentInfo
	Data        []byte
}

// ContentInfo describes a MIME part.
type ContentInfo struct {
	Type             ContentType
	ID               string
	Description      string
	TransferEncoding string
	Disposition      ContentDisposition
}

// ContentType is a parsed Content-Type field. Type, Subtype and parameter
// names are lowercase.
type ContentType struct {
	Type       string
	Subtype    string
	Parameters map[string]string
}

// ContentDisposition is a parsed Content-Disposition field. Type and
// parameter names are lowercase.
type ContentDisposition struct {
	Type       string
	Parameters map[string]string
}
