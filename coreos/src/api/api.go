// api.go: Queries db, exposes JSON API:
// 1. GET /monkeys.json: lists all entries
// 2. POST /monkeys.json: create entity of the type
// 3. GET /monkeys/[enc id].json retrieves a specific entity
// 4. PUT /monkeys/[enc id].json updates a specific entity
// 5. DELETE /monkey/[enc id].json: deletes that entity
package api

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"

	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/gorilla/mux"
)

var (
	// TODO: pass in etcdctl-read value for -db_addr in .service
	// file. (we still would need to handle a bad connection by reading
	// from etcd internally to find the new host, but it would be a
	// start).
	dbAddrFlag   = flag.String("db_addr", "", "If set, TCP host for the DB. If not set, address is read from etcd")
	dbAddr       = ""
	buildVersion = flag.String("build_version", "unknown revision", "Build version")
	// Note that we always bind to the same port inside the container; the
	// .service file can map it to any external port that's desired
	// based on which stage we're running.
	bindAddr                        = ":9100"
	stage                           = ""      // prod|staging|testN|dev|unittest
	maxRequestSize            int64 = 1048576 // largest allowed request, in bytes
	statusUnprocessableEntity       = 422
)

// Monkey is an entity we deal with in the API.
type (
	Monkey struct {
		Id        int       `json:"id"`
		Name      string    `json:"name"`
		Birthdate time.Time `json:"birthdate"`
	}
	// Monkeys are a collection of monkey.
	Monkeys []*Monkey
	// MonkeyAPI defines the interface on how we interact with monkeys.
	MonkeyAPI interface {
		GetMonkey(int) (*Monkey, error)
		GetMonkeys() (*Monkeys, error)
		AddMonkey(Monkey) error
		// TODO: add UpdateMonkey, DeleteMonkey.
	}
)

// String returns a human-readable description of the monkey.
func (m Monkey) String() string {
	return fmt.Sprintf("%s (%d) was born on %v", m.Name, m.Id, m.Birthdate.Format("Mon, 02 Jan 2006"))
}

// String returns a human-readable description of the monkeys.
func (ms Monkeys) String() string {
	r := ""
	for i, m := range ms {
		if i > 0 {
			r += ", "
		}
		r += m.String()
	}
	return r
}

// getDbAddr returns the DB address, taken from the -db_addr flag if
// specified, otherwise read from etcd.
func getDBAddr() (string, error) {
	if *dbAddrFlag != "" {
		glog.Infof("-db_addr is specified, so using it: %s\n", *dbAddrFlag)
		return *dbAddrFlag, nil
	}
	addr, err := getEtcdHost()
	if err != nil {
		glog.Errorf("failed to get DB address from etcd: %v", err)
		return "", err
	}
	glog.Infof("etcd says DB can be found at: %s\n", addr)
	return addr, nil
}

// Serve blocks forever, serving the API on bindAddr.
func Serve() {
	flag.Parse()
	stage = os.Getenv("STAGE")
	if stage == "" {
		log.Fatalln("FATAL: no STAGE set as environment variable")
	}
	glog.Errorf("FIXME: stage=%s, -build_version=%s, -db_addr=%s\n", stage, *buildVersion, *dbAddrFlag)
	var err error
	dbAddr, err = getDBAddr()
	if err != nil { // TODO: instead serve 503 Service Unavailable and keep trying to find DB.
		log.Fatalf("FATAL: no DB addr could be found: %v\n", err)
	}
	glog.Infof("[%s] api layer for stage %q starting..\n", *buildVersion, stage)
	glog.Infof("binding to %s\n", bindAddr)
	log.Fatal(http.ListenAndServe(bindAddr, newRouter(apiHandler{jsonAPI{}})))
}

// getEtcdHost returns the Host info from etcd.
func getEtcdHost() (string, error) {
	//	peers := []string{"http://172.17.42.1:4001", "http://10.1.42.1:4001"}
	peers := []string{"http://172.17.42.1:4001"}
	path := fmt.Sprintf("/services/db/%s", stage)
	c := etcd.NewClient(peers)
	// TODO: From within a container we can't just go to 127.0.0.1:4001 for etcd; we need the docker0 interface's IP:
	// https://coreos.com/docs/distributed-configuration/getting-started-with-etcd/#reading-and-writing-from-inside-a-container
	// Is there something simpler than the following?
	// ETCD_ENDPOINT="$(ifconfig docker0 | awk '/\<inet\>/ { print $2}'):4001"

	r, err := c.Get(path, false, false)
	if err != nil {
		return "", fmt.Errorf("failed to read etcd path %s from peers %v: %v", path, peers, err)
	}
	v := r.Node.Value
	glog.Infof("read value %q from %s\n", v, path)
	addr := struct {
		Host    string `json:"host"`
		Port    int    `json:"port"`
		Version string `json:"version"`
	}{}
	err = json.Unmarshal([]byte(v), &addr)
	if err != nil {
		return "", fmt.Errorf("failed to interpret etcd value %q: %v", v, err)
	}
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port), nil
}

