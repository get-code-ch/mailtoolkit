package mailtoolkit

import "regexp"

type Mail struct {
	Header      Header
	Contents    map[string]Content
	Attachments map[string]Attachment
}

type Header struct {
	IsMime      bool
	From        string
	To          string
	Cc          string
	Cci         string
	Date        string
	Subject     string
	Elements    map[string]string
	ContentInfo ContentInfo
}

type Content struct {
	ContentInfo ContentInfo
	Data        []byte
	Contents    map[string]Content
}

type Attachment struct {
	ContentInfo ContentInfo
	Data        []byte
}

type ContentInfo struct {
	Type             ContentType
	ID               string
	Description      string
	TransferEncoding string
	Disposition      ContentDisposition
}

type ContentType struct {
	Type       string
	Subtype    string
	Parameters map[string]string
}

type ContentDisposition struct {
	Type       string
	Parameters map[string]string
}

// Global variable

var MIMEContentTypes = []string{"application", "audio", "image", "multipart", "text", "video"}
var MIMEContentDispositionTypes = []string{"inline", "attachment"}

var firstLineRegex = regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
var headerRegex = regexp.MustCompile(`(?mi)(^[\w_-]+)(?::\s+"?)(.*)(?:\r?\n)((?:\s*(?:\s+).*(?:\r?\n+))*)`)
var whitespaceRegex = regexp.MustCompile(`[ ]{2,}|[\t|\0|\n|\r]+`)
