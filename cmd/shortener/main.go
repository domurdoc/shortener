package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/grpc"
	"github.com/domurdoc/shortener/internal/handler"
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, cfgErr := config.LoadConfig()
	if cfgErr != nil {
		log.Fatal(cfgErr)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	closer := utils.NewCloser()

	h := handler.New(a)
	r := router.NewRouter(
		h,
		a.Auth,
		a.Log,
		cfg.TrustedSubnet,
	)
	mainSrv := router.NewServer(
		ctx,
		r,
		a.Log,
		cfg.ServerAddress,
		cfg.ServerCloseTimeout,
	)

	if cfg.ServerEnableHTTPS {
		go func() {
			mainSrv.StartTLS(cfg.ServerCertFile, cfg.ServerKeyFile)
			cancel()
		}()
	} else {
		go func() {
			mainSrv.Start()
			cancel()
		}()
	}
	closer.Register(mainSrv.Close)

	if cfg.ProfilerAddress != "" {
		profSrv := profiler.NewServer(ctx, cfg.ProfilerAddress, a.Log, cfg.ProfilerCloseTimeout)
		go func() {
			profSrv.Start()
			cancel()
		}()
		closer.Register(profSrv.Close)
	}

	if cfg.GRPCPort != 0 {
		grpcSvc := grpc.NewShortenerServiceServer(a)
		grpcSrv := grpc.NewServer(grpcSvc, a.Log, a.Auth, cfg.GRPCPort)
		go func() {
			grpcSrv.Start()
			cancel()
		}()
		closer.Register(grpcSrv.Close)
	}

	<-ctx.Done()
	cancel()
	stop()

	closingErr := closer.Close()
	if err := a.Close(closingErr); err != nil {
		log.Fatal(err)
	}
	log.Println("app stopped gracefully")
}

func printBuildParams() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
