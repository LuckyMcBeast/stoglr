package lib

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/LuckyMcBeast/stoglr/model"
)

const (
	TOGGLE_PATH = "/api/toggle/"
	TOGGLE_TYPE = "?type="
	EXECUTES    = "&executes="
)

type random interface {
	randomNumber() int
}

type realRandom struct {}

func (r realRandom) randomNumber() int  {
	return rand.IntN(100) + 1
}

type StoglrClient struct {
	Url    string
	client *http.Client
	random random
}

func NewStoglrClient(url string) *StoglrClient {
	return &StoglrClient{
		Url:    url,
		client: &http.Client{},
		random: realRandom{},
	}
}

type Stoglr struct {
	ToggleName string
	ToggleType model.ToggleType
	Executes   int
}

func ReleaseStoglr(toggleName string) *Stoglr {
	return &Stoglr{
		ToggleName: toggleName,
		ToggleType: model.RELEASE,
		Executes:   100,
	}
}

func OpsStoglr(toggleName string) *Stoglr {
	return &Stoglr{
		ToggleName: toggleName,
		ToggleType: model.OPS,
		Executes:   100,
	}
}

func ABStoglr(toggleName string, executes int) *Stoglr {
	return &Stoglr{
		ToggleName: toggleName,
		ToggleType: model.AB,
		Executes:   executes,
	}
}

func (s *StoglrClient) reqUrl(toggleName string, toggleType model.ToggleType, executes int) string {
	return s.Url + TOGGLE_PATH + toggleName + TOGGLE_TYPE + string(toggleType) + EXECUTES + strconv.Itoa(executes)
}

func (s *StoglrClient) IsEnabled(stoglr *Stoglr) bool {
	toggle, err := s.retrieveToggle(stoglr)
	if err != nil {
		return false
	}
	log.Println(toggle)
	return s.shouldExecute(toggle)
}

func (s *StoglrClient) shouldExecute(t *model.Toggle) bool {
	if t.Status == model.ENABLED {
		if !(t.ToggleType == model.AB) {
			return true
		}
		if t.Executes >= s.random.randomNumber() {
			return true
		}
	}
	return false
}

func (s *StoglrClient) retrieveToggle(stoglr *Stoglr) (*model.Toggle, error) {
	toggleName, toggleType, executes := stoglr.ToggleName, stoglr.ToggleType, stoglr.Executes
	req, err := http.NewRequest(http.MethodPost, s.reqUrl(toggleName, toggleType, executes), nil)
	if err != nil {
		log.Printf("Treating toggle %v as disabled, error in creating request: %v", toggleName, err)
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Treating toggle %v as disabled, error retrieving toggle: %v", toggleName, err)
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var toggle model.Toggle
	err = decoder.Decode(&toggle)
	if err != nil {
		log.Printf("Treating toggle %v as disabled, error parsing toggle: %v", toggleName, err)
		return nil, err
	}
	defer resp.Body.Close()
	return &toggle, nil
}
