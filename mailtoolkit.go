package mailtoolkit

func Parse(buffer []byte) Mail {
	var Mail Mail

	Mail.Header = ParseHeader(buffer)
	Mail.Contents, Mail.Attachments = ParseContents(buffer, Mail.Header.ContentInfo, nil)

	return Mail
}
