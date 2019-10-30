package app

import (
	"strconv"
	"strings"

	"github.com/Ronmi/pastebin"
)

// Create stats from vote
func createStat(votes *Vote) string {

	transformed := make(map[string][]string)
	for user, value := range votes.Votes {
		transformed[value] = append(transformed[value], user)
	}

	sum := 0
	for _, users := range transformed {
		sum += len(users)
	}

	str := "Total amount: " + strconv.Itoa(sum) + "\n"
	for value, users := range transformed {
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
