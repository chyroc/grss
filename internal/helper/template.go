package helper

import (
	"bytes"
	"text/template"
)

func BuildTemplate(templateString string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	t, err := template.New("").Parse(templateString)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
