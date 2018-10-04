package mailtoolkit

import (
	"fmt"
	"regexp"
)

type Mail struct {
	Header     Header
	Content    []Content
	Attachment map[string]Attachment
}

func Parse(buffer []byte) Mail {
	var re *regexp.Regexp
	var Mail Mail

	Mail.Header = ParseHeader(buffer)

	// Get body start
	re = regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
	BodyStart := re.FindIndex(buffer)

	if Mail.Header.IsMime {
		b, exist := Mail.Header.ContentInfo.Type.Parameters["boundary"]
		if exist {
			re = regexp.MustCompile(`(?mi)^--` + b)
			match := re.FindAllIndex(buffer[BodyStart[1]:], -1)
			if match != nil {
				fmt.Printf("%v\n", match)
			}
		} else {
			Mail.Content = append(Mail.Content, ParseContent(buffer[BodyStart[1]:]))
		}

	} else {
		Mail.Content = append(Mail.Content, ParseContent(buffer[BodyStart[1]:]))

	}
	return Mail
}
