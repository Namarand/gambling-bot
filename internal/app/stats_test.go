package app

import (
	"testing"
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

	if st.Total != 2 {
		t.Errorf("[NewStatistics] Expected : %d, Get : %d", 2, st.Total)
	}

	if len(st.Transformed["depraz"]) != 1 {
		t.Errorf("[NewStatistics] Expected : %d, Get : %d", 1, len(st.Transformed["depraz"]))
	}
}

func TestCreateStats(t *testing.T) {
	v := generateVote()

	expected := `Total: 2
levy (1): alice
depraz (1): bob
`

	if createStat(v) != expected {
		t.Errorf("[createStat] Expected : \n%s\n\nGet : \n%s\n", expected, createStat(v))
	}
}
