package main

import (
	"log"

	"github.com/tstranex/u2f"
)

func main() {
	app_id := "http://localhost"
	c, err := u2f.NewChallenge(app_id, []string{app_id})
	if err != nil {
		log.Fatalf("NewChallenge failed: %v\n", err)
	}

	existingTokens := []u2f.Registration{}
	req := u2f.NewWebRegisterRequest(c, existingTokens)
	// Send registration request to the browser.
	log.Printf("Should return request %v to browser\n", req)
	// High-level API
	// u2f.register(<Application id>,
	// [<RegisterRequest instance>, ...],
	// [<RegisteredKey for known token 1>, ...],
	// registerResponseHandler);

	// Read resp from the browser.
	var resp u2f.RegisterResponse
	reg, err := u2f.Register(resp, *c, nil)
	if err != nil {
		log.Fatalf("Register failed: %v\n", err)
	}

	log.Printf("Should store registration %v in the database to browser.\n", reg)
}
