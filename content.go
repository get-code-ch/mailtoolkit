package mailtoolkit

import (
	"regexp"
	"strconv"
)

type delimiter struct {
	start int
	end   int
}

func ParseContent(buffer []byte, contentInfo ContentInfo, depth int) map[string]Content {
	var parts []delimiter

	content := make(map[string]Content)

	if depth < 1 {
		depth = 1
	}

	// Get first line of part
	re := regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
	contentStart := re.FindIndex(buffer)[1] + 1
	rawContent := buffer[contentStart:]

	// Allways return a RAW mail content
	content["raw"] = Content{Data: rawContent, ContentInfo: contentInfo}

	if contentInfo.Type.Type != "multipart" {
		return content
	}

	// Get content starting of mail based on boundary
	re = regexp.MustCompile(`(?m)^--` + contentInfo.Type.Parameters["boundary"] + `\s*$`)
	contentStart = re.FindIndex(rawContent)[0]
	rawContent = rawContent[contentStart:]

	// Extract all content as RAW format, multipart header and separation are ignored
	re = regexp.MustCompile(`(?m)^--` + contentInfo.Type.Parameters["boundary"] + `--\s*$`)
	contentEnd := re.FindIndex(rawContent)[0]

	rawContent = rawContent[:contentEnd]
	content["raw"] = Content{Data: rawContent}

	// Get
	re = regexp.MustCompile(`(?m)--` + contentInfo.Type.Parameters["boundary"] + `\s*$`)
	contents := re.FindAllIndex(rawContent, -1)

	for i := range contents {
		delimiter := delimiter{start: contents[i][1] + 1}
		parts = append(parts, delimiter)
		if i < len(contents)-1 {
			parts[i].end = contents[i+1][0]
		} else {
			parts[i].end = contentEnd
		}
		ci := getContentInfo(rawContent[parts[i].start:parts[i].end])
		c := rawContent[parts[i].start:parts[i].end]
		if ci.Type.Type != "multipart" {
			if ci.Disposition.Type == "attachment" {
				/*
					_, ok := content[strconv.Itoa(i*depth)]
					if !ok {
						content[strconv.Itoa(i*depth)] = Content{}
						content[strconv.Itoa(i*depth)].Attachments[ci.Disposition.Parameters["filename"]] = Attachment{}
					}
					a := Attachment{Data: c, ContentInfo: ci}
					content[strconv.Itoa(i*depth)].Attachments[ci.Disposition.Parameters["filename"]] = a
				*/
				content[strconv.Itoa(i*depth)] = Content{Data: c, ContentInfo: ci}
			} else {
				content[strconv.Itoa(i*depth)] = Content{Data: c, ContentInfo: ci}
			}
		} else {
			content = ParseContent(c, ci, depth*10)
		}
	}

	return content
}
