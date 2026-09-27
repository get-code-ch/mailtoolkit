package mailtoolkit

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"path"
	"strconv"
	"strings"
)

// Decode returns the content data with its Content-Transfer-Encoding removed.
func (c Content) Decode() ([]byte, error) {
	return decode(c.ContentInfo.TransferEncoding, c.Data)
}

// Decode returns the attachment data with its Content-Transfer-Encoding
// removed.
func (a Attachment) Decode() ([]byte, error) {
	return decode(a.ContentInfo.TransferEncoding, a.Data)
}

func decode(encoding string, data []byte) ([]byte, error) {
	switch strings.ToLower(encoding) {
	case "base64":
		return io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data)))
	case "quoted-printable":
		return io.ReadAll(quotedprintable.NewReader(bytes.NewReader(data)))
	default:
		return data, nil
	}
}

type bodyParser struct {
	contents    map[string]Content
	attachments map[string]Attachment
	nextID      int
}

func parseBody(body io.Reader, contentInfo ContentInfo) (map[string]Content, map[string]Attachment, error) {
	p := bodyParser{
		contents:    make(map[string]Content),
		attachments: make(map[string]Attachment),
	}
	err := p.walk(body, contentInfo)
	return p.contents, p.attachments, err
}

// walk reads a part, recursing into multipart ones. Leaf parts are kept even
// when they could only be partially read.
func (p *bodyParser) walk(r io.Reader, contentInfo ContentInfo) error {
	if contentInfo.Type.Type != "multipart" {
		data, err := io.ReadAll(r)
		p.add(contentInfo, data)
		if err != nil {
			return fmt.Errorf("mailtoolkit: reading %s/%s part: %w", contentInfo.Type.Type, contentInfo.Type.Subtype, err)
		}
		return nil
	}

	boundary := contentInfo.Type.Parameters["boundary"]
	if boundary == "" {
		return errors.New("mailtoolkit: multipart part without boundary")
	}
	reader := multipart.NewReader(r, boundary)
	for {
		// NextRawPart keeps the Content-Transfer-Encoding untouched, it is
		// removed on demand by Decode.
		part, err := reader.NextRawPart()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("mailtoolkit: reading multipart/%s: %w", contentInfo.Type.Subtype, err)
		}
		if err := p.walk(part, newContentInfo(part.Header)); err != nil {
			return err
		}
	}
}

func (p *bodyParser) add(contentInfo ContentInfo, data []byte) {
	if contentInfo.Disposition.Type == "attachment" {
		p.attachments[p.attachmentKey(contentInfo)] = Attachment{contentInfo, data}
		return
	}

	key := contentInfo.ID
	if _, exists := p.contents[key]; key == "" || exists {
		key = strconv.Itoa(p.nextID)
		p.nextID++
	}
	p.contents[key] = Content{contentInfo, data}
}

// attachmentKey returns the attachment file name, suffixed with " (n)" when
// the name is already used.
func (p *bodyParser) attachmentKey(contentInfo ContentInfo) string {
	name := contentInfo.Disposition.Parameters["filename"]
	if name == "" {
		name = contentInfo.Type.Parameters["name"]
	}
	if name == "" {
		name = "attachment"
	}

	key := name
	ext := path.Ext(name)
	base := strings.TrimSuffix(name, ext)
	for n := 2; ; n++ {
		if _, exists := p.attachments[key]; !exists {
			return key
		}
		key = fmt.Sprintf("%s (%d)%s", base, n, ext)
	}
}
