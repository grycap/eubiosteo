package pergamo

type Config struct {
	APIPort int
	SPAPort int

	OnedataPath   string
	SchedulerPath string
	Dispatcher    string

	DBDriver string
	DBPath   string
}

func DefaultConfig() *Config {
	c := &Config{
		APIPort:       10001,
		SPAPort:       3000,
		OnedataPath:   "/tmp/storage",
		SchedulerPath: "/tmp/.galen",
		Dispatcher:    "local",
		DBDriver:      "sqlite3",
		DBPath:        ":memory:",
	}

	return c
}
