package datastore

import (
	"stoglr/model"
)

type Datastore interface {
	CreateOrGetToggle(name string, toggleType string, executes string) model.Toggle
	EnableToggle(name string) model.Toggle
	DisableToggle(name string) model.Toggle
	GetAllToggles() []model.Toggle
	DeleteToggle(name string) model.Toggle
	SetExecution(name string, executes string) model.Toggle
}
