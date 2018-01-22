package simulator

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func clear(t *testing.T, file string) {
	if err := os.Remove(file); err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("Failed to clear file %s: %s", file, err.Error())
		}
	}
}

func TestFlushInterval(t *testing.T) {
	var (
		fp uint = 10
		mt uint = 1000
		ff      = "test.json"
		k       = "test_key"
		v       = "test_value"
	)
	clear(t, ff)
	sim := NewSimulator(fp, mt, ff)
	sim.Input() <- NewTransaction(k, v)

	time.Sleep(time.Duration(fp+1) * time.Second)

	_, err := os.Stat(ff)
	if os.IsNotExist(err) {
		t.Errorf("Expected - file %s was created. Actual - file doesn`t exist.", ff)
	}
	data, err := ioutil.ReadFile(ff)
	if err != nil {
		t.Errorf("Error reading file %s : %s.", ff, err.Error())
	}
	b := &block{}
	json.Unmarshal(data, b)
	if len(b.Transactions) == 0 {
		t.Error("Expected there is one transaction in block. Actual - no transactions.")
	}
	if b.Transactions[0].Key != k || b.Transactions[0].Value != v {
		t.Errorf(
			"Expected Transactions[0].Key to be %s, b.Transactions[0].Value to be %s. Got key: %s, value %s",
			k, v, b.Transactions[0].Key, b.Transactions[0].Value)
	}

	clear(t, ff)
}
