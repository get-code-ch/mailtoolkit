package mailtoolkit

import (
	"mime"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"
)

var wordDecoder = &mime.WordDecoder{}

// emailRegex is only a fallback for address fields net/mail cannot parse.
var emailRegex = regexp.MustCompile(`[^\s<>"',;:()@]+@[^\s<>"',;:()@]+`)

func newHeader(h mail.Header) Header {
	header := Header{Elements: make(map[string]string, len(h))}

	for key, values := range h {
		header.Elements[strings.ToLower(key)] = strings.Join(values, "\n")
	}

	_, header.IsMime = h["Mime-Version"]
	header.ContentInfo = newContentInfo(textproto.MIMEHeader(h))

	header.FromList = parseAddresses(h.Get("From"))
	if _, ok := h["To"]; ok {
		header.ToList = parseAddresses(h.Get("To"))
	} else {
		header.ToList = parseAddresses(h.Get("Delivered-To"))
	}
	header.CcList = parseAddresses(h.Get("Cc"))
	header.BccList = parseAddresses(h.Get("Bcc"))

	header.From = joinAddresses(header.FromList)
	header.To = joinAddresses(header.ToList)
	header.Cc = joinAddresses(header.CcList)
	header.Bcc = joinAddresses(header.BccList)

	header.Subject = decodeHeader(h.Get("Subject"))
	header.Date = h.Get("Date")
	if t, err := mail.ParseDate(header.Date); err == nil {
		header.Time = t
	}

	return header
}

// decodeHeader decodes RFC 2047 encoded-words, returning value unchanged if
// it cannot be decoded.
func decodeHeader(value string) string {
	decoded, err := wordDecoder.DecodeHeader(value)
	if err != nil {
		return value
	}
	return decoded
}

func parseAddresses(value string) []*mail.Address {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parser := mail.AddressParser{WordDecoder: wordDecoder}
	if list, err := parser.ParseList(value); err == nil {
		return list
	}
	var list []*mail.Address
	for _, address := range emailRegex.FindAllString(value, -1) {
		list = append(list, &mail.Address{Address: address})
	}
	return list
}

func joinAddresses(list []*mail.Address) string {
	addresses := make([]string, len(list))
	for i, a := range list {
		addresses[i] = a.Address
	}
	return strings.Join(addresses, ", ")
}
