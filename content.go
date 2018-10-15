package mailtoolkit

import (
	"regexp"
	"strconv"
	_ "strconv"
)

func ParseContents(buffer []byte, contentInfo ContentInfo, contentID *int) (map[string]Content, map[string]Attachment) {
	var start int
	var end int
	var indexes [][]int
	var contentStart int
	var contentEnd int
	var ci ContentInfo
	var c []byte

	contents := make(map[string]Content)
	attachment := make(map[string]Attachment)

	if contentID == nil {
		contentID = new(int)
		*contentID = 0
	}

	// Get first line of part
	contentStart = firstLineRegex.FindIndex(buffer)[1] + 1
	rawContent := buffer[contentStart:]

	if contentInfo.Type.Type == "multipart" {
		// Get contents starting of mail based on boundary
		startBoundaryRegex := regexp.MustCompile(`(?m)^--` + contentInfo.Type.Parameters["boundary"] + `\s*$`)
		contentStart = startBoundaryRegex.FindIndex(rawContent)[0]
		rawContent = rawContent[contentStart:]

		// Extract all contents as RAW format, multipart header and separation are ignored
		endBoundaryRegex := regexp.MustCompile(`(?m)^--` + contentInfo.Type.Parameters["boundary"] + `--\s*$`)
		contentEnd = endBoundaryRegex.FindIndex(rawContent)[0]
		rawContent = rawContent[:contentEnd]

		// Get parts of content
		indexes = startBoundaryRegex.FindAllIndex(rawContent, -1)
	} else {
		indexes = [][]int{{0, -1}}
		contentEnd = len(rawContent)
	}

	for i := range indexes {
		start = indexes[i][1] + 1
		if i < len(indexes)-1 {
			end = indexes[i+1][0]
		} else {
			end = contentEnd
		}
		ci = ContentInfo{}
		c = []byte{}
		if contentInfo.Type.Type == "multipart" {
			ci = getContentInfo(rawContent[start:end])
			c = rawContent[start:end]
			nc, a := ParseContents(c, ci, contentID)
			attachment = a
			//contents[strconv.Itoa(*contentID)] = Content{ci, c, nc}
			for k, v := range nc {
				contents[k] = v
			}
			*contentID += 1
		} else {
			ci = contentInfo
			c = rawContent[start:end]
			if ci.Disposition.Type == "attachment" {
				attachment[ci.Disposition.Parameters["filename"]] = Attachment{ci, c}
			} else {
				contents[strconv.Itoa(*contentID)] = Content{ci, c, nil}
				*contentID += 1
			}
		}
	}
	return contents, attachment
}

func extractContent(buffer []byte, contentInfo ContentInfo) Content {
	return Content{}
}
