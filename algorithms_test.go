package derk_test

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/zack-alex/derk"
)

type Config struct {
	Salt               string `json:"salt"`
	User               string `json:"user"`
	MasterPasswordHash string `json:"master_password_hash"`
}

var update = flag.Bool("update", false, "update the test data")

func TestAlgorithms(t *testing.T) {
	file, err := os.Open("test-data.json")
	if err != nil {
		t.Fatalf("Failed to open data file: %v", err)
	}

	var data [][]map[string]string
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}
	file.Close()

	changed := false
	var fixedData [][]map[string]string
	for _, item := range data {
		spec := item[0]
		expect := item[1]["secret"]
		fullspec := make(map[string]string)
		fullspec["method"] = spec["method"]
		fullspec["domain"] = "test_domain"
		fullspec["username"] = "test_username"
		res, err := derk.DeriveAndFormat("test_master_password", fullspec)
		if err != nil {
			t.Fatalf("Error deriving and formatting password: %v", err)
		}
		if res != expect {
			changed = true
		}
		fixedData = append(fixedData, []map[string]string{spec, {"secret": res}})
	}

	if *update {
		file, err = os.Create("test-data.json")
		if err != nil {
			t.Fatalf("Failed to open data file for writing: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(fixedData); err != nil {
			t.Fatalf("Failed to encode JSON: %v", err)
		}
	}
	if changed {
		t.Fatalf("Something changed")
	}
}

func TestUnknownAlgorithm(t *testing.T) {
	_, err := derk.DeriveAndFormat(
		"test_master_password",
		map[string]string{
			"method":   "test_unknown_method",
			"domain":   "test_domain",
			"username": "test_username",
		},
	)
	if err == nil {
		t.Fatalf("Expected an error, but got none")
	}
	if !strings.Contains(err.Error(), "Unknown method: test_unknown_method") {
		t.Fatalf("Bad error message")
	}
}
