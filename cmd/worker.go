package cmd

import (
	"dev.sigpipe.me/dashie/reel2bits/models"
	"dev.sigpipe.me/dashie/reel2bits/pkg/mailer"
	"dev.sigpipe.me/dashie/reel2bits/setting"
	"dev.sigpipe.me/dashie/reel2bits/workers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Worker cli target
var Worker = cli.Command{
	Name:        "worker",
	Usage:       "Start workers",
	Description: "It starts the reel2bits workers",
	Action:      runWorker,
	Flags: []cli.Flag{
		stringFlag("config, c", "config/app.ini", "Custom config file path"),
	},
}

func runWorker(ctx *cli.Context) error {
	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
	}

	setting.InitConfig()
	models.InitDb()
	mailer.NewContext()

	server, err := workers.CreateServer()
	if err != nil {
		return err
	}

	worker := server.NewWorker("transcoding_infos")
	err = worker.Launch()
	if err != nil {
		log.Errorf("Launching worker transcoding_infos error: %s", err)
		return err
	}

	return nil
}
