// web.go: Exposes results from api as HTML.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang/glog"

	"hkjn.me/junk/coreos/src/api"
	"hkjn.me/junk/coreos/src/etcdwrapper"
)

var (
	apiServer    = flag.String("api_server", "", "If set, HTTP address of API server. If not set, address is read from etcd")
	buildVersion = flag.String("web_version", "unknown revision", "Build version of web server")
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
	// TODO: use html.tmpl.
	fmt.Fprintf(w, "<html><body><h1>%s</h1><h2>Here's your monkeys:</h2><p>%s</p></body></html>", msg, m)
}

// getURL returns the URL for a specific API endpoint.
//
// getURL uses specified value from -api_server if present, otherwise
// reads /services/api/[stage] from etcd.
func getURL(endpoint string) (string, error) {
	server := *apiServer
	if server == "" {
		addr, err := etcdwrapper.Read(fmt.Sprintf("/services/api/%s", stage))
		if err != nil {
			glog.Errorf("failed to get API server from etcd: %v", err)
			return "", err
		}
		glog.Infof("etcd says API server can be found at: %s\n", addr)
		server = addr
	}
	return fmt.Sprintf("http://%s%s", server, endpoint), nil
}

func (p jsonAPI) GetMonkey(id int) (*api.Monkey, error) {
	// TODO: query JSON API here.
	return nil, fmt.Errorf("TODO: implement getMonkey")
}

func (p jsonAPI) GetMonkeys() (*api.Monkeys, error) {
	target, err := getURL("/monkeys")
	if err != nil {
		// TODO: Should return 503 Service Unavailable here - only can
		// happen if etcd doesn't know about our API backend.
		return nil, fmt.Errorf("failed to find URL: %v", err)
	}
	r, err := http.Get(target)
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
	glog.V(2).Infof("web starting with stage=%s, -web_version=%s, -api_server=%s\n", stage, *buildVersion, *apiServer)
	if stage == "" {
		log.Fatalf("FATAL: no STAGE set as environment variable")
	}
	fmt.Printf("[%s] web layer for stage %q binding to %s..\n", *buildVersion, stage, bindAddr)
	http.Handle("/", webHandler{jsonAPI{}})
	log.Fatal(http.ListenAndServe(bindAddr, nil))
}
