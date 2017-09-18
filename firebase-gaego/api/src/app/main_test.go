package app

import (
	"fmt"
	"os"
	"testing"

	// do testerator feature setup
	_ "github.com/favclip/testerator/datastore"
	_ "github.com/favclip/testerator/memcache"

	"github.com/favclip/testerator"
)

func TestMain(m *testing.M) {
	_, _, err := testerator.SpinUp()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	status := m.Run()

	err = testerator.SpinDown()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}
