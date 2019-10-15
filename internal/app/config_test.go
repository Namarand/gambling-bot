// Package config testing
package config

import (
	"testing"
)

func TestGetConf(t *testing.T) {

	var expectedKey string = "KEY"
	var expectedToken string = "TOKEN"
	var expectedAdminLen int = 3

	var c Conf

	c.getConf("../../tests/config.yml")

	if c.Pastebin.Key != expectedKey {
		t.Errorf("[c.Pastebin.key] Expected : %s, Get : %s", expectedKey, c.Pastebin.Key)
	}

	if c.Twitch.Oauth != expectedToken {
		t.Errorf("[c.Twitch.Oauth] Expected : %s, Get : %s", expectedToken, c.Pastebin.Key)
	}

	if len(c.Admins) != 3 {
		t.Errorf("[c.Admins] Expected : %d, Get : %d", expectedAdminLen, len(c.Admins))
	}

}
