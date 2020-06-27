package app

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCommand(t *testing.T) {

	var cmdExpected string = "create"
	var argExpected string = "val"

	var message = "!gamble create val pl"

	cmd, args := extractCommand(message)

	assert.Equal(t, cmdExpected, cmd)

	assert.Equal(t, argExpected, args[0])

}

func testFilterPossibilities(t *testing.T) {

	var data = []string{"test", "test", "Test", "choice"}
	var expected = []string{"test", "choice"}

	res := filterPossibilities(data)

	assert.Equal(t, expected[0], res[0])

	assert.Equal(t, 2, len(res))
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
