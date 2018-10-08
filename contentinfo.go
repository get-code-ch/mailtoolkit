package mailtoolkit

import (
	"regexp"
	"strings"
)

func getContentInfo(buffer []byte) ContentInfo {
	var match [][]byte
	contentInfo := ContentInfo{}
	contentInfo.Type.Parameters = make(map[string]string)
	contentInfo.Disposition.Parameters = make(map[string]string)

	// Get End of Header (blank line)
	re := regexp.MustCompile(`(?m)(^[\n|\n\r]?$)`)
	end := re.FindIndex(buffer)[0]

	// Get Attachments-Type
	contentInfo.Type = getContentType(buffer[:end])

	// Get Attachments-Disposition
	contentInfo.Disposition = getContentDisposition(buffer[:end])

	// Get Attachments-Transfer-Encoding Attachments-Transfer-Encoding:
	re = regexp.MustCompile(`(?mi)(?:^\s*Content-Transfer-Encoding:\s+"?)(.*)(?:"?\n?)`)
	match = re.FindSubmatch(buffer[:end])
	if match != nil {
		contentInfo.TransferEncoding = string(match[1])
	}

	// Get Attachments-ID
	re = regexp.MustCompile(`(?mi)(?:^\s*Content-ID:\s+"?)(.*)(?:"?\n?)`)
	match = re.FindSubmatch(buffer[:end])
	if match != nil {
		contentInfo.ID = string(match[1])
	}
	// Get Attachments-Description
	re = regexp.MustCompile(`(?mi)(?:^\s*Content-Description:\s+"?)(.*)(?:"?\n?)`)
	match = re.FindSubmatch(buffer[:end])
	if match != nil {
		contentInfo.Description = string(match[1])
	}

	return contentInfo
}

func getContentType(buffer []byte) ContentType {
	contentType := ContentType{}
	contentType.Parameters = make(map[string]string)
	wrkCT := "" // working content-type string

	// Find Attachments-Type
	re := regexp.MustCompile(`(?im)(?:^Content-Type: ?)(.+)(?:\r?\n)((?:\s*(?:\s+).*(?:\r?\n+))*)`)
	wrkBuffer := re.FindSubmatch(buffer)
	// If Attachments-Type is not found we assume content type is Text/plain (non MIME email) rest of datas are nil
	if wrkBuffer == nil {
		contentType.Type = "text"
		contentType.Subtype = "plain"
		return contentType
	}
	// Concatenate parameters to one string and removing whitespaces
	for _, wb := range wrkBuffer[1:] {
		wrkCT += string(wb)
	}
	re = regexp.MustCompile(`[ ]{2,}|[\t|\0|\n|\r]+`)
	wrkCT = re.ReplaceAllString(wrkCT, ``)

	//Extract parameters to a slice
	re = regexp.MustCompile(`;`)
	parameters := re.Split(wrkCT, -1)

	// Get Attachments Type and Subtype
	re = regexp.MustCompile(`/`)
	se := re.Split(parameters[0], -1)
	contentType.Type = se[0]
	contentType.Subtype = se[1]

	// Get Parameters
	for _, param := range parameters[1:] {
		re = regexp.MustCompile(`(?i)^(?:\s*)([\w-]+)(?:\s*=\s*"?)(.*[^"])(?:"?\s*)$`)
		se = re.FindStringSubmatch(param)
		if se != nil {
			// Parameter attribute are normalized to lowercase
			contentType.Parameters[strings.ToLower(se[1])] = se[2]
		}
	}
	return contentType
}

func getContentDisposition(buffer []byte) ContentDisposition {
	contentDisposition := ContentDisposition{}
	contentDisposition.Parameters = make(map[string]string)
	wrkCD := "" // working content-disposition string

	// Find Attachments-Type
	re := regexp.MustCompile(`(?im)(?:^Content-Disposition: ?)(.+)(?:\r?\n)((?:\s*(?:\s+).*(?:\r?\n+))*)`)
	wrkBuffer := re.FindSubmatch(buffer)
	// If Attachments-Type is not found we return empty structure
	if wrkBuffer == nil {
		return contentDisposition
	}
	// Concatenate parameters to one string and removing whitespaces
	for _, wb := range wrkBuffer[1:] {
		wrkCD += string(wb)
	}
	re = regexp.MustCompile(`[ ]{2,}|[\t|\0|\n|\r]+`)
	wrkCD = re.ReplaceAllString(wrkCD, ``)

	//Extract parameters to a slice
	re = regexp.MustCompile(`;`)
	parameters := re.Split(wrkCD, -1)

	// Get Attachments Type and Subtype
	re = regexp.MustCompile(`/`)
	se := re.Split(parameters[0], -1)
	contentDisposition.Type = se[0]

	// Get Parameters
	for _, param := range parameters[1:] {
		re = regexp.MustCompile(`(?i)^(?:\s*)([\w-]+)(?:\s*=\s*"?)(.*[^"])(?:"?\s*)$`)
		se = re.FindStringSubmatch(param)
		if se != nil {
			// Parameter attribute are normalized to lowercase
			contentDisposition.Parameters[strings.ToLower(se[1])] = se[2]
		}
	}
	return contentDisposition
}
