package control

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Logger is a STDERR logger
var Logger = log.New(os.Stderr, "", 0)

type Filter map[string]string

// Configuration is the overall data structure unmarshalled from JSON
type Configuration struct {
	Filters map[string]Filter
	BaseDir string
}

// Get fetches the control from a location and returns a io.ReadCloser
func Get(location string, headers []string) (io.ReadCloser, error) {
	u, err := url.ParseRequestURI(location)
	if err == nil {
		if u.Scheme == "http" || u.Scheme == "https" {
			Logger.Println("[INFO] Fetching control from URL", location)

			var client = &http.Client{
				Timeout: time.Second * 10,
			}

			req, err := http.NewRequest(http.MethodGet, location, nil)
			if err != nil {
				return nil, err
			}

			for _, h := range headers {
				header := strings.SplitN(h, "=", 2)
				if len(header) < 2 {
					e := fmt.Sprintf("Unable to parse HTTP header: %s", h)
					return nil, errors.New(e)
				}

				req.Header.Set(header[0], header[1])
			}

			res, err := client.Do(req)
			if err != nil {
				Logger.Println("[ERROR] Unable to get file from URL location", err)
				return nil, err
			}
			return res.Body, nil
		}

		if u.Scheme == "ssm" {
			Logger.Println("[INFO] Fetching control from SSM location", location)
			svc := NewSSM()
			res, err := svc.GetParameter(u.RequestURI())
			if err != nil {
				Logger.Println("[ERROR] Unable to get file from SSM", err)
				return nil, err
			}

			return ioutil.NopCloser(strings.NewReader(res)), nil
		}
	}

	Logger.Println("[INFO] Using control from file", location)
	r, err := os.Open(location)
	if err != nil {
		Logger.Println("[ERROR] unable to open control file!", err.Error())
		return nil, err
	}
	return r, nil
}

// Read reads in the configuration and returns the object
func (c *Configuration) Read(location string, headers []string, encoded bool) error {
	r, err := Get(location, headers)
	if err != nil {
		return err
	}
	defer r.Close()

	if encoded {
		return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, r)).Decode(c)
	}

	return json.NewDecoder(r).Decode(c)
}

// Print displays the configuration object
func (c *Configuration) Print() {
	Logger.Printf("%+v", c)
}

// DoFilters filters the files listed in the Configuration object
func (c *Configuration) DoFilters() error {
	for f, filters := range c.Filters {
		if c.BaseDir != "" {
			f = filepath.Join(c.BaseDir, f)
		}

		Logger.Println("Filtering", f)
		if err := ExecFilter(f, filters); err != nil {
			Logger.Println("[ERROR] Failed filtering template", err)
			return err
		}
	}

	return nil
}

// ExecFilter filters an individual file
func ExecFilter(file string, filters map[string]string) error {

	funcMap := template.FuncMap{
		"b64dec":  func(v string) string { return base64decode(v) },
		"b64enc":  func(v string) string { return base64encode(v) },
		"decrypt": func(v string) string { return decrypt(v) },
	}

	blob, err := ioutil.ReadFile(file)
	if err != nil {
		Logger.Println("[ERROR] Unable to read file", err)
		return err
	}

	tmpl, err := template.New("config").Funcs(funcMap).Parse(string(blob))
	if err != nil {
		Logger.Println("[ERROR] Unable to parse template file", err)
		return err
	}

	var b bytes.Buffer
	parsedTemplate := bufio.NewWriter(&b)
	err = tmpl.Execute(parsedTemplate, filters)
	if err != nil {
		Logger.Println("[ERROR] Unable to execute parsed template", err)
		return err
	}
	parsedTemplate.Flush()

	if err := ioutil.WriteFile(file, b.Bytes(), 0444); err != nil {
		return err
	}

	return nil
}

// base64decode decodes a base64 encoded string
func base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		Logger.Println("[ERROR] Failed decoding base64 encoded string", err)
		return err.Error()
	}
	return string(data)
}

// base64encode base64 encodes a string
func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func decrypt(v string) string {
	encryptionKey := os.Getenv("DECO_ENCRYPTION_KEY")
	if encryptionKey == "" {
		Logger.Println("[ERROR] Failed decrypt string, missing DECO_ENCRYPTION_KEY")
		return v
	}

	keyBytes, err := hex.DecodeString(encryptionKey)
	if err != nil {
		Logger.Println("[ERROR] Failed decrypt string, invalid key encoding", err)
		return v
	}

	var key [32]byte
	copy(key[:], keyBytes)

	cipherBytes, err := hex.DecodeString(v)
	if err != nil {
		Logger.Println("[ERROR] Failed decrypt string, invalid ciphertext encoding", err)
		return v
	}

	plainText, err := Decrypt(cipherBytes, &key)
	if err != nil {
		Logger.Println("[ERROR] Failed decrypt string", err)
		return v
	}

	return string(plainText)
}
