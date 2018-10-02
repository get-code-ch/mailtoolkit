package mailtoolkit

type Mail struct {
	Header     Header
	Content    []Content
	Attachment map[string]Attachment
}

func Parse(buffer []byte) Mail {
	var Mail Mail

	Mail.Header = ParseHeader(buffer)

	return Mail
}

func ParseBody(buffer []byte) Content {
	return Content{}
}
