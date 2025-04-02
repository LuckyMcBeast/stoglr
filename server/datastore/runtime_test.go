package datastore

import (
	"github.com/LuckyMcBeast/stoglr/model"
	"reflect"
	"testing"
)

var createOrGetTestCases = []struct {
	name       string
	toggleName string
	toggleType string
	executes   string
	expected   model.Toggle
}{
	{
		name:       "name type and executes provided",
		toggleName: "test",
		toggleType: "AB",
		executes:   "50",
		expected:   model.Toggle{Name: "test", Status: model.DISABLED, ToggleType: model.AB, Executes: 50},
	},
	{
		name:       "existing toggle name provided",
		toggleName: "test-2",
		toggleType: "RELEASE",
		executes:   "",
		expected:   model.Toggle{Name: "test-2", Status: model.DISABLED, ToggleType: model.AB, Executes: 23},
	},
}

func TestRuntimeDatastore_CreateOrGetToggle(t *testing.T) {
	for _, tt := range createOrGetTestCases {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewRuntimeDatastore()
			rd.db["test-2"] = model.Toggle{
				Name:       "test-2",
				Status:     model.DISABLED,
				ToggleType: model.AB,
				Executes:   23,
			}

			actual := rd.CreateOrGetToggle(tt.toggleName, tt.toggleType, tt.executes)

			if actual != tt.expected {
				t.Errorf("expected: %v, actual: %v", tt.expected, actual)
			}
		})
	}
}

var enableTestCases = []struct {
	name       string
	toggleName string
	expected   model.Toggle
}{
	{
		name:       "should return enabled toggle",
		toggleName: "test",
		expected:   model.Toggle{Name: "test", Status: model.ENABLED, ToggleType: model.RELEASE, Executes: 100},
	},
	{
		name:       "should return not found toggle",
		toggleName: "test-2",
		expected:   model.Toggle{Name: "test-2", Status: model.NOTFOUND, ToggleType: model.NONE, Executes: 0},
	},
}

func TestRuntimeDatastore_EnableToggle(t *testing.T) {
	for _, tt := range enableTestCases {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewRuntimeDatastore()
			rd.db["test"] = model.Toggle{
				Name:       "test",
				Status:     model.DISABLED,
				ToggleType: model.RELEASE,
				Executes:   100,
			}

			actual := rd.ChangeToggle(tt.toggleName)

			if actual != tt.expected {
				t.Errorf("expected: %v, actual: %v", tt.expected, actual)
			}
		})
	}
}

var disableTestCases = []struct {
	name       string
	toggleName string
	expected   model.Toggle
}{
	{
		name:       "should return disabled toggle",
		toggleName: "test",
		expected:   model.Toggle{Name: "test", Status: model.DISABLED, ToggleType: model.RELEASE, Executes: 100},
	},
	{
		name:       "should return not found toggle",
		toggleName: "test-2",
		expected:   model.Toggle{Name: "test-2", Status: model.NOTFOUND, ToggleType: model.NONE, Executes: 0},
	},
}

func TestRuntimeDatastore_DisableToggle(t *testing.T) {
	for _, tt := range disableTestCases {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewRuntimeDatastore()
			rd.db["test"] = model.Toggle{
				Name:       "test",
				Status:     model.ENABLED,
				ToggleType: model.RELEASE,
				Executes:   100,
			}

			actual := rd.ChangeToggle(tt.toggleName)

			if actual != tt.expected {
				t.Errorf("expected: %v, actual: %v", tt.expected, actual)
			}
		})
	}
}

var deleteToggleCases = []struct {
	name       string
	toggleName string
	expected   model.Toggle
}{
	{
		name:       "should return removed toggle",
		toggleName: "test",
		expected:   model.Toggle{Name: "test", Status: model.REMOVED, ToggleType: model.RELEASE, Executes: 100},
	},
	{
		name:       "should return not found toggle",
		toggleName: "test-2",
		expected:   model.Toggle{Name: "test-2", Status: model.NOTFOUND, ToggleType: model.NONE, Executes: 0},
	},
}

func TestRuntimeDatastore_DeleteToggle(t *testing.T) {
	for _, tt := range deleteToggleCases {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewRuntimeDatastore()
			rd.db["test"] = model.Toggle{
				Name:       "test",
				Status:     model.DISABLED,
				ToggleType: model.RELEASE,
				Executes:   100,
			}

			actual := rd.DeleteToggle(tt.toggleName)

			if actual != tt.expected {
				t.Errorf("expected: %v, actual: %v", tt.expected, actual)
			}
		})
	}
}

var setExecutionCases = []struct {
	name       string
	toggleName string
	executes   string
	expected   model.Toggle
}{
	{
		name:       "should return toggle with updated execution",
		toggleName: "test",
		executes:   "50",
		expected:   model.Toggle{Name: "test", Status: model.DISABLED, ToggleType: model.AB, Executes: 50},
	},
	{
		name:       "should not update execution if toggle is RELEASE",
		toggleName: "test1",
		executes:   "50",
		expected:   model.Toggle{Name: "test1", Status: model.DISABLED, ToggleType: model.RELEASE, Executes: 100},
	},
	{
		name:       "should not update execution if toggle is OPS",
		toggleName: "test2",
		executes:   "50",
		expected:   model.Toggle{Name: "test2", Status: model.DISABLED, ToggleType: model.OPS, Executes: 100},
	},
	{
		name:       "should return not found toggle",
		toggleName: "test-2",
		executes:   "50",
		expected:   model.Toggle{Name: "test-2", Status: model.NOTFOUND, ToggleType: model.NONE, Executes: 0},
	},
}

func TestRuntimeDatastore_SetExecution(t *testing.T) {
	for _, tt := range setExecutionCases {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewRuntimeDatastore()
			rd.db["test"] = model.Toggle{
				Name:       "test",
				Status:     model.DISABLED,
				ToggleType: model.AB,
				Executes:   10,
			}
			rd.db["test1"] = model.Toggle{
				Name:       "test1",
				Status:     model.DISABLED,
				ToggleType: model.RELEASE,
				Executes:   100,
			}
			rd.db["test2"] = model.Toggle{
				Name:       "test2",
				Status:     model.DISABLED,
				ToggleType: model.OPS,
				Executes:   100,
			}

			actual := rd.SetExecution(tt.toggleName, tt.executes)

			if actual != tt.expected {
				t.Errorf("expected: %v, actual: %v", tt.expected, actual)
			}
		})
	}
}

func TestRuntimeDatastore_GetToggles(t *testing.T) {
	rd := NewRuntimeDatastore()
	rd.db["test1"] = model.Toggle{
		Name:       "test1",
		Status:     model.DISABLED,
		ToggleType: model.RELEASE,
		Executes:   100,
	}
	rd.db["test2"] = model.Toggle{
		Name:       "test2",
		Status:     model.DISABLED,
		ToggleType: model.RELEASE,
		Executes:   100,
	}
	rd.db["test3"] = model.Toggle{
		Name:       "test3",
		Status:     model.DISABLED,
		ToggleType: model.RELEASE,
		Executes:   100,
	}
	expected := []model.Toggle{
		{Name: "test1", Status: model.DISABLED, ToggleType: model.RELEASE, Executes: 100},
		{Name: "test2", Status: model.DISABLED, ToggleType: model.RELEASE, Executes: 100},
		{Name: "test3", Status: model.DISABLED, ToggleType: model.RELEASE, Executes: 100},
	}

	actual := rd.GetAllToggles()

	if len(actual) != len(expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}
