package scripts

import (
	"github.com/caseyjmorris/threadsquish/testHelpers"
	"reflect"
	"testing"
)

func Test_ParseINIFile(t *testing.T) {
	file, err := ParseINIFile(testHelpers.GetFixturePath("simpleSample.cmd"))
	if err != nil {
		t.Error(err)
		return
	}
	expected := FormFields{
		Name:        "Sonic the Hedgehog",
		FormatName:  "CRI Movie 2 (*.usm)",
		Format:      "*.usm",
		Example:     "hedge.usm",
		Description: "Gotta go fast!",
		Options: []MenuOption{
			{
				Default:     "Select aspect ratio",
				Description: "",
				Cases: map[string]string{
					"16_9": "16:9 (HD)",
					"21_9": "21:9 (Ultra-wide)",
					"32_9": "32:9 (Super ultra-wide)",
				},
			},
			{
				Default:     "Select resolution profile",
				Description: "You can pick your favorite resolution profile.",
				Cases: map[string]string{
					"high": "Optimal (Higher quality)",
					"low":  "Reduced (Higher performance)",
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, file) {
		t.Errorf("Records don't match.  \nExpected:  %v\nActual:  %v", expected, file)
	}
}
