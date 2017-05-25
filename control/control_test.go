package control_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"git.yale.edu/docker/deco/control"
)

var testDecoString string = `
{
    "filters": {
        "test/file1": {
            "string1": "value1",
            "string2": "value2",
            "string3": "value3"
        },
        "test/file2": {
            "string1": "othervalue1"
        }
    }
}
`

var file1 = map[string]string{
	"string1": "value1",
	"string2": "value2",
	"string3": "value3",
}

var file2 = map[string]string{
	"string1": "othervalue1",
}

var testDecoStruct = control.Configuration{
	Filters: map[string]map[string]string{
		"test/file1": file1,
		"test/file2": file2,
	},
}

func TestRead(t *testing.T) {
	testFile := createTemporaryConfigFile()
	defer os.Remove(testFile.Name())

	filename := testFile.Name()
	var actual control.Configuration
	actual.Read(filename)

	for filterFile, filterMap := range testDecoStruct.Filters {
		if actualFilterMap := actual.Filters[filterFile]; actualFilterMap != nil {
			for find, replace := range filterMap {
				if actualFilterMap[find] != replace {
					t.Errorf("control.Read(%s) for key '%s', got replacement '%s', expected '%s'", filename, find, actualFilterMap[find], replace)
				}
			}
		} else {
			t.Errorf("control.Read(%s) returned nil filter map for file.", filename)
		}
	}
}

func createTemporaryConfigFile() *os.File {
	var tmpfile *os.File
	content := []byte(testDecoString)
	tmpfile, err := ioutil.TempFile("", "decotest")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile
}
