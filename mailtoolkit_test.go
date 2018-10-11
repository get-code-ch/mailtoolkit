package mailtoolkit

import (
	"io/ioutil"
	"testing"
)

const checkMark = "\u2714"
const ballotX = "\u2718"
const cross = "\u271A"

const file1 = "./files/test_multipart_complex.eml"
const file2 = "./files/test_nonmime.eml"

func TestParse(t *testing.T) {
	buffer1, err := ioutil.ReadFile(file1)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Logf("Testing mail parse %s", file1)
	{
		mail := Parse(buffer1)
		for _, e := range mail.Contents {
			t.Logf("First part of mail content:\n%v", string(e.Data))
		}
	}
	buffer2, err := ioutil.ReadFile(file2)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Logf("Testing mail parse %s", file2)
	{
		mail := Parse(buffer2)
		t.Logf("RAW mail content:\n%v", string(mail.Contents[0].Data))
	}
}
