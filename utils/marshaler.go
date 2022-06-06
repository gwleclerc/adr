package utils

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func MarshalYAML(v any) ([]byte, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2) // this is what you're looking for
	err := yamlEncoder.Encode(v)
	return b.Bytes(), err
}
