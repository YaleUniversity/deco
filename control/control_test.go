package control_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/YaleUniversity/deco/control"
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

func TestReadFile(t *testing.T) {
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

func TestReadURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testDecoString)
	}))
	defer ts.Close()

	var actual control.Configuration
	err := actual.Read(ts.URL)
	if err != nil {
		t.Errorf("Expected to successfully read for test URL")
	}

	t.Log("Got control file from test URL:", actual)

	for filterFile, filterMap := range testDecoStruct.Filters {
		if actualFilterMap := actual.Filters[filterFile]; actualFilterMap != nil {
			for find, replace := range filterMap {
				if actualFilterMap[find] != replace {
					t.Errorf("control.Read(%s) for key '%s', got replacement '%s', expected '%s'", ts.URL, find, actualFilterMap[find], replace)
				}
			}
		} else {
			t.Errorf("control.Read(%s) returned nil filter map for URL.", ts.URL)
		}
	}
}
