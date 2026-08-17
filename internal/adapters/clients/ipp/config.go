package ipp

type Config struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"631"`
	Username string `env:"USERNAME" envDefault:""`
	Password string `env:"PASSWORD" envDefault:""`
	UseTLS   bool   `env:"USE_TLS" envDefault:"false"`
}
