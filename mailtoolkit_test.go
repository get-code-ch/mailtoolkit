package mailtoolkit

import (
	"io/ioutil"
	"testing"
)

const checkMark = "\u2714"
const ballotX = "\u2718"
const cross = "\u271A"

const file1 = "./files/test1.eml"

func TestParseHeader(t *testing.T) {
	buffer, err := ioutil.ReadFile(file1)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Log("Testing mail header")
	{
		contentType := "multipart"
		boundary := "----=_NextPart_000_006D_01D4415F.8115DFE0"
		header := ParseHeader(buffer)
		if header.ContentInfo.Type.Type != contentType {
			t.Fatalf("\t%v Error wrong Content-Type, should be \"%s\"", ballotX, contentType)
		}
		if header.From == "" {
			t.Errorf("\t%v Missing From value", ballotX)
		}
		if header.To == "" {
			t.Errorf("\t%v Missing To value", ballotX)
		}
		if header.Subject == "" {
			t.Logf("\t%v Subject is blank", cross)
		}
		_, exist := header.Elements["received"]
		if !exist {
			t.Errorf("\t%v Missing received field", ballotX)
		}
		b, exist := header.ContentInfo.Type.Parameters["boundary"]
		if !exist {
			t.Errorf("\t%v Missing boundary parameter", ballotX)
		} else {
			if b != boundary {
				t.Errorf("\t%v Wrong boundary value, schould be \"%s\"", ballotX, boundary)
			}
		}
	}
}

func TestParseBody(t *testing.T) {
	buffer, err := ioutil.ReadFile(file1)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Log("Testing mail content")
	{
		_ = buffer
	}

}
