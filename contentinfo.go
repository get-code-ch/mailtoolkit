package mailtoolkit

import (
	"strings"
)

func getContentInfo(buffer []byte) ContentInfo {
	var match [][]byte
	contentInfo := ContentInfo{}
	contentInfo.Type.Parameters = make(map[string]string)
	contentInfo.Disposition.Parameters = make(map[string]string)

	// Get End of Header (blank line)
	end := firstLineRegex.FindIndex(buffer)[0]

	// Get Attachments-Type
	contentInfo.Type = getContentType(buffer[:end])

	// Get Attachments-Disposition
	contentInfo.Disposition = getContentDisposition(buffer[:end])

	// Get Attachments-Transfer-Encoding Attachments-Transfer-Encoding:
	match = contentTransferEncodingRegex.FindSubmatch(buffer[:end])
	if match != nil {
		contentInfo.TransferEncoding = removeQuotesRegex.ReplaceAllString(string(match[1]), ``)
	}

	// Get Attachments-ID
	match = contentIDRegex.FindSubmatch(buffer[:end])
	if match != nil {
		contentInfo.ID = removeQuotesRegex.ReplaceAllString(string(match[1]), ``)
	}
	// Get Attachments-Description
	match = contentDescriptionRegex.FindSubmatch(buffer[:end])
	if match != nil {
		contentInfo.Description = removeQuotesRegex.ReplaceAllString(string(match[1]), ``)
	}

	return contentInfo
}

func getContentType(buffer []byte) ContentType {
	contentType := ContentType{}
	contentType.Parameters = make(map[string]string)
	wrkCT := "" // working content-type string

	// Find Attachments-Type
	wrkBuffer := contentTypeRegex.FindSubmatch(buffer)
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
	wrkCT = whitespaceRegex.ReplaceAllString(wrkCT, ``)

	//Extract parameters to a slice
	parameters := semiColonRegex.Split(wrkCT, -1)

	// Get Attachments Type and Subtype
	se := slashRegex.Split(parameters[0], -1)
	contentType.Type = se[0]
	contentType.Subtype = se[1]

	// Get Parameters
	for _, param := range parameters[1:] {
		se = parametersRegex.FindStringSubmatch(param)
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
	wrkBuffer := contentDispositionRegex.FindSubmatch(buffer)
	// If Attachments-Type is not found we return empty structure
	if wrkBuffer == nil {
		return contentDisposition
	}
	// Concatenate parameters to one string and removing whitespaces
	for _, wb := range wrkBuffer[1:] {
		wrkCD += string(wb)
	}
	wrkCD = whitespaceRegex.ReplaceAllString(wrkCD, ``)

	//Extract parameters to a slice
	parameters := semiColonRegex.Split(wrkCD, -1)

	// Get Attachments Type and Subtype
	se := slashRegex.Split(parameters[0], -1)
	contentDisposition.Type = se[0]

	// Get Parameters
	for _, param := range parameters[1:] {
		se = parametersRegex.FindStringSubmatch(param)
		if se != nil {
			// Parameter attribute are normalized to lowercase
			contentDisposition.Parameters[strings.ToLower(se[1])] = se[2]
		}
	}
	return contentDisposition
}
