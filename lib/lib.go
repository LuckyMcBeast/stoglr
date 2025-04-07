package lib

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/LuckyMcBeast/stoglr/model"
)

const (
	TOGGLE_PATH   = "/api/toggle"
	TOGGLE_TYPE   = "?type="
	EXECUTES      = "&executes="
	POLL_INTERVAL = 5 * time.Second
)

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

type random interface {
	randomNumber() int
}

type realRandom struct{}

func (r realRandom) randomNumber() int {
	return rand.IntN(100) + 1
}

type StoglrClient struct {
	Url          string
	client       *http.Client
	random       random
	tStore       []model.Toggle
	pollInterval time.Duration
}

func NewStoglrClient(url string) *StoglrClient {
	return &StoglrClient{
		Url:          url,
		client:       &http.Client{},
		random:       realRandom{},
		tStore:       []model.Toggle{},
		pollInterval: POLL_INTERVAL,
	}
}

func NewStoglrClientWithPollInterval(url string, interval time.Duration) *StoglrClient {
	return &StoglrClient{
		Url:          url,
		client:       &http.Client{},
		random:       realRandom{},
		tStore:       []model.Toggle{},
		pollInterval: interval,
	}
}

func (s *StoglrClient) PollToggles() chan os.Signal {
	done := make(chan os.Signal, 1)
	toggles := make(chan []model.Toggle, 1)
	signal.Notify(done, os.Interrupt)

	go s.retrieveAllToggles(toggles)
	go func() {
		for {
			select {
			case toggleList := <-toggles:
				s.tStore = toggleList
			case <-done:
				close(toggles)
				return
			}
		}
	}()
	return done
}

func (s *StoglrClient) retrieveAllToggles(tChannel chan []model.Toggle) {
	for {
		resp, err := s.client.Get(s.Url + TOGGLE_PATH)
		if err != nil {
			log.Println("Error retrieving all toggles", err)
			continue
		}
		toggles := []model.Toggle{}
		err = json.NewDecoder(resp.Body).Decode(&toggles)
		if err != nil {
			log.Println("Error parsing all toggles", err)
			continue
		}

		tChannel <- toggles
		log.Println("Retrieved all toggles", len(toggles))
		time.Sleep(POLL_INTERVAL)
		resp.Body.Close()
	}
}

func (s *StoglrClient) reqUrl(toggleName string, toggleType model.ToggleType, executes int) string {
	return s.Url + TOGGLE_PATH + "/" + toggleName + TOGGLE_TYPE + string(toggleType) + EXECUTES + strconv.Itoa(executes)
}

func (s *StoglrClient) IsEnabled(stoglr *Stoglr) bool {
	toggle := s.checkStore(stoglr)
	if toggle != nil {
		return s.shouldExecute(toggle)
	}
	toggle, err := s.retrieveToggle(stoglr)
	if err != nil {
		return false
	}
	log.Println(toggle)
	return s.shouldExecute(toggle)
}

func (s *StoglrClient) shouldExecute(t *model.Toggle) bool {
	if t.Status == model.ENABLED {
		if t.ToggleType != model.AB {
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

func (s *StoglrClient) checkStore(stoglr *Stoglr) *model.Toggle {
	for _, t := range s.tStore {
		if t.Name == stoglr.ToggleName {
			return &t
		}
	}
	return nil
}
