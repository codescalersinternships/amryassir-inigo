package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Section represents a section in an INI file.
// Each section has a name and a collection of key-value pairs.
type Section struct {
	Keys map[string]string
}

// IniParser is a parser for INI files.
// Stores all the sections and their corresponding key-value pairs.
type IniParser struct {
	sections map[string]Section
}

// List of all section names.
func (ini *IniParser) GetSectionNames() []string {
	var sectionNames []string
	for sectionName := range ini.sections {
		sectionNames = append(sectionNames, sectionName)
	}
	return sectionNames
}

// Serialize convert into a dictionary/map { section_name: {key1: val1, key2, val2} ...}.
func (ini *IniParser) GetSections() map[string]map[string]string {
	sections := make(map[string]map[string]string)
	for sectionName, section := range ini.sections {
		sectionMap := make(map[string]string)
		for key, value := range section.Keys {
			sectionMap[key] = value
		}
		sections[sectionName] = sectionMap
	}
	return sections
}

// Gets the value of key key in section section_name.
func (ini *IniParser) Get(sectionName, key string) (string, error) {
	section, exists := ini.sections[sectionName]
	if !exists {
		return "", errors.New("section does not exist")
	}

	value, exists := section.Keys[key]
	if !exists {
		return "", errors.New("key does not exist in section")
	}

	return value, nil
}

// Sets a key in section section_name to value value.
func (ini *IniParser) Set(sectionName, key, value string) {
	section, exists := ini.sections[sectionName]
	if !exists {
		section = Section{
			Keys: make(map[string]string),
		}
		ini.sections[sectionName] = section
	}
	section.Keys[key] = value
	ini.sections[sectionName] = section
}

// Converts the ini file to string
func (ini *IniParser) String() string {
	var result strings.Builder
	for sectionName, section := range ini.sections {
		result.WriteString(fmt.Sprintf("[%v]\n", sectionName))
		for key, value := range section.Keys {
			result.WriteString(fmt.Sprintf("%v = %v\n", key, value))
		}
		result.WriteString("\n")
	}
	return result.String()
}

// Saves data to a new file
func (ini *IniParser) SaveToFile() error {
	d1 := []byte(ini.String())
	err := os.WriteFile("NewFile.ini", d1, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Load data from an input string
func LoadFromString(input string) (*IniParser, error) {
	ini := &IniParser{
		sections: make(map[string]Section),
	}
	reader := bufio.NewReader(strings.NewReader(input))
	var currentSection string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if err := ini.parseIniLine(&currentSection, line); err != nil {
			return nil, err
		}
	}
	return ini, nil
}

// Loads data from ini file
func LoadFromFile(fileName string) (*IniParser, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return LoadFromString(string(content))
}

// parsing each line of an ini file or string
func (ini *IniParser) parseIniLine(currentSection *string, line string) error {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '#' || line[0] == ';' {
		return nil
	}
	if line[0] == '[' && line[len(line)-1] == ']' {
		*currentSection = line[1 : len(line)-1]
		ini.sections[*currentSection] = Section{
			Keys: make(map[string]string),
		}
	} else if *currentSection != "" {
		pairs := strings.SplitN(line, "=", 2)
		if len(pairs) < 2 {
			return fmt.Errorf("invalid key-value pair: %s", line)
		}
		key := strings.TrimSpace(pairs[0])
		value := strings.TrimSpace(pairs[1])
		section := ini.sections[*currentSection]
		section.Keys[key] = value
		ini.sections[*currentSection] = section
	} else {
		return fmt.Errorf("line does not belong to any section: %s", line)
	}
	return nil
}
