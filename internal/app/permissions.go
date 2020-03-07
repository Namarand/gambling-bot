package app

import "github.com/apex/log"

// Check permission, is the used an admin ?
func checkPermission(user string, admins []string) bool {

	// loop over admins defined in configuration
	for _, e := range admins {
		// if user is admin
		if user == e {
			// ok
			return true
		}
	}

	// if user is not admin
	// not ok
	// log
	log.WithField("user", user).Warn("User is not listed as admin")
	// then return
	return false
}
