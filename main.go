package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/logocomune/keybasedocker/container"
	"github.com/logocomune/keybasedocker/message"
	"github.com/logocomune/keybasedocker/processor"
	"github.com/logocomune/keybasedocker/webhook"
	"github.com/pkg/errors"
)

const (
	httpClientTimeOut = 3 * time.Second
	aggregationTime   = 3 * time.Second
)

var (
	//version is the application version. (Injected by make)
	version string
	//shortVersion is the application version. It's used also in frontend and api response (Injected by make)
	shortVersion = "v0.0.0"
	//commit commit hash. (Injected by make)
	commit string
	//branch current branch name. (Injected by make)
	branch string
	//buildDate Build date. (Injected by make)
	buildDate string

	build = "develop"
)

type cfgArgs struct {
	NodeName string
	Docker   struct {
		ShowRunning bool `conf:"default:false"`
		Listen      struct {
			ContainerEvents bool `conf:"default:true"`
			NetworkEvents   bool `conf:"default:true"`
			VolumeEvents    bool `conf:"default:true"`
		}
	}
	Keybase struct {
		Endpoint string `conf:"required"`
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.LUTC)

	if err := run(); err != nil {
		var cDoneErro *container.ContextDoneError
		if errors.As(err, &cDoneErro) {
			os.Exit(0)
		}

		log.Println("error :", err)

		os.Exit(1)
	}
}

func run() error {
	var cfg cfgArgs

	log.Printf("main : Started : Application initializing : version %s (Built: %s)", shortVersion, buildDate)
	log.Println(os.Args[0], "-h", "for help")

	if err := conf.Parse(os.Args[1:], "KD", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("KD", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}

			fmt.Println(usage)

			return nil
		}

		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}

	log.Printf("main : Config :\n%v\n", out)

	defer log.Println("main : Completed")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	defer func() {
		signal.Stop(shutdown)
		cancel()
	}()

	go func() {
		select {
		case <-shutdown:
			log.Println("Shutdown request")
			cancel()
		case <-ctx.Done():
		}
	}()

	nodeName := cfg.NodeName
	if nodeName == "" {
		nodeName, _ = os.Hostname()
	}

	formatter, err := message.NewFormatter(nodeName)
	if err != nil {
		log.Fatalln("formatter error", err)
	}

	kb := webhook.NewKB(cfg.Keybase.Endpoint, httpClientTimeOut, formatter)

	processor := processor.NewProcessor(ctx, aggregationTime, formatter, kb)

	return container.DockerEvents(ctx, processor.Q, container.DockerCfg{
		ContainerEvents: cfg.Docker.Listen.ContainerEvents,
		VolumeEvents:    cfg.Docker.Listen.VolumeEvents,
		NetworkEvents:   cfg.Docker.Listen.ContainerEvents,
		ShowRunning:     cfg.Docker.ShowRunning,
	})
}
