package mailtoolkit

type Content struct {
	ContentInfo ContentInfo
	Data        []byte
	Content     []Attachment
}

type Attachment struct {
	ContentInfo ContentInfo
	Data        []byte
}

func ParseContent(buffer []byte) Content {
	content := Content{}
	content.Data = buffer
	return content
}
