package mailtoolkit

import (
	"regexp"
	"strings"
)

func ParseHeader(buffer *[]byte) Header {
	var re *regexp.Regexp
	var header Header

	// Get End of Header (blank line)
	re = regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
	end := re.FindIndex(*buffer)

	// Get Header Elements
	header.Elements = make(map[string]string)
	re = regexp.MustCompile(`(?mi)(^[\w_-]+)(?::\s+"?)(.*)(?:\r?\n)((?:\s*(?:\s+).*(?:\r?\n+))*)`)
	elements := re.FindAllSubmatch((*buffer)[:end[0]], -1)
	for _, element := range elements {
		value := ""
		for _, fieldValue := range element[2:] {
			value += string(fieldValue)
		}
		// cleaning string from whitespaces and newline....
		re = regexp.MustCompile(`[ ]{2,}|[\t|\0|\n|\r]+`)
		value = re.ReplaceAllString(value, ``)

		key := strings.ToLower(string(element[1]))
		_, exist := header.Elements[key]
		if exist {
			header.Elements[key] += "\n" + value
		} else {
			header.Elements[key] = value
		}
	}

	_, header.IsMime = header.Elements["mime-version"]
	header.ContentInfo = getContentInfo((*buffer)[:end[0]])

	e, ok := header.Elements["from"]
	if ok {
		header.From = e
	}

	e, ok = header.Elements["delivered-to"]
	if ok {
		header.To = e
	}

	e, ok = header.Elements["cc"]
	if ok {
		header.Cc = e
	}

	e, ok = header.Elements["cci"]
	if ok {
		header.Cci = e
	}

	e, ok = header.Elements["subject"]
	if ok {
		header.Subject = e
	}

	e, ok = header.Elements["date"]
	if ok {
		header.Date = e
	}

	return header
}
