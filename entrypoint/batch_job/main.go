package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/averak/gamebox/app/core/build_info"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/core/logger"
	"github.com/averak/gamebox/app/infrastructure/db"
	"github.com/averak/gamebox/app/infrastructure/trace"
	"github.com/google/uuid"
	"github.com/urfave/cli"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Init("batch-job", build_info.ServerVersion())
	if config.Get().GetGoogleCloud().GetTrace().GetEnabled() {
		trace.Init(ctx, "batch-job", build_info.ServerVersion(), config.Get().GetGoogleCloud().GetTrace().GetSamplingRate())
	}

	if err := newCliApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "batch_job"
	app.HelpName = "batch_job"
	app.Usage = "batch job tools"

	var jobCommands cli.Commands
	for name, job := range registry {
		jobCommands = append(jobCommands, cli.Command{
			Name:  name,
			Usage: job.Desc(),
			Action: func(c *cli.Context) error {
				ctx := context.Background()
				logger.Notice(ctx, map[string]any{
					"message": "run batch job",
					"job":     name,
				})

				conn, err := db.NewConnection()
				if err != nil {
					return err
				}
				err = job.Run(ctx, game_context.NewGameContext(uuid.New(), time.Now()), conn)
				if err != nil {
					logger.Error(ctx, map[string]any{
						"message": fmt.Sprintf("failed to run %s", name),
						"error":   err.Error(),
					})
					return err
				}
				return nil
			},
		})
	}
	app.Commands = jobCommands
	return app
}
