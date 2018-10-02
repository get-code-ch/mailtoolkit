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
