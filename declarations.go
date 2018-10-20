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
	Bcc         string
	Date        string
	Subject     string
	Elements    map[string]string
	ContentInfo ContentInfo
}

type Content struct {
	ContentInfo ContentInfo
	Data        []byte
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
var headerRegex = regexp.MustCompile(`(?mi)(^[\w_-]+)(?::\s+)(.*)(?:\r?\n)((?:\s*(?:\s+).*(?:\r?\n+))*)`)
var parametersRegex = regexp.MustCompile(`(?i)^(?:\s*)([\w-]+)(?:\s*=\s*"?)(.*[^"])(?:"?\s*)$`)
var whitespaceRegex = regexp.MustCompile(`[ ]{2,}|[\t|\0|\n|\r]+`)
var quotesRegex = regexp.MustCompile(`^([<|"|']?)+|([[>|"|']?)+[\r|\n|\s]+$`)
var contentTransferEncodingRegex = regexp.MustCompile(`(?mi)(?:^\s*Content-Transfer-Encoding:\s+"?)(.*)(?:"?\n?)`)
var contentIDRegex = regexp.MustCompile(`(?mi)(?:^\s*Content-ID:\s+"?)(.*)(?:"?\n?)`)
var contentDescriptionRegex = regexp.MustCompile(`(?mi)(?:^\s*Content-Description:\s+"?)(.*)(?:"?\n?)`)
var contentTypeRegex = regexp.MustCompile(`(?im)(?:^Content-Type: ?)(.+)(?:\r?\n)((?:\s*(?:\s+).*(?:\r?\n+))*)`)
var contentDispositionRegex = regexp.MustCompile(`(?im)(?:^Content-Disposition: ?)(.+)(?:\r?\n)((?:\s*(?:\s+).*(?:\r?\n+))*)`)
var semiColonRegex = regexp.MustCompile(`;`)
var slashRegex = regexp.MustCompile(`/`)
var cidRegex = regexp.MustCompile(`(?im)(?:cid:)((?:[^"]|\\")*)`)
var emailRegex = regexp.MustCompile(`(?i)(?:\<|\"|\x60|\')?([\w-\+\.]*@[a-z0-9][\w-\+\.]*)(?:\>|\"|\x60|\')?`)
