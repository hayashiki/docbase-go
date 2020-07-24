package testutil

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func LoadFixture(t *testing.T, n string) string {
	fixtureDir := "./testdata"
	p := filepath.Join(fixtureDir, n)

	b, err := ioutil.ReadFile(p)

	if err != nil {
		t.Fatalf("Error while trying to read %s: %v\n", n, err)
	}

	return string(b)
}
