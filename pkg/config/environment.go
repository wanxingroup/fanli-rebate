package config

type Environment string

const (
	Debug   Environment = "debug"
	Release Environment = "release"
	Test    Environment = "test"
)

func (env Environment) String() string {

	if env == "" {
		return string(Release)
	}

	return string(env)
}
