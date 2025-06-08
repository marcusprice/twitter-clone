package constants

const (
	DEV_ENV  = "DEVELOPMENT"
	TEST_ENV = "TEST"
	PROD_ENV = "PRODUCTION"
)

func DEV_ENVS() []string {
	return []string{DEV_ENV, TEST_ENV}
}
