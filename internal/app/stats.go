package app

import (
	"strconv"
	"strings"

	"github.com/Ronmi/pastebin"
)

// Create stats from vote and paste it to pastebin
func createStat(key string, votes *Vote) (string, error) {
	api := pastebin.API{Key: key}

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

	return api.Post(&pastebin.Paste{
		Title:    "Stat Vote",
		Content:  str,
		ExpireAt: pastebin.In1D,
	})
}
