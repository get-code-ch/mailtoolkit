package mailtoolkit

import (
	"regexp"
)

func ParseContents(buffer []byte, contentInfo ContentInfo) ([]Content, map[string]Attachment) {
	var start int
	var end int

	contents := []Content{}
	attachment := make(map[string]Attachment)

	// Get first line of part
	re := regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
	contentStart := re.FindIndex(buffer)[1] + 1
	rawContent := buffer[contentStart:]

	// Allways return a RAW mail contents
	contents = append(contents, Content{Data: rawContent, ContentInfo: contentInfo})

	if contentInfo.Type.Type != "multipart" {
		return contents, nil
	}

	// Get contents starting of mail based on boundary
	re = regexp.MustCompile(`(?m)^--` + contentInfo.Type.Parameters["boundary"] + `\s*$`)
	contentStart = re.FindIndex(rawContent)[0]
	rawContent = rawContent[contentStart:]

	// Extract all contents as RAW format, multipart header and separation are ignored
	re = regexp.MustCompile(`(?m)^--` + contentInfo.Type.Parameters["boundary"] + `--\s*$`)
	contentEnd := re.FindIndex(rawContent)[0]

	rawContent = rawContent[:contentEnd]
	contents = append(contents, Content{Data: rawContent})

	// Get
	re = regexp.MustCompile(`(?m)--` + contentInfo.Type.Parameters["boundary"] + `\s*$`)
	indexes := re.FindAllIndex(rawContent, -1)

	for i := range indexes {
		start = indexes[i][1] + 1
		if i < len(indexes)-1 {
			end = indexes[i+1][0]
		} else {
			end = contentEnd
		}
		ci := getContentInfo(rawContent[start:end])

		// Get first line of part
		re := regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
		contentStart := re.FindIndex(rawContent)[1] + 1

		c := rawContent[contentStart:end]
		if ci.Type.Type != "multipart" {
			if ci.Disposition.Type == "attachment" {
				attachment[ci.Disposition.Parameters["filename"]] = Attachment{Data: c, ContentInfo: ci}
			} else {
				contents = append(contents, Content{Data: c, ContentInfo: ci})
			}
		} else {
			contents, _ = ParseContents(c, ci)
		}
	}

	return contents, attachment
}
