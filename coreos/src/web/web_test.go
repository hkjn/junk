// Tests for web service.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"hkjn.me/junk/coreos/src/api"
	"hkjn.me/timeutils"
)

// fakeAPI is a fake monkeyAPI implementation for testing.
type fakeAPI struct{}

func (fakeAPI) GetMonkey(id int) (*api.Monkey, error) {
	return nil, fmt.Errorf("bad request: unexpected GetMonkey call")
}

func (fakeAPI) GetMonkeys() (*api.Monkeys, error) {
	return &api.Monkeys{
		&api.Monkey{6, "Noel", timeutils.Must(timeutils.ParseStd("2006-02-21 18:15"))},
		&api.Monkey{14, "Ethan", timeutils.Must(timeutils.ParseStd("2010-12-02 05:52"))},
	}, nil
}

func (fakeAPI) AddMonkey(m api.Monkey) error {
	return fmt.Errorf("bad request: unexpected AddMonkey call")
}

func TestGetURL(t *testing.T) {
	stage = "unittest"
	*apiServer = "FAKE_API_SERVER"
	cases := []struct {
		in   string
		want string
	}{
		{"/", "http://FAKE_API_SERVER/"},
		{"/monkeys", "http://FAKE_API_SERVER/monkeys"},
	}
	for i, tt := range cases {
		got, err := getURL(tt.in)
		if err != nil {
			t.Errorf("[%d] getURL(%q) got error %v, want %v\n", i, tt.in, err, tt.want)
		} else if got != tt.want {
			t.Errorf("[%d] getURL(%q) got %v, want %v\n", i, tt.in, got, tt.want)
		}
	}
}

func TestWeb(t *testing.T) {
	stage = "unittest"
	h := webHandler{fakeAPI{}}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("failed to construct request: %v\n", err)
	}

	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("want status %d, got %d, with body %v\n", http.StatusOK, resp.Code, resp.Body)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v\n", err)
	}
	want := "<html><body><h1>Yes, automatic tester, I'm working as intended.</h1><p>Noel (6) was born on Tue, 21 Feb 2006, Ethan (14) was born on Thu, 02 Dec 2010</p></body></html>"
	got := string(b)
	if got != want {
		t.Fatalf("want response %q, got %q\n", want, got)
	}
}
