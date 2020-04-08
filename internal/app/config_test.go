package app

import (
	"testing"
)

func TestGetConf(t *testing.T) {

	var expectedToken string = "TOKEN"
	var expectedAdminLen int = 3

	var c Conf

	c.getConf("../../tests/config.yml")

	if c.Twitch.Oauth != expectedToken {
		t.Errorf("[c.Twitch.Oauth] Expected : %s, Get : %s", expectedToken, c.Pastebin.Key)
	}

	if len(c.Admins) != 3 {
		t.Errorf("[c.Admins] Expected : %d, Get : %d", expectedAdminLen, len(c.Admins))
	}

}
