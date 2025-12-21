package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/profiler"
	"github.com/domurdoc/shortener/internal/router"
	"github.com/domurdoc/shortener/internal/utils"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildParams()
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	cfg := config.New()
	if err := config.ParseEnv(cfg); err != nil {
		log.Fatal(err)
	}
	config.ParseArgs(cfg)
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	serverCloser := utils.NewCloser()

	s := router.NewServer(ctx, cfg.Server.Address, a, a.Log, cfg.Server.CloseTimeout)
	go s.Start()
	serverCloser.Register(s.Close)

	if cfg.Profiler.Address != "" {
		p := profiler.New(ctx, cfg.Profiler.Address, a.Log, cfg.Profiler.CloseTimeout)
		go p.Start()
		serverCloser.Register(p.Close)
	}

	<-ctx.Done()
	stop()

	serverClosingErr := serverCloser.Close()
	if err := a.Close(serverClosingErr); err != nil {
		log.Fatal(err)
	}
	log.Println("app stopped gracefully")
}

func printBuildParams() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
