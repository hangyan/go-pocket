package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
)

const (
	VERSION      = "0.0.1"
	POCKET_URL   = "https://getpocket.com/v3"
	CONSUMER_KEY = "31547-629fbec67e06e4af52cf970b"
	REDIRECT_URI = "http://yayua.github.io"
)

var (
	logger = logrus.New()
)

func main() {
	app := cli.NewApp()
	app.Name = "silence"
	app.Usage = "Read your pocket in silence"
	app.Version = VERSION
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		authCommand,
	}

	app.Run(os.Args)
}
