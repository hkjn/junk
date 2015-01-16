// Tests for api service.
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"hkjn.me/timeutils"
)

// fakeAPI is a MonkeyAPI implementation returning static data for testing.
type fakeAPI struct{}

func (api fakeAPI) GetMonkey(int) (*Monkey, error) {
	return &Monkey{1, "Claude", timeutils.Must(timeutils.ParseStd("2008-11-15 01:05"))}, nil
}

func (api fakeAPI) GetMonkeys() (*Monkeys, error) {
	return &Monkeys{
		&Monkey{2, "Bobby", timeutils.Must(timeutils.ParseStd("2013-07-31 12:45"))},
		&Monkey{3, "Jean", timeutils.Must(timeutils.ParseStd("2012-01-15 17:54"))},
	}, nil
}

func (api fakeAPI) AddMonkey(m Monkey) error { return nil }

func TestGetMonkeys(t *testing.T) {
	stage = "unittest"
	router := newRouter(apiHandler{fakeAPI{}})
	req, err := http.NewRequest("GET", "/monkeys", nil)
	if err != nil {
		t.Fatalf("failed to construct request: %v\n", err)
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("want status %d, got %d, with body %q\n", http.StatusOK, resp.Code, resp.Body)
	}
	got := []*Monkey{}
	if err = json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("couldn't decode response: %v\n", err)
	}

	want := []*Monkey{
		&Monkey{2, "Bobby", timeutils.Must(timeutils.ParseStd("2013-07-31 12:45"))},
		&Monkey{3, "Jean", timeutils.Must(timeutils.ParseStd("2012-01-15 17:54"))},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want response %+v, got %+v\n", want, got)
	}
}

func TestGetMonkey(t *testing.T) {
	router := newRouter(apiHandler{fakeAPI{}})
	stage = "unittest"
	req, err := http.NewRequest("GET", "/monkeys/1234", nil)
	if err != nil {
		t.Fatalf("failed to construct request: %v\n", err)
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("want status %d, got %d, with body %q\n", http.StatusOK, resp.Code, resp.Body)
	}
	got := Monkey{}
	if err = json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("couldn't decode response: %v\n", err)
	}

	want := Monkey{1, "Claude", timeutils.Must(timeutils.ParseStd("2008-11-15 01:05"))}
	if got != want {
		t.Errorf("want response %+v, got %+v\n", want, got)
	}
}
