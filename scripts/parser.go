package scripts

import (
	"bytes"
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"strings"
)

func ParseINIFile(path string) (FormFields, error) {
	str, err := getTrimmedINI(path)
	if err != nil {
		return FormFields{}, err
	}

	cfg, err := ini.Load(str)
	if err != nil {
		return FormFields{}, fmt.Errorf("error parsing INI file %v: %s", path, err)
	}

	profile := cfg.Section("PROFILE")
	format := profile.Key("format").String()
	formatParts := strings.Split(format, "|")
	if len(formatParts) != 2 {
		return FormFields{}, fmt.Errorf("invalid format string %v", format)
	}

	menuIdx := 2
	current := profile

	opts := make([]MenuOption, 0)

	for {
		key := fmt.Sprintf("MENU%d", menuIdx)
		current = cfg.Section(key)
		if len(current.Keys()) == 0 {
			break
		}
		opts = append(opts, readMenuItem(current))
		menuIdx++
	}

	record := FormFields{
		Name:        profile.Key("name").String(),
		FormatName:  formatParts[0],
		Format:      formatParts[1],
		Example:     profile.Key("example").String(),
		Description: profile.Key("description").String(),
		Options:     opts,
	}

	return record, nil
}

func readMenuItem(section *ini.Section) MenuOption {
	cases := make(map[string]string)
	i := 0
	for {
		caseName := fmt.Sprintf("case%d", i)
		value := section.Key(caseName).String()
		parts := strings.Split(value, "|")
		if len(parts) != 2 {
			break
		}
		cases[parts[1]] = parts[0]
		i++
	}
	return MenuOption{
		Default:     section.Key("default").String(),
		Description: section.Key("description").String(),
		Cases:       cases,
	}
}

func getTrimmedINI(path string) ([]byte, error) {
	byteArray, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening %v:  %s", path, err)
	}
	profileIdx := bytes.Index(byteArray, []byte("[PROFILE]"))
	batchIdx := bytes.LastIndex(byteArray, []byte("[BATCH]"))
	if profileIdx == -1 || batchIdx == -1 {
		return nil, fmt.Errorf("%v is not a valid profile", path)
	}
	return byteArray[profileIdx:batchIdx], nil
}
