package main

import (
	"log"
	"os"

	"instant-messaging-app/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "instant-messaging-app",
		Usage: "An instant messaging app with which utilizes RabbitMQ",
		Commands: []*cli.Command{
			{
				Name:  "api",
				Usage: "Start the api gateway",
				Action: func(c *cli.Context) error {
					cmd.StartWebServer()
					return nil
				},
			},
			{
				Name:  "user",
				Usage: "Start the UserService daemon",
				Action: func(c *cli.Context) error {
					cmd.StartUserService()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}