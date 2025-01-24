package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"git.sr.ht/~icikowski/goosymock/config"
	"git.sr.ht/~icikowski/goosymock/constants"
	"git.sr.ht/~icikowski/goosymock/logs"
	"git.sr.ht/~icikowski/goosymock/meta"
	"git.sr.ht/~icikowski/goosymock/services"
	"git.sr.ht/~icikowski/goosymock/services/admin"
	"git.sr.ht/~icikowski/goosymock/services/content"
	"git.sr.ht/~icikowski/goosymock/services/probes"
	"github.com/caarlos0/env/v6"
	"pkg.icikowski.pl/kubeprobes"
)

var (
	cfg config.Config
	lf  *logs.LoggerFactory
)

func init() {
	el := logs.GetEmergencyLogger(os.Stdout)

	cfg = config.Config{
		AdminAPIService: config.ServiceConfig{
			Address:        constants.DefaultCfgAdminAPIAddr,
			SecuredAddress: constants.DefaultCfgAdminAPISecuredAddr,
		},
		ContentService: config.ServiceConfig{
			Address:        constants.DefaultCfgContentAddr,
			SecuredAddress: constants.DefaultCfgContentSecuredAddr,
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
	app, _ := kubeprobes.NewManualProbe("app")
	adminApi, _ := kubeprobes.NewManualProbe("adminAPI")
	cnt, _ := kubeprobes.NewManualProbe("contentAPI")

	log.Debug().Msg("building service manager")
	mgr := services.NewServicesManager(
		lf.InstanceFor(constants.ComponentServicesManager),
		probes.NewProbesService(
			lf.InstanceFor(constants.ComponentHealthProbesService),
			cfg.HealthProbesAddr,
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
		app,
		cnt,
		adminApi,
	)

	log.Debug().Msg("starting service manager")
	mgr.Run(ctx)
}
