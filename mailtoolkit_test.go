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
const file3 = "./files/test_multipart.eml"

func TestParse(t *testing.T) {

	// Test complex (nested) multipart mail
	buffer1, err := ioutil.ReadFile(file1)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Logf("Testing mail parse %s", file1)
	{
		mail := Parse(buffer1)
		for key, content := range mail.Contents {
			if content.Data != nil {
				t.Logf("%s part of mail content:\n%v", key, string(content.Data[:200]))
			}
		}
	}
	// Test non MIME mail
	buffer2, err := ioutil.ReadFile(file2)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Logf("Testing mail parse %s", file2)
	{
		mail := Parse(buffer2)
		for key, content := range mail.Contents {
			if content.Data != nil {
				t.Logf("%s part of mail content:\n%v", key, string(content.Data[:200]))
			}
		}
	}

	buffer3, err := ioutil.ReadFile(file3)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Logf("Testing mail parse %s", file3)
	{
		mail := Parse(buffer3)
		for key, content := range mail.Contents {
			if content.Data != nil {
				t.Logf("%s part of mail content:\n%v", key, string(content.Data[:200]))
			}
		}
	}
}
