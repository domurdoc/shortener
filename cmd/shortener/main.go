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

	cfg, cfgErr := config.LoadConfig()
	if cfgErr != nil {
		log.Fatal(cfgErr)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	serverCloser := utils.NewCloser()

	s := router.NewServer(
		ctx,
		cfg.ServerAddress,
		cfg.ServerEnableHTTPS,
		cfg.ServerCertFile,
		cfg.ServerKeyFile,
		cfg.ServerCloseTimeout,
		a,
		a.Log,
	)
	go s.Start()
	serverCloser.Register(s.Close)

	if cfg.ProfilerAddress != "" {
		p := profiler.NewServer(ctx, cfg.ProfilerAddress, a.Log, cfg.ProfilerCloseTimeout)
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
