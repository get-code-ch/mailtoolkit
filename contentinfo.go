package mailtoolkit

import (
	"mime"
	"net/textproto"
	"strings"
)

func newContentInfo(h textproto.MIMEHeader) ContentInfo {
	return ContentInfo{
		Type:             parseContentType(h.Get("Content-Type")),
		Disposition:      parseContentDisposition(h.Get("Content-Disposition")),
		TransferEncoding: strings.ToLower(strings.TrimSpace(h.Get("Content-Transfer-Encoding"))),
		ID:               strings.Trim(strings.TrimSpace(h.Get("Content-ID")), "<>"),
		Description:      decodeHeader(h.Get("Content-Description")),
	}
}

// parseContentType defaults to text/plain (RFC 2045 §5.2) when the field is
// missing or invalid.
func parseContentType(value string) ContentType {
	contentType := ContentType{Type: "text", Subtype: "plain", Parameters: map[string]string{}}

	mediaType, params := parseMediaType(value)
	main, sub, ok := strings.Cut(mediaType, "/")
	if !ok || main == "" || sub == "" {
		return contentType
	}
	contentType.Type, contentType.Subtype, contentType.Parameters = main, sub, params
	return contentType
}

func parseContentDisposition(value string) ContentDisposition {
	disposition, params := parseMediaType(value)
	return ContentDisposition{Type: disposition, Parameters: params}
}

// parseMediaType wraps mime.ParseMediaType, keeping the media type when only
// the parameters are invalid, and decodes RFC 2047 encoded file names as sent
// by some mail clients.
func parseMediaType(value string) (string, map[string]string) {
	params := map[string]string{}
	if strings.TrimSpace(value) == "" {
		return "", params
	}

	mediaType, parsed, err := mime.ParseMediaType(value)
	if err != nil && err != mime.ErrInvalidMediaParameter {
		return "", params
	}
	for key, param := range parsed {
		if key == "filename" || key == "name" {
			param = decodeHeader(param)
		}
		params[key] = param
	}
	return mediaType, params
}
