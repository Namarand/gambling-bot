package app

import (
	"testing"
)

func TestExtractCommand(t *testing.T) {

	var cmdExpected string = "create"
	var argExpected string = "val"

	var message = "!gamble create val pl"

	cmd, args := extractCommand(message)

	if cmd != cmdExpected {
		t.Errorf("[extractCommand] Expected : %s, Get : %s", cmdExpected, cmd)
	}

	if args[0] != argExpected {
		t.Errorf("[extractCommand] Expected : %s, Get : %s", argExpected, args[0])
	}

}
