package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/karlskewes/pricing"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.Print("pricing server starting")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := pricing.NewApp(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()

		log.Print("received OS signal to shutdown, use Ctrl+C again to force")

		// reset signals so a second ctrl+c will terminate the application.
		stop()

		return nil
	})

	group.Go(func() error {
		return app.Run(ctx)
	})

	if err := group.Wait(); err != nil {
		log.Print(err)
	}
}
