package model

import "testing"

var newToggleCases = []struct {
	name       string
	toggleName string
	toggleType string
	executes   string
	expected   Toggle
}{
	{
		name:       "only name provided",
		toggleName: "test",
		toggleType: "",
		executes:   "",
		expected:   Toggle{Name: "test", Status: DISABLED, ToggleType: RELEASE, Executes: 100},
	},
	{
		name:       "name and type provided",
		toggleName: "test",
		toggleType: "OPS",
		executes:   "",
		expected:   Toggle{Name: "test", Status: DISABLED, ToggleType: OPS, Executes: 100},
	},
	{
		name:       "name type and executes provided",
		toggleName: "test",
		toggleType: "AB",
		executes:   "50",
		expected:   Toggle{Name: "test", Status: DISABLED, ToggleType: AB, Executes: 50},
	},
	{
		name:       "invalid when AB and executes is not parsable",
		toggleName: "test",
		toggleType: "AB",
		executes:   "asdfasd",
		expected:   Toggle{Name: "test", Status: INVALID, ToggleType: NONE, Executes: 0},
	},
	{
		name:       "executes should be 100 when AB and no executes is provided",
		toggleName: "test",
		toggleType: "AB",
		executes:   "",
		expected:   Toggle{Name: "test", Status: DISABLED, ToggleType: AB, Executes: 100},
	},
	{
		name:       "executes provided with RELEASE overridden",
		toggleName: "test",
		toggleType: "RELEASE",
		executes:   "50",
		expected:   Toggle{Name: "test", Status: DISABLED, ToggleType: RELEASE, Executes: 100},
	},
	{
		name:       "executes provided with OPS overridden",
		toggleName: "test",
		toggleType: "OPS",
		executes:   "25",
		expected:   Toggle{Name: "test", Status: DISABLED, ToggleType: OPS, Executes: 100},
	},
}

func TestShouldMakeNewToggle(t *testing.T) {
	for _, tt := range newToggleCases {
		t.Run(tt.name, func(t *testing.T) {

			actual := NewToggle(tt.toggleName, tt.toggleType, tt.executes)

			if *actual != tt.expected {
				t.Errorf("expected: %v, actual: %v", tt.expected, actual)
			}
		})
	}
}
