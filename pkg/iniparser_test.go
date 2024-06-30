package pkg

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// StringInput is a sample INI configuration string
const StringInput = `[DEFAULT]
ServerAliveInterval = 45
Compression = yes

[forge.example]
User = hg

`

func TestGetSectionNames(t *testing.T) {
	ini, _ := LoadFromString(StringInput)
	got := ini.GetSectionNames()
	want := []string{"DEFAULT", "forge.example"}

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
			"Compression":         "yes",
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
	t.Run("Setting in existing section", func(t *testing.T) {
		iniData.Set("DEFAULT", "key1", "value1")
		want := "value1"
		got := iniData.Sections["DEFAULT"].Keys["key1"]
		if got != want {
			t.Errorf(" got %q, want %q", got, want)
		}
	})
	t.Run("Setting in new section", func(t *testing.T) {
		iniData.Set("NewSection", "key2", "value2")
		want := "value2"
		got := iniData.Sections["NewSection"].Keys["key2"]
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("Overwriting existing key-value pair", func(t *testing.T) {
		iniData.Set("DEFAULT", "Compression", "no")
		want := "no"
		got := iniData.Sections["DEFAULT"].Keys["Compression"]
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestString(t *testing.T) {
	iniData, err := LoadFromString(StringInput)
	if err != nil {
		t.Fatalf("Failed to load data from input: %v", err)
	}
	got := iniData.String()
	want := StringInput

	got = strings.TrimSpace(got)
	want = strings.TrimSpace(want)

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
	got, err := LoadFromString(StringInput)
	if err != nil {
		t.Fatalf("Failed to load data from input: %v", err)
	}

	want := &IniParser{
		Sections: map[string]Section{
			"DEFAULT": {
				Name: "DEFAULT",
				Keys: map[string]string{
					"ServerAliveInterval": "45",
					"Compression":         "yes",
				},
			},
			"forge.example": {
				Name: "forge.example",
				Keys: map[string]string{
					"User": "hg",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestLoadFromFile(t *testing.T) {
	got, err := LoadFromFile("testdata/input.ini")
	if err != nil {
		t.Fatalf("Failed to load data from file: %v", err)
	}
	content, err := os.ReadFile("testdata/input.ini")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	want, err := LoadFromString(string(content))
	if err != nil {
		t.Fatalf("Failed to load expected data from string: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
