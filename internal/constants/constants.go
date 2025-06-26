package constants

const (
	DEV_ENV     = "DEVELOPMENT"
	TEST_ENV    = "TEST"
	PROD_ENV    = "PRODUCTION"
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

const (
	DALE_COOPER_USER_ID = 3
)

func DEV_ENVS() []string {
	return []string{DEV_ENV, TEST_ENV}
}
