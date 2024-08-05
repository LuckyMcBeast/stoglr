package datastore

import (
	"stoglr/model"
)

type RuntimeDatastore struct {
	db map[string]model.Toggle
}

func NewRuntimeDatastore() *RuntimeDatastore {
	return &RuntimeDatastore{db: make(map[string]model.Toggle)}
}

func (rd *RuntimeDatastore) CreateOrGetToggle(name string, toggleType string, executes string) model.Toggle {
	toggle := model.NewToggle(name, toggleType, executes)
	t, ok := rd.db[name]
	if ok {
		return t
	}
	rd.db[name] = *toggle
	return *toggle
}

func (rd *RuntimeDatastore) ChangeToggle(name string) model.Toggle {
	t, ok := rd.db[name]
	if ok {
		switch t.Status {
		case model.ENABLED:
			t.Status = model.DISABLED
			rd.db[name] = t
		case model.DISABLED:
			t.Status = model.ENABLED
			rd.db[name] = t
		}
		return t
	}
	return *model.NotFound(name)
}

func (rd *RuntimeDatastore) GetAllToggles() []model.Toggle {
	current := rd.db
	ts := make([]model.Toggle, 0, len(current))
	for _, t := range current {
		ts = append(ts, t)
	}
	return ts
}

func (rd *RuntimeDatastore) DeleteToggle(name string) model.Toggle {
	t, ok := rd.db[name]
	if ok {
		delete(rd.db, name)
		t.Status = model.REMOVED
		return t
	}
	return *model.NotFound(name)
}

func (rd *RuntimeDatastore) SetExecution(name string, executes string) model.Toggle {
	t, ok := rd.db[name]
	if ok {
		return rd.updateExecutesIfAB(name, executes, t)
	}
	return *model.NotFound(name)
}

func (rd *RuntimeDatastore) updateExecutesIfAB(name string, executes string, t model.Toggle) model.Toggle {
	if t.ToggleType == model.AB {
		t.UpdateExecutes(executes)
		rd.db[name] = t
	}
	return t
}
