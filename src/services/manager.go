package services

import (
	"context"

	"github.com/Icikowski/GoosyMock/constants"
	"github.com/Icikowski/GoosyMock/data"
	"github.com/Icikowski/GoosyMock/services/admin"
	"github.com/Icikowski/GoosyMock/services/content"
	"github.com/Icikowski/GoosyMock/services/probes"
	"github.com/Icikowski/GoosyMock/utils"
	"github.com/Icikowski/kubeprobes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ServiceManager represents the service manager
type ServicesManager struct {
	log      zerolog.Logger
	probes   *probes.ProbesService
	content  *content.ContentService
	adminApi *admin.AdminAPIService

	routesStoreLog   zerolog.Logger
	payloadsStoreLog zerolog.Logger

	appProbe      kubeprobes.ProbeFunction
	contentProbe  kubeprobes.ProbeFunction
	adminApiProbe kubeprobes.ProbeFunction
}

// NewServicesManager creates ner ServiceManager instance
func NewServicesManager(
	log zerolog.Logger,
	probes *probes.ProbesService,
	content *content.ContentService,
	adminApi *admin.AdminAPIService,
	routesStoreLog zerolog.Logger,
	payloadsStoreLog zerolog.Logger,
	appProbe, contentProbe, adminApiProbe kubeprobes.ProbeFunction,
) *ServicesManager {
	return &ServicesManager{
		log:           log,
		probes:        probes,
		content:       content,
		adminApi:      adminApi,
		appProbe:      appProbe,
		contentProbe:  contentProbe,
		adminApiProbe: adminApiProbe,
	}
}

const (
	fieldService         string = "service"
	msgUnderlyingService string = "starting underlying service"
	msgWaitingForStop    string = "waiting for service to stop"
)

// Run starts the ServicesManager with given context
func (sm *ServicesManager) Run(ctx context.Context) {
	sm.log.Info().Msg("starting Services Manager")

	sm.log.Debug().Msg("generating master context for services cancellation")
	gCtx, cancel := context.WithCancel(context.Background())

	sm.log.Debug().Msg("creating routes store")
	routesStore := data.NewRoutesStore(sm.routesStoreLog)

	sm.log.Debug().Msg("creating payloads store")
	payloadsStore, err := data.NewPayloadsStore(sm.payloadsStoreLog)
	if err != nil {
		sm.log.Fatal().Err(err).Msg("unable to create payloads store")
		cancel()
		return
	}

	sm.log.Debug().Str(fieldService, constants.ComponentHealthProbesService).Msg(msgUnderlyingService)
	sm.probes.Run(gCtx)

	sm.log.Debug().Str(fieldService, constants.ComponentContentService).Msg(msgUnderlyingService)
	sm.content.Run(gCtx, routesStore, payloadsStore)

	sm.log.Debug().Str(fieldService, constants.ComponentAdminAPIService).Msg(msgUnderlyingService)
	sm.adminApi.Run(gCtx, routesStore, payloadsStore)

	sm.log.Info().Msg("services started successfully")

	<-ctx.Done()

	sm.log.Info().Msg("stopping underlying services")
	cancel()

	sm.log.Debug().Str(fieldService, constants.ComponentAdminAPIService).Msg(msgWaitingForStop)
	utils.WaitForProbeDown(sm.adminApiProbe)

	sm.log.Debug().Str(fieldService, constants.ComponentContentService).Msg(msgWaitingForStop)
	utils.WaitForProbeDown(sm.contentProbe)

	sm.log.Debug().Str(fieldService, constants.ComponentHealthProbesService).Msg(msgWaitingForStop)
	utils.WaitForProbeDown(sm.appProbe)

	sm.log.Debug().Str(fieldService, constants.ComponentPayloadsStore).Msg("closing payloads store")
	if err := payloadsStore.Close(); err != nil {
		log.Warn().Err(err).Msg("problem occurred while closing payloads store")
	}

	sm.log.Info().Msg("services stopped successfully")
}
