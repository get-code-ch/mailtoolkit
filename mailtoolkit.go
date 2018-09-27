package mailtoolkit

type Mail struct {
	Header     Header
	Content    []Content
	Attachment map[string]Attachment
}

type Header struct {
	IsMime   bool
	From     string
	To       string
	Date     string
	Subject  string
	Elements map[string]string
}

type Content struct {
	ContentInfo ContentInfo
	Data        []byte
	Related     map[string]Attachment
}

type Attachment struct {
	ContentInfo ContentInfo
	Data        []byte
}

type ContentInfo struct {
	Type              string
	Subtype           string
	Parameters        map[string]string
	ID                string
	Description       string
	TransfertEncoding string
}

var HandledContentType = map[string]bool{
	"application": false,
	"audio":       false,
	"image":       false,
	"message":     false,
	"multipart":   true,
	"text":        true,
	"video":       false,
	"x-*":         false,
}

func Parse(buffer []byte) Mail {
	return Mail{}
}

func ParseHeader(buffer []byte) Header {
	return Header{}
}
