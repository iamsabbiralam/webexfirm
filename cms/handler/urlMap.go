package handler

func adminViewURLs() map[string]string {
	return map[string]string{
		"register":  registrationURL,
		"login":     loginURL,
		"home":      homeURL,
		"dashboard": dashboardPath,
		"userList":  getAllUsersPath,
	}
}
