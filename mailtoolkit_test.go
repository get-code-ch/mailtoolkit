package mailtoolkit

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

const checkMark = "\u2714"
const ballotX = "\u2718"
const cross = "\u271A"

const testFilesFolder = "./files/"

func TestParse(t *testing.T) {
	buffer, err := ioutil.ReadFile("mailtoolkit_test.json")
	if err != nil {
		t.Fatal("Error opening test directive:", err)
	}
	var directives []map[string]interface{}
	json.Unmarshal([]byte(buffer), &directives)

	for idx := range directives {
		directive := directives[idx]
		filename := directive["filename"].(string)
		buffer, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal("Error opening test file:", err)
		}
		t.Logf("Testing mail parse %s", filename)

		mail := Parse(buffer)
		for key, content := range mail.Contents {
			if content.Data != nil {
				l := len(content.Data)
				if l > 200 {
					l = 200
				}
				t.Logf("%s part of mail content:\n%v", key, string(content.Data[:l]))
			}
		}
		header := directive["header"].(map[string]interface{})
		for key, value := range header {
			fieldValue := reflect.ValueOf(mail.Header).FieldByName(key).String()
			if fieldValue != value {
				t.Errorf("\t%v Error wrong value for %s, result is \"%v\" should be \"%v\"", ballotX, key, fieldValue, value)
			} else {
				t.Logf("\t%v Ok value for %s, match with \"%v\"", checkMark, key, value)
			}
		}
	}
}
