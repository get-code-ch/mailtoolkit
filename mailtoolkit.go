package mailtoolkit

func Parse(buffer *[]byte) Mail {
	var Mail Mail

	Mail.Header = ParseHeader(buffer)
	Mail.Content = ParseContent(buffer, Mail.Header)

	return Mail
}
