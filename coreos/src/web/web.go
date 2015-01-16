// web.go: Exposes results from api as HTML.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"hkjn.me/junk/coreos/src/api"
)

var (
	apiServer = flag.String("api_server", "", "If set, HTTP address of API server. If not set, address is read from etcd")
	// Note that we always bind to the same port internally; the
	// .service file can map it to any external port that's desired
	// based on which stage we're running.
	bindAddr = ":9000"
	stage    = "" // prod|staging|testN|dev|unittest
)

type (
	// webHandler handles HTTP requests for monkeys.
	webHandler struct {
		p api.MonkeyAPI // provider of the monkeys
	}

	// jsonAPI implements api.MonkeyAPI by HTTP requests to the JSON API service.
	jsonAPI struct{}
)

// ServeHTTP serves the index page.
func (h webHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Note: We just do some arbitrary checks for stage here to prove
	// that we can have different internal behavior in the app
	// depending on which stage we're running on (prod, test, ..). In
	// a real deployment this should be abstracted away in e.g. a
	// shared package that reads some config defining the stages.
	msg := fmt.Sprintf("Hi from web layer on %s?! I don't even know what I'm supposed to do in this kind of environment!", stage)
	if stage == "test" {
		msg = "Hi from web layer on test!"
	} else if stage == "unittest" {
		msg = "Yes, automatic tester, I'm working as intended."
	} else if stage == "prod" {
		msg = "Hi, I'm the sooper productionized prod web layer!"
	}
	m, err := h.p.GetMonkeys()
	if err != nil {
		log.Printf("Error from API: %v\n", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	// TODO: use tmpl?
	fmt.Fprintf(w, "<html><body><h1>%s</h1><p>%s</p></body></html>", msg, m)
}

// getURL returns the URL for specific API endpoint.
//
// getURL uses specified value from -api_server if present, otherwise
// reads /services/api/[stage] from etcd.
func (p jsonAPI) getURL(endpoint string) string {
	if *apiServer == "" {
		// TODO: read /services/api/%stage% from etcd, get "host" and "port" keys
		// TODO: cache value or read each time?
		log.Fatalf("not implemented yet\n")
		return ""
	} else {
		return fmt.Sprintf("%s/%s", *apiServer, endpoint)
	}
}

func (p jsonAPI) GetMonkey(id int) (*api.Monkey, error) {
	// TODO: query JSON API here.
	return nil, fmt.Errorf("TODO: implement getMonkey")
}

func (p jsonAPI) GetMonkeys() (*api.Monkeys, error) {
	r, err := http.Get(p.getURL("/monkeys"))
	if err != nil {
		return nil, fmt.Errorf("failed to GET /monkeys: %v", err)
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-success status from GET /monkeys: %s", r.Status)
	}
	m := api.Monkeys{}
	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("couldn't decode response: %v", err)
	}
	return &m, nil
}

func (p jsonAPI) AddMonkey(m api.Monkey) error {
	// TODO: pass data to JSON API here.
	return fmt.Errorf("TODO: implement addMonkey")
}

func main() {
	flag.Parse()
	stage = os.Getenv("STAGE")
	if stage == "" {
		log.Fatalf("FATAL: no STAGE set as environment variable")
	}
	fmt.Printf("web layer for %s stage starting..\n", stage)
	http.Handle("/", webHandler{jsonAPI{}})
	log.Fatal(http.ListenAndServe(bindAddr, nil))
}
