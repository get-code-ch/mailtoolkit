package mailtoolkit

func Parse(buffer []byte) Mail {
	var Mail Mail

	Mail.Header = ParseHeader(buffer)
	Mail.Content = ParseContent(buffer, Mail.Header.ContentInfo, 1)

	return Mail
}
