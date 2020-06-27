package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// generateVote is used to generate a test vote struct
func generateVote() *Vote {
	vote := new(Vote)
	vote.Votes = make(map[string]string)
	vote.Votes["alice"] = "levy"
	vote.Votes["bob"] = "depraz"

	return vote

}

func TestNewStatistics(t *testing.T) {

	v := generateVote()

	st := NewStatistics(v)

	assert.Equal(t, 2, st.Total)

}

func TestCreateStats(t *testing.T) {
	v := generateVote()

	expected := `Total: 2
levy (1): alice
depraz (1): bob
`

	assert.Equal(t, expected, createStat(v))
}
