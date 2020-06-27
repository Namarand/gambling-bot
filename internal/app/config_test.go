package app

import (
	"testing"

	"github.com/tj/assert"
)

func TestGetConf(t *testing.T) {

	var expectedToken string = "TOKEN"
	var expectedAdminLen int = 3

	var c Conf

	c.getConf("../../tests/config.yml")

	assert.Equal(t, expectedToken, c.Twitch.Oauth)

	assert.Equal(t, expectedAdminLen, len(c.Admins))
}
