package pkg

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// Section represents a section in an INI file.
// Each section has a name and a collection of key-value pairs.
type Section struct {
	Name string
	Keys map[string]string
}

// IniParser is a parser for INI files.
// Stores all the sections and their corresponding key-value pairs.
type IniParser struct {
	Sections map[string]Section
}

// List of all section names.
func (ini *IniParser) GetSectionNames() []string {
	var sectionNames []string
	for sectionName := range ini.Sections {
		sectionNames = append(sectionNames, sectionName)
	}
	return sectionNames
}

// Serialize convert into a dictionary/map { section_name: {key1: val1, key2, val2} ...}.
func (ini *IniParser) GetSections() map[string]map[string]string {
	sections := make(map[string]map[string]string)
	for sectionName, section := range ini.Sections {
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
	section, exists := ini.Sections[sectionName]
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
	section, exists := ini.Sections[sectionName]
	if !exists {
		section = Section{
			Name: sectionName,
			Keys: make(map[string]string),
		}
		ini.Sections[sectionName] = section
	}
	section.Keys[key] = value
}

// Converts the ini file to string
func (ini *IniParser) ToString() string {
	var sb strings.Builder

	for sectionName, section := range ini.Sections {
		sb.WriteString("[" + sectionName + "]\n")
		for key, value := range section.Keys {
			sb.WriteString(key + " = " + value + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// Saves data to a new file
func (ini *IniParser) SaveToFile() error {
	f, err := os.Create("NewFile.ini")
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(ini.ToString())
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Load data from an input string
func LoadFromString(input string) (*IniParser, error) {
	ini := &IniParser{
		Sections: make(map[string]Section),
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
		parseIniLine(ini, &currentSection, line)
	}
	return ini, nil
}

// Loads data from ini file
func LoadFromFile(fileName string) (*IniParser, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ini := &IniParser{
		Sections: make(map[string]Section),
	}
	reader := bufio.NewReader(file)
	var currentSection string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		parseIniLine(ini, &currentSection, line)
	}
	return ini, nil
}

// parsing each line of an ini file or string
func parseIniLine(ini *IniParser, currentSection *string, line string) {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '#' || line[0] == ';' {
		return
	}
	if line[0] == '[' && line[len(line)-1] == ']' {
		*currentSection = line[1 : len(line)-1]
		ini.Sections[*currentSection] = Section{
			Name: *currentSection,
			Keys: make(map[string]string),
		}
	} else if *currentSection != "" {
		pairs := strings.SplitN(line, "=", 2)
		if len(pairs) < 2 {
			return
		}
		key := strings.TrimSpace(pairs[0])
		value := strings.TrimSpace(pairs[1])
		section := ini.Sections[*currentSection]
		section.Keys[key] = value
		ini.Sections[*currentSection] = section
	}
}
