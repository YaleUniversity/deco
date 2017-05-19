package control

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"
)

// Configuration is the overall data structure unmarshalled from JSON
type Configuration struct {
	Filters map[string]map[string]string
}

// Read reads in the configuration and returns the object
func (c *Configuration) Read(file string) error {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("[ERROR] unable to read file!", err.Error())
		return err
	}

	if err := json.Unmarshal(raw, c); err != nil {
		return err
	}

	return nil
}

// Print displays the configuration object
func (c *Configuration) Print() {
	fmt.Printf("%+v", c)
}

// DoFilters filters the files listed in the Configuration object
func (c *Configuration) DoFilters() error {
	for f, filters := range c.Filters {
		fmt.Println("Filtering", f)
		if err := Filter(f, filters); err != nil {
			fmt.Println("Error filtering template", err)
			return err
		}
	}

	return nil
}

// Filter filters an individual file
func Filter(file string, filters map[string]string) error {

	funcMap := template.FuncMap{
		"b64dec": base64decode,
		"b64enc": base64encode,
	}

	blob, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Unable to read file", err)
		return err
	}

	tmpl, err := template.New("config").Funcs(funcMap).Parse(string(blob))
	if err != nil {
		fmt.Println("Unable to parse template file", err)
		return err
	}

	var b bytes.Buffer
	parsedTemplate := bufio.NewWriter(&b)
	err = tmpl.Execute(parsedTemplate, filters)
	if err != nil {
		fmt.Println("Unable to execute parsed template", err)
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
		fmt.Println("Error decoding base64 encoded string", err)
		return err.Error()
	}
	return string(data)
}

// base64encode base64 encodes a string
func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}
