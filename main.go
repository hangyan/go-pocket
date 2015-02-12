package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
)

const (
	VERSION      = "0.0.1"
	POCKET_URL   = "https://getpocket.com"
	CONSUMER_KEY = "37723-f51be1d3abd159ef958ceea6"
	REDIRECT_URI = "https://github.com/hangyan/Silence"
)

var (
	logger = logrus.New()
)

func main() {
	app := cli.NewApp()
	app.Name = "go-pocket"
	app.Usage = "use pocket in terminal"
	app.Version = VERSION
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		authCommand,
		addCommand,
		getCommand,
	}

	app.Run(os.Args)
}
