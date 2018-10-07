package mailtoolkit

import (
	"regexp"
	"strconv"
)

type delimiter struct {
	start int
	end   int
}

func ParseContent(buffer *[]byte, header Header) map[string]Content {
	var parts []delimiter

	content := make(map[string]Content)

	// Get first line of part
	re := regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
	contentStart := re.FindIndex(*buffer)[1] + 1

	// Allways return a RAW mail content
	content["raw"] = Content{Data: (*buffer)[contentStart:], ContentInfo: header.ContentInfo}

	if !header.IsMime || header.ContentInfo.Type.Type != "multipart" {
		return content
	}

	// Get parts of mail based on boundary
	re = regexp.MustCompile(`(?m)--` + header.ContentInfo.Type.Parameters["boundary"] + `\s*$`)
	partStarts := re.FindAllIndex((*buffer)[contentStart:], -1)
	contentStart = contentStart + partStarts[0][1] + 1

	// Extract all content as RAW format multipart separation is ignored
	re = regexp.MustCompile(`(?m)--` + header.ContentInfo.Type.Parameters["boundary"] + `--\s*$`)
	contentEnd := re.FindIndex((*buffer)[contentStart:])[0]
	content["raw"] = Content{Data: (*buffer)[contentStart : contentStart+contentEnd]}

	for i := range partStarts {
		delimiter := delimiter{start: contentStart + partStarts[i][1] + 1}
		parts = append(parts, delimiter)
		if i < len(partStarts)-1 {
			parts[i].end = contentStart + partStarts[i+1][0]
		} else {
			parts[i].end = contentStart + contentEnd
		}

		content[strconv.Itoa(i)] = Content{Data: (*buffer)[parts[i].start:parts[i].end], ContentInfo: getContentInfo((*buffer)[parts[i].start:parts[i].end])}
	}

	return content
}
