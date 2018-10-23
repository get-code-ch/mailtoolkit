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

const directiveFile = "mailtoolkit_test.json"

func TestParse(t *testing.T) {
	var displayContent bool

	// Import test directive from JSON files
	buffer, err := ioutil.ReadFile(directiveFile)
	if err != nil {
		t.Fatal("Error opening test directive:", err)
	}

	// Parse JSON to directives object array
	var directives []map[string]interface{}
	json.Unmarshal([]byte(buffer), &directives)

	// Browse array of directives an run test
	for idx := range directives {
		directive := directives[idx]
		filename := directive["filename"].(string)

		// open email file and parse content (call functions to be tested)
		buffer, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal("Error opening test file:", err)
		}
		t.Logf("======= Testing mail parse %s =======", filename)
		mail := Parse(buffer)

		// Check result
		d := directive["displaycontent"]
		if d == nil {
			displayContent = false
		} else {
			displayContent = d.(bool)
		}
		if displayContent {
			for key, content := range mail.Contents {
				if content.Data != nil {
					l := len(content.Data)
					if l > 200 {
						l = 200
					}
					t.Logf("content part %s:\n%v", key, string(content.Data[:l]))
				}
			}
		}

		// Check returned header of mail file
		header := directive["header"].(map[string]interface{})
		for key, value := range header {
			if key == "ContentInfo" {
				contentInfo := header[key].(map[string]interface{})["Type"].(map[string]interface{})
				contentType := contentInfo["Type"]
				contentSubtype := contentInfo["Subtype"]
				mailType := reflect.ValueOf(mail.Header.ContentInfo.Type.Type).String()
				mailSubtype := reflect.ValueOf(mail.Header.ContentInfo.Type.Subtype).String()
				if contentType != mailType || contentSubtype != mailSubtype {
					t.Errorf("\t%v Error wrong value for type/subtype, result is \"%s/%s\" should be \"%s/%s\"", ballotX, mailType, mailSubtype, contentType, contentSubtype)
				} else {
					t.Logf("\t%v Ok value for type/subtype, match with \"%s/%s\"", checkMark, mailType, mailSubtype)
				}
				continue
			}
			// Check fields of Header
			fieldValue := reflect.ValueOf(mail.Header).FieldByName(key).String()
			if fieldValue != value {
				t.Errorf("\t%v Error wrong value for %s, result is \"%v\" should be \"%v\"", ballotX, key, fieldValue, value)
			} else {
				t.Logf("\t%v Ok value for %s, match with \"%v\"", checkMark, key, value)
			}
		}
	}
}