func getDB() (*sql.DB, error) {
	user := ""
	password := ""
	// Note: Obviously not secure, in real use we'd have an encrypted
	// config.
	if stage == "test" {
		user = "testuser"
		password = "testsecret"
	} else if stage == "prod" {
		user = "produser"
		password = "prodsecret"
	}
	sqlSource := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		user, password, dbAddr, "monkeydb")
	glog.V(1).Infof("connecting to MySQL at %s..\n", sqlSource)
	db, err := sql.Open("mysql", sqlSource)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

type apiHandler struct {
	api MonkeyAPI
}

type jsonAPI struct{}

func (api jsonAPI) GetMonkey(id int) (*Monkey, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to reach DB: %v", err)
	}
	row := db.QueryRow(`
      SELECT monkeyName, birthDate
      FROM monkeys
      WHERE monkeyId=?`, id)
	name := ""
	sec := int64(0)
	if err = row.Scan(&name, &sec); err != nil {
		return nil, fmt.Errorf("failed to scan: %v", err)
	}
	// Note: If this was exposed to users, we'd need to display it in
	// their own timezone (explicitly selected).
	birthdate := time.Unix(sec, 0).UTC()
	return &Monkey{id, name, birthdate}, nil
}

func (api jsonAPI) GetMonkeys() (*Monkeys, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to contact DB: %v", err)
	}
	rows, err := db.Query(`
      SELECT monkeyId, monkeyName, birthDate
      FROM monkeys
      LIMIT 1000;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query DB: %v", err)
	}
	defer rows.Close()
	monkeys := Monkeys{}
	for rows.Next() {
		id := 0
		name := ""
		sec := int64(0)

		if err = rows.Scan(&id, &name, &sec); err != nil {
			return nil, fmt.Errorf("failed to scan: %v", err)
		}
		// Note: If this was exposed to users, we'd need to display it in
		// their own timezone (explicitly selected).
		birthdate := time.Unix(sec, 0).UTC()
		monkeys = append(monkeys, &Monkey{id, name, birthdate})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %v", err)
	}
	return &monkeys, nil
}

func (api jsonAPI) AddMonkey(m Monkey) error {
	// TODO: insert data into MySQL db here.
	return fmt.Errorf("TODO: implement addMonkey")
}

// newRouter returns a new HTTP router for the endpoints of the API.
func newRouter(h apiHandler) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/monkeys", h.getMonkeys).Methods("GET")
	r.HandleFunc("/monkeys", h.createMonkey).Methods("POST")
	r.HandleFunc("/monkeys/{key}", h.getMonkey).Methods("GET")
	r.HandleFunc("/monkey/{key}", h.updateMonkey).Methods("PUT")
	r.HandleFunc("/monkeys/{key}", h.deleteMonkey).Methods("DELETE")
	return r
}

// getMonkey fetches all monkeys.
func (h apiHandler) getMonkeys(w http.ResponseWriter, r *http.Request) {
	m, err := h.api.GetMonkeys()
	if err != nil {
		glog.Errorf("failed to fetch monkeys: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		glog.Errorf("failed to encode monkeys: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

// createMonkey creates a new monkey.
func (h apiHandler) createMonkey(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxRequestSize))
	if err != nil {
		glog.Errorf("failed to read monkey: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		glog.Errorf("failed to close request: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	m := Monkey{}
	if err := json.Unmarshal(body, &m); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(statusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			glog.Errorf("failed to write encoding error: %v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		return
	}
	err = h.api.AddMonkey(m)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		glog.Errorf("failed to write encoding error: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

// getMonkey fetches a specific monkey.
func (h apiHandler) getMonkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	// Note: In a production environment, we likely should expose hashes
	// of database ids, not the raw ids.
	id, err := strconv.Atoi(vars["key"])
	if err != nil {
		glog.Errorf("bad monkey id %q: %v", vars["key"], err)
		http.Error(w, fmt.Sprintf("No such id %q.", vars["key"]), http.StatusBadRequest)
		return
	}

	m, err := h.api.GetMonkey(id)
	if err != nil {
		glog.Errorf("failed to fetch monkey: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	if m == nil {
		glog.Errorf("no monkey with id %d\n", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		glog.Errorf("failed to encode monkey: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

// updateMonkey updates a monkey.
func (h apiHandler) updateMonkey(w http.ResponseWriter, r *http.Request) {
	msg := "TODO: implement updateMonkey\n"
	glog.Errorf(msg)
	http.Error(w, msg, http.StatusInternalServerError)
}

// deleteMonkey deletes a monkey.
func (h apiHandler) deleteMonkey(w http.ResponseWriter, r *http.Request) {
	msg := "TODO: implement deleteMonkey\n"
	glog.Errorf(msg)
	http.Error(w, msg, http.StatusInternalServerError)
}
