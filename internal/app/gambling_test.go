package app

import (
	"sync"
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

func TestWrongCommand(t *testing.T) {

	var message = "!gamble"

	extractCommand(message)

}

func testFilterPossibilities(t *testing.T) {

	var data = []string{"test", "test", "Test", "choice"}
	var expected = []string{"test", "choice"}

	res := filterPossibilities(data)

	if res[0] != expected[0] {
		t.Errorf("[filterPossibilities] Expected : %s, Get %s", expected, res)
	}

	if len(res) != 2 {
		t.Errorf("[filterPossibilities] Choices not filtered correctly")
	}
}

func TestVote(t *testing.T) {
	g := &Gambling{
		CurrentVote: &Vote{
			Possibilities: make([]string, 0),
			Votes:         make(map[string]string),
			Acks:          NewAcks(),
		},
	}
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		g.SendAcks()
		wait.Done()
	}()
	g.CurrentVote.Acks.Drop <- true
	wait.Wait()
}
