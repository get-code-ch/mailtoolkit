package mailtoolkit

type Mail struct {
	Header     Header
	Content    map[string]Content
	Attachment map[string]Attachment
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
	Content     []Attachment
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

var MIMEContentTypes = []string{"application", "audio", "image", "multipart", "text", "video"}
var MIMEContentDispositionTypes = []string{"inline", "attachment"}
