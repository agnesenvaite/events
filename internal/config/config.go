package config

const (
	EnvLocal      = "local"
	EnvProduction = "production"
)

type Config struct {
	APIPort             int    `conf:"env:API_PORT"`
	Environment         string `conf:"env:ENVIRONMENT,default:local"`
	MaxInvitees         int    `conf:"env:MAX_INVITEES"`
	DefaultVideoQuality string `conf:"env:DEFAULT_VIDEO_QUALITY"`
	DefaultAudioQuality string `conf:"env:DEFAULT_AUDIO_QUALITY"`
	Database            database
	Docs                docs
}

func (c *Config) InDevelopment() bool {
	return c.Environment == "" || c.Environment == EnvLocal
}

func (c *Config) InProduction() bool {
	return c.Environment == EnvProduction
}

type database struct {
	URL string `conf:"env:DATABASE_URL"`
}

type docs struct {
	Host string `conf:"env:DOCS_HOST,default:localhost:8080"`
}
