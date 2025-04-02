package datastore

import (
	"github.com/LuckyMcBeast/stoglr/model"
)

type Datastore interface {
	CreateOrGetToggle(name string, toggleType string, executes string) model.Toggle
	ChangeToggle(name string) model.Toggle
	GetAllToggles() []model.Toggle
	DeleteToggle(name string) model.Toggle
	SetExecution(name string, executes string) model.Toggle
}
