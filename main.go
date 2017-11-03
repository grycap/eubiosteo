package main

import (
	"os"
	"pergamo/pergamo"

	"github.com/urfave/cli"
)

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

func main() {
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"

	app.Flags = []cli.Flag{
		cli.Int64Flag{
			Name:  "api",
			Value: 10001,
			Usage: "Port of the rest api",
		},
		cli.Int64Flag{
			Name:  "spa",
			Value: 3000,
			Usage: "Port of the web app",
		},
		cli.StringFlag{
			Name:  "onedata",
			Value: "/tmp/storage",
			Usage: "Onedata sync folder",
		},
		cli.StringFlag{
			Name:  "scheduler",
			Value: "/tmp/.galen",
			Usage: "Scheduler utility folder",
		},
		cli.StringFlag{
			Name:  "dispatcher",
			Value: "local",
			Usage: "Two types: local or slurm",
		},
		cli.StringFlag{
			Name:  "dbdriver",
			Value: "sqlite3",
			Usage: "Driver of the db",
		},
		cli.StringFlag{
			Name:  "dbpath",
			Value: ":memory:",
			Usage: "Path of the database",
		},
	}

	app.Action = func(c *cli.Context) error {

		config := pergamo.DefaultConfig()
		config.APIPort = c.Int("api")
		config.SPAPort = c.Int("spa")
		config.OnedataPath = c.String("onedata")
		config.SchedulerPath = c.String("scheduler")
		config.Dispatcher = c.String("dispatcher")
		config.DBDriver = c.String("dbdriver")
		config.DBPath = c.String("dbpath")

		server, err := pergamo.NewServer(config)
		if err != nil {
			panic(err)
		}

		server.Listen()

		return nil
	}

	app.Run(os.Args)
}
