package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Ronmi/pastebin"
)

// Statistics is a struct containing generated stats for a given vote
type Statistics struct {
	Total       int
	Transformed map[string][]string
}

// NewStatistics if used to transform a vote into a statistics struct
func NewStatistics(votes *Vote) Statistics {

	tr := make(map[string][]string)
	for u, v := range votes.Votes {
		tr[v] = append(tr[v], u)
	}

	return Statistics{
		Total:       len(votes.Votes),
		Transformed: tr,
	}

}

// Create stats from vote
func createStat(votes *Vote) string {

	stats := NewStatistics(votes)

	str := "Total: " + strconv.Itoa(stats.Total) + "\n"
	for value, users := range stats.Transformed {
		str += value + " (" + strconv.Itoa(len(users)) + "): " + strings.Join(users, ", ") + "\n"
	}

	return str

}

// Push stats to pastebin as string
func statsToPastebin(key string, stats string) (string, error) {
	api := pastebin.API{Key: key}

	return api.Post(&pastebin.Paste{
		Title:    "Stat Vote",
		Content:  stats,
		ExpireAt: pastebin.In1D,
	})

}

// Write stats into a file inside a base directory
func statsToFile(stats string, dir string) error {

	// Create a directory using current date
	// get current date
	dt := time.Now()

	// forge base dir
	basedir := fmt.Sprintf("%s/%s", dir, dt.Format("2006-01-02"))
	// create base dir if not exists
	if _, err := os.Stat(basedir); os.IsNotExist(err) {
		os.Mkdir(basedir, 0755)
	}

	// Create file
	f, err := os.Create(fmt.Sprintf("%s/stats", basedir))
	if err != nil {
		return err
	}
	// Ensure file is closed at the end of the func
	defer f.Close()

	// Write stuff and return err
	_, err = f.WriteString(stats)

	return err

}
