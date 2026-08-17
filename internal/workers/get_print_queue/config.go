package get_print_queue

import "time"

type Config struct {
	Name     string        `env:"NAME" envDefault:"get_print_queue_worker"`
	Enabled  bool          `env:"ENABLED" envDefault:"true"`
	Interval time.Duration `env:"INTERVAL" envDefault:"5s"`
}
