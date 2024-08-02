package server

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/steinfletcher/apitest"
	"io"
	"net/http"
	"stoglr/server/datastore"
	"testing"
)

var db = datastore.NewRuntimeDatastore()
var tr = NewToggleRouter(db)

var createRouterTestCases = []struct {
	name            string
	method          string
	url             string
	expectedPattern string
	expectedHandler http.HandlerFunc
}{
	{
		name:            "Health",
		url:             "/api/health",
		method:          http.MethodGet,
		expectedPattern: "GET /api/health",
		expectedHandler: tr.getHealth,
	},
	{
		name:            "Get All",
		url:             "/api/toggle",
		method:          http.MethodGet,
		expectedPattern: "GET /api/toggle",
		expectedHandler: tr.getAll,
	},
	{
		name:            "Create or Get",
		url:             "/api/toggle/test1",
		method:          http.MethodPost,
		expectedPattern: "POST /api/toggle/{name}",
		expectedHandler: tr.createOrGet,
	},
	{
		name:            "Enable",
		url:             "/api/toggle/test1/enable",
		method:          http.MethodPut,
		expectedPattern: "PUT /api/toggle/{name}/enable",
		expectedHandler: tr.enable,
	},
	{
		name:            "Disable",
		url:             "/api/toggle/test3/disable",
		method:          http.MethodPut,
		expectedPattern: "PUT /api/toggle/{name}/disable",
		expectedHandler: tr.disable,
	},
	{
		name:            "Set Executes",
		url:             "/api/toggle/test3/execute/50",
		method:          http.MethodPut,
		expectedPattern: "PUT /api/toggle/{name}/execute/{executes}",
		expectedHandler: tr.executes,
	},
	{
		name:            "Delete",
		url:             "/api/toggle/test2",
		method:          http.MethodDelete,
		expectedPattern: "DELETE /api/toggle/{name}",
		expectedHandler: tr.delete,
	},
}

func TestToggleRouter_CreateRouter(t *testing.T) {
	db.CreateOrGetToggle("test1", "RELEASE", "")
	db.CreateOrGetToggle("test2", "OPS", "")
	db.CreateOrGetToggle("test3", "AB", "25")
	for _, tt := range createRouterTestCases {
		t.Run(tt.name, func(t *testing.T) {
			tr = NewToggleRouter(db)
			req := createRequest(tt.method, tt.url, nil)

			var actual = tr.CreateRouter()
			handler, pattern := actual.Handler(req)

			if handler == nil {
				t.Errorf("handler should not be nil")
			}
			if pattern != tt.expectedPattern {
				t.Errorf("expected: %v, actual: %v", tt.expectedPattern, pattern)
			}
			if spew.Sdump(handler) != spew.Sdump(tt.expectedHandler) {
				t.Errorf("expected: %v, actual: %v", spew.Sdump(tt.expectedHandler), spew.Sdump(handler))
			}

			apitest.Handler(handler).
				Method(tt.method).
				URL(tt.url).
				Expect(t).
				Status(200).
				End()
		})
	}
}

func createRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return req
}
