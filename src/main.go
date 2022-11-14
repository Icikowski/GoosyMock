package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Icikowski/GoosyMock/config"
	"github.com/Icikowski/GoosyMock/constants"
	"github.com/Icikowski/GoosyMock/logs"
	"github.com/Icikowski/GoosyMock/meta"
	"github.com/Icikowski/GoosyMock/services"
	"github.com/Icikowski/GoosyMock/services/admin"
	"github.com/Icikowski/GoosyMock/services/content"
	"github.com/Icikowski/GoosyMock/services/probes"
	"github.com/Icikowski/kubeprobes"
	"github.com/caarlos0/env/v6"
)

var (
	cfg config.Config
	lf  *logs.LoggerFactory
)

func init() {
	el := logs.GetEmergencyLogger(os.Stdout)

	cfg = config.Config{
		AdminAPIService: config.ServiceConfig{
			Port:        constants.DefaultCfgAdminAPIPort,
			SecuredPort: constants.DefaultCfgAdminAPISecuredPort,
		},
		ContentService: config.ServiceConfig{
			Port:        constants.DefaultCfgContentPort,
			SecuredPort: constants.DefaultCfgContentSecuredPort,
		},
	}

	if err := env.Parse(&cfg, env.Options{
		Prefix: "GM_",
	}); err != nil {
		el.Fatal().Err(err).Msg("unable to parse configuration from environment variables")
	}

	if err := cfg.AdminAPIService.LoadCerts(); err != nil {
		el.Fatal().Err(err).Msg("unable to load Admin API service certificates")
	}

	if err := cfg.ContentService.LoadCerts(); err != nil {
		el.Fatal().Err(err).Msg("unable to load Content service certificates")
	}

	lf = logs.NewLoggerFactory(cfg.Logging, os.Stdout)
}

func main() {
	log := lf.InstanceFor("cli")
	log.Info().
		Str("version", meta.Version).
		Str("gitCommit", meta.GitCommit).
		Str("binaryType", meta.BinaryType).
		Str("buildTime", meta.BuildTime).
		Msg("starting GoosyMock")

	log.Info().Interface("cfg", cfg).Msg("configuration parsed")

	log.Debug().Msg("creating application context")
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	log.Debug().Msg("creating probes")
	app, adminApi, cnt := kubeprobes.NewStatefulProbe(), kubeprobes.NewStatefulProbe(), kubeprobes.NewStatefulProbe()

	log.Debug().Msg("building service manager")
	mgr := services.NewServicesManager(
		lf.InstanceFor(constants.ComponentServicesManager),
		probes.NewProbesService(
			lf.InstanceFor(constants.ComponentHealthProbesService),
			cfg.HealthProbesPort,
			app, adminApi, cnt,
		),
		content.NewContentService(
			lf.InstanceFor(constants.ComponentContentService),
			cfg.ContentService,
			cnt,
		),
		admin.NewAdminAPIService(
			lf.InstanceFor(constants.ComponentAdminAPIService),
			cfg.AdminAPIService,
			adminApi,
			cfg.MaximumPayloadSize,
		),
		lf.InstanceFor(constants.ComponentRoutesStore),
		lf.InstanceFor(constants.ComponentPayloadsStore),
		app.GetProbeFunction(),
		cnt.GetProbeFunction(),
		adminApi.GetProbeFunction(),
	)

	log.Debug().Msg("starting service manager")
	mgr.Run(ctx)
}
