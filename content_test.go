package mailtoolkit

import (
	"io/ioutil"
	"testing"
)

func TestParseContent(t *testing.T) {
	buffer, err := ioutil.ReadFile(file1)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	t.Log("Testing mail content")
	{
		_ = buffer
	}

}
