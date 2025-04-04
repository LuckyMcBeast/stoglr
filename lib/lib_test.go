package lib

import (
	"math/rand/v2"
	"net/http"
	"sync"
	"testing"

	"github.com/LuckyMcBeast/stoglr/model"
	"github.com/LuckyMcBeast/stoglr/server"
	"github.com/LuckyMcBeast/stoglr/server/datastore"
)

func TestReleaseStoglr(t *testing.T) {
	stoglr := ReleaseStoglr("test-toggle")
	if stoglr.ToggleType != model.RELEASE || stoglr.Executes != 100 {
		t.Errorf("Expected RELEASE toggle type and 100 executes, got %v", stoglr)
	}
}

func TestOpsStoglr(t *testing.T) {
	stoglr := OpsStoglr("test-toggle")
	if stoglr.ToggleType != model.OPS || stoglr.Executes != 100 {
		t.Errorf("Expected OPS toggle type and 100 executes, got %v", stoglr)
	}
}

func TestABStoglr(t *testing.T) {
	stoglr := ABStoglr("test-toggle", 50)
	if stoglr.ToggleType != model.AB || stoglr.Executes != 50 {
		t.Errorf("Expected AB toggle type and 50 executes, got %v", stoglr)
	}
}

func TestStoglrClient_reqUrl(t *testing.T) {
	tests := []struct {
		name     string
		toggle   string
		typ      model.ToggleType
		executes int
		want     string
	}{
		{
			name:     "RELEASE toggle with default executes",
			toggle:   "test-toggle",
			typ:      model.RELEASE,
			executes: 100,
			want:     "/api/toggle/test-toggle?type=RELEASE&executes=100",
		},
		{
			name:     "OPS toggle with default executes",
			toggle:   "test-toggle",
			typ:      model.OPS,
			executes: 100,
			want:     "/api/toggle/test-toggle?type=OPS&executes=100",
		},
		{
			name:     "AB toggle with custom executes",
			toggle:   "test-toggle",
			typ:      model.AB,
			executes: 50,
			want:     "/api/toggle/test-toggle?type=AB&executes=50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &StoglrClient{Url: ""}
			got := client.reqUrl(tt.toggle, tt.typ, tt.executes)
			if got != tt.want {
				t.Errorf("reqUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

type belowFifty struct {}

func (b belowFifty) randomNumber() int {
	return rand.IntN(49) + 1
}

type aboveFifty struct {}

func  (a aboveFifty) randomNumber() int {
	return rand.IntN(51) + 49
}

func TestStoglrClient_IsEnabled(t *testing.T) {

	tests := []struct {
		name   string
		client *StoglrClient
		stoglr Stoglr
		want   bool
	}{
		{
			name: "should create release toggle",
			client: NewStoglrClient("http://localhost:4444"),
			stoglr: Stoglr{
				ToggleName: "test-toggle1",
				ToggleType: model.RELEASE,
				Executes:   100,
			},
			want: false,
		},
		{
			name: "should create AB toggle",
			client: &StoglrClient{Url: "http://localhost:4444", client: &http.Client{}, random: belowFifty{}},
			stoglr: Stoglr{
				ToggleName: "test-toggle2",
				ToggleType: model.AB,
				Executes:   50,
			},
			want: false,
		},
		{
			name: "should create another AB toggle",
			client: &StoglrClient{Url: "http://localhost:4444", client: &http.Client{}, random: aboveFifty{}},
			stoglr: Stoglr{
				ToggleName: "test-toggle2-another",
				ToggleType: model.AB,
				Executes:   50,
			},
			want: false,
		},
		{
			name: "should create OPs toggle",
			client: NewStoglrClient("http://localhost:4444"),
			stoglr: Stoglr{
				ToggleName: "test-toggle3",
				ToggleType: model.OPS,
				Executes:   100,
			},
			want: false,
		},
		{
			name: "should get release toggle as enabled",
			client: NewStoglrClient("http://localhost:4444"),
			stoglr: Stoglr{
				ToggleName: "test-toggle1",
				ToggleType: model.RELEASE,
				Executes:   100,
			},
			want: true,
		},
		{
			name: "should get AB toggle as enabled when random number is below executes",
			client: &StoglrClient{Url: "http://localhost:4444", client: &http.Client{}, random: belowFifty{}},
			stoglr: Stoglr{
				ToggleName: "test-toggle2",
				ToggleType: model.AB,
				Executes:   50,
			},
			want: true,
		},
		{
			name: "should get AB toggle as disabled when random number is above executes",
			client: &StoglrClient{Url: "http://localhost:4444", client: &http.Client{}, random: aboveFifty{}},
			stoglr: Stoglr{
				ToggleName: "test-toggle2-another",
				ToggleType: model.AB,
				Executes:   50,
			},
			want: false,
		},
		{
			name: "should get OPS toggle as enabled",
			client: NewStoglrClient("http://localhost:4444"),
			stoglr: Stoglr{
				ToggleName: "test-toggle3",
				ToggleType: model.OPS,
				Executes:   100,
			},
			want: true,
		},
	}

	port := "4444"
	db := datastore.NewRuntimeDatastore()
	ts := *server.NewToggleServer(port, db)
	var wg sync.WaitGroup
	wg.Add(1)
	go serverRunner(&ts, &wg)

	for _, tt := range tests {
		actual := tt.client.IsEnabled(&tt.stoglr)
		if actual != tt.want {
			t.Errorf("got %v, want %v", actual, tt.want)
		}
		//provides different values, will not change what is in the database if exists
		toggle := db.CreateOrGetToggle(tt.stoglr.ToggleName, string(model.NONE), "0")
		if !compareStoglrToToggle(tt.stoglr, toggle) {
			t.Errorf("got %v, want %v", toggle, tt.stoglr)
		}
		db.ChangeToggle(tt.stoglr.ToggleName)
	}
	wg.Done()
}

func serverRunner(ts *server.ToggleServer, wg *sync.WaitGroup) {
	defer wg.Done()
	ts.Start()
}

func compareStoglrToToggle(stoglr Stoglr, toggle model.Toggle) bool {
	return stoglr.ToggleName == toggle.Name &&
		string(stoglr.ToggleType) == string(toggle.ToggleType) &&
		stoglr.Executes == toggle.Executes
}
