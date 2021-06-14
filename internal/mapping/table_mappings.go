package mapping

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ColumnMapping map[string]*string

// LoadMappings loads all column mappings for multiple tables from a YAML file
func LoadMappings(mappingFile string) (map[string]ColumnMapping, error) {
	f, err := os.Open(mappingFile)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	allMappings := make(map[string]ColumnMapping, 0)
	err = yaml.Unmarshal(contents, &allMappings)
	if err != nil {
		return nil, err
	}
	return allMappings, nil
}

// SaveMappingToFile outputs the mapping for a table to a YAML file
func SaveMappingToFile(fileName, tableName string, mapped ColumnMapping, dry bool) error {
	existing, err := LoadMappings(fileName)
	if err != nil && errors.Unwrap(err).Error() != "no such file or directory" {
		return err
	}
	if existing == nil {
		existing = make(map[string]ColumnMapping, 0)
	}
	existing[tableName] = mapped
	yamlOut, err := yaml.Marshal(existing)
	if err != nil {
		return err
	}
	if !dry {
		err = ioutil.WriteFile(fileName, yamlOut, 644)
		if err != nil {
			return err
		}
	}
	return nil
}
