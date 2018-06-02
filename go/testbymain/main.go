package main

import "testing"

func Hello(name string) string {
	return "Hello, Emacs!"
}

func TestHello(t *testing.T) {
	if Hello("Vim") != "Hello, Vim!" {
		t.Error("Use Vim!!!")
	}
}

func matchString(a, b string) (bool, error) {
	return true, nil
}

func main() {
	testSuite := []testing.InternalTest{
		{Name: "TestHello", F: TestHello},
	}
	testing.Main(matchString, testSuite, nil, nil)
}
