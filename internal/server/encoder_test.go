package server

import (
	"demo/internal/errors"
	"encoding/json"

	"fmt"
	"testing"
)

func TestEncoder(t *testing.T) {
	a := &errors.HTTPError{
		Errors: make(map[string][]string),
	}
	a.Errors["body"] = []string{"empty"}
	b, _ := json.Marshal(a)
	fmt.Println(string(b))
}
