package ini

import (
	"os"
	"reflect"
	"testing"
)

const StringInput = `[DEFAULT]
ServerAliveInterval = 45
Compression = yes

[forge.example]
User = hg

`

func TestGetSectionNames(t *testing.T) {
	ini, _ := LoadFromString(StringInput)
	got := ini.GetSectionNames()
	want := []string {"DEFAULT", "forge.example"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %q want %q", got, want)
	}
}

func TestGetSections(t *testing.T) {
	iniData, _ := LoadFromString(StringInput)
	got := iniData.GetSections()
	want := map[string]map[string]string{
		"DEFAULT": {
			"ServerAliveInterval": "45",
			"Compression": "yes",
		},
		"forge.example": {
			"User": "hg",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGet(t *testing.T) {
	iniData, _ := LoadFromString(StringInput)
	got, err := iniData.Get("DEFAULT", "ServerAliveInterval")
	want := "45"
	if err != nil {
		t.Errorf("Returned error: %v", err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSet(t *testing.T) {
	iniData, _ := LoadFromString(StringInput)
	
	// Seting in Existing section 
	iniData.Set("DEFAULT", "key1", "value1")
	want := "value1"
	got := iniData.Sections["DEFAULT"].Keys["key1"]
	if got != want {
		t.Errorf(" got %q, want %q", got, want)
	}

	// Setting in a new section
	iniData.Set("NewSection", "key2", "value2")
	want = "value2"
	got = iniData.Sections["NewSection"].Keys["key2"]
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	// Overwriting an existing key-value pair
	iniData.Set("DEFAULT", "Compression", "no")
	want = "no"
	got = iniData.Sections["DEFAULT"].Keys["Compression"]
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestToString(t *testing.T) {
	iniData, _ := LoadFromString(StringInput)
	got := iniData.ToString()
	want := StringInput

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSaveToFile(t *testing.T) {
	iniData, _ := LoadFromString(StringInput)
	err := iniData.SaveToFile()
	if err != nil {
		t.Errorf("Can't save data to a file")
	}

	got, err := os.ReadFile("NewFile.ini")
	if err != nil {
		t.Fatalf("Error reading saved file: %v", err)
	}

	want := []byte("[DEFAULT]\nServerAliveInterval = 45\nCompression = yes\n\n[forge.example]\nUser = hg\n\n")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Saved file content does not match expected content.\ngot %q want %q", got, want)
	}
	os.Remove("NewFile.ini")
}

func TestLoadFromString(t *testing.T) {
	_, err := LoadFromString(StringInput)
	if err != nil {
		t.Fatalf("Failed to load data from input: %v", err)
	}
}

func TestLoadFromFile(t *testing.T) {
	_, err := LoadFromFile("testdata/input.ini")
	if err != nil {
		t.Fatalf("Failed to load data from file: %v", err)
	}
}