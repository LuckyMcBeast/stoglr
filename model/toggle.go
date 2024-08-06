package model

import (
	"log"
	"strconv"
	"strings"
)

type Status string

const (
	DISABLED Status = "DISABLED"
	ENABLED  Status = "ENABLED"
	REMOVED  Status = "REMOVED"
	NOTFOUND Status = "NOTFOUND"
	INVALID  Status = "INVALID"
)

type ToggleType string

const (
	RELEASE ToggleType = "RELEASE"
	OPS     ToggleType = "OPS"
	AB      ToggleType = "AB"
	NONE    ToggleType = "NONE"
)

func fromString(s string) ToggleType {
	switch strings.ToUpper(s) {
	case "":
		return RELEASE
	case "RELEASE":
		return RELEASE
	case "OPS":
		return OPS
	case "AB":
		return AB
	default:
		return NONE
	}
}

type Toggle struct {
	Name       string     `json:"name"`
	Status     Status     `json:"status"`
	ToggleType ToggleType `json:"toggleType"`
	Executes   int        `json:"executes"`
}

func (t *Toggle) UpdateExecutes(i string) {
	exe, err := executesFromString(i)
	if err != nil {
		log.Println(err)
		return
	}
	t.Executes = keepWithinRange(exe)
}

func NewToggle(name string, toggleTypeStr string, executesStr string) *Toggle {
	toggleType := fromString(toggleTypeStr)
	executes := 100
	if toggleType == AB && executesStr != "" {
		exe, err := executesFromString(executesStr)
		if err != nil {
			return Invalid(name)
		}
		executes = exe
	}
	return &Toggle{
		Name:       name,
		Status:     DISABLED,
		ToggleType: toggleType,
		Executes:   keepWithinRange(executes),
	}
}

func NotFound(name string) *Toggle {
	return &Toggle{Name: name, Status: NOTFOUND, ToggleType: NONE, Executes: 0}
}

func Invalid(name string) *Toggle {
	return &Toggle{Name: name, Status: INVALID, ToggleType: NONE, Executes: 0}
}

func executesFromString(s string) (int, error) {
	executes, err := strconv.Atoi(s)
	if err != nil || s == "" {
		log.Printf("Error converting executes (%s) to int when executes was required\n", s)
		return 0, err
	}
	return executes, nil
}

func keepWithinRange(executes int) int {
	switch {
	case executes <= 0:
		return 0
	case executes >= 100:
		return 100
	default:
		return executes
	}
}

type TogglesByType struct {
	Release []Toggle
	Ops     []Toggle
	Ab      []Toggle
}

func NewTogglesByType(t []Toggle) *TogglesByType {
	var r []Toggle
	var o []Toggle
	var a []Toggle
	for i := range t {
		switch t[i].ToggleType {
		case RELEASE:
			r = append(r, t[i])
		case OPS:
			o = append(o, t[i])
		case AB:
			a = append(a, t[i])
		}
	}
	return &TogglesByType{
		Release: r,
		Ops:     o,
		Ab:      a,
	}
}
