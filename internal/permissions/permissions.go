package permissions

type Role int

const (
	USER_ROLE   Role = 1 // everyday users of the app
	ADMIN_ROLE  Role = 2 // admin users
	SYSTEM_ROLE Role = 3 // LLM bots, task runners etc
)
