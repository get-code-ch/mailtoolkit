package mailtoolkit

import (
	"regexp"
	"strconv"
)

type delimiter struct {
	start int
	end   int
}

func ParseContent(buffer []byte, header Header) map[string]Content {
	var parts []delimiter

	content := make(map[string]Content)

	// Get first line of part
	re := regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
	contentStart := re.FindIndex(buffer)[1] + 1
	rawContent := buffer[contentStart:]

	// Allways return a RAW mail content
	content["raw"] = Content{Data: rawContent, ContentInfo: header.ContentInfo}

	if !header.IsMime || header.ContentInfo.Type.Type != "multipart" {
		return content
	}

	// Get content starting of mail based on boundary
	re = regexp.MustCompile(`(?m)^--` + header.ContentInfo.Type.Parameters["boundary"] + `\s*$`)
	contentStart = re.FindIndex(rawContent)[0]
	rawContent = rawContent[contentStart:]

	// Extract all content as RAW format, multipart header and separation are ignored
	re = regexp.MustCompile(`(?m)^--` + header.ContentInfo.Type.Parameters["boundary"] + `--\s*$`)
	contentEnd := re.FindIndex(rawContent)[0]

	rawContent = rawContent[:contentEnd]
	content["raw"] = Content{Data: rawContent}

	// Get
	re = regexp.MustCompile(`(?m)--` + header.ContentInfo.Type.Parameters["boundary"] + `\s*$`)
	contents := re.FindAllIndex(rawContent, -1)

	for i := range contents {
		delimiter := delimiter{start: contents[i][1] + 1}
		parts = append(parts, delimiter)
		if i < len(contents)-1 {
			parts[i].end = contents[i+1][0]
		} else {
			parts[i].end = contentEnd
		}
		c := Content{Data: rawContent[parts[i].start:parts[i].end], ContentInfo: getContentInfo(rawContent[parts[i].start:parts[i].end])}
		if c.ContentInfo.Type.Type != "multipart" {
			content[strconv.Itoa(i)] = c
		} else {
			content[strconv.Itoa(i)] = c
		}
	}

	return content
}
