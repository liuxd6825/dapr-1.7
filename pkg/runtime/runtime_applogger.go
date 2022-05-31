package runtime

import (
	"github.com/liuxd6825/components-contrib/liuxd/common"
	components_v1alpha1 "github.com/liuxd6825/dapr/pkg/apis/components/v1alpha1"
	"github.com/liuxd6825/dapr/pkg/components/liuxd/applogger"
	diag "github.com/liuxd6825/dapr/pkg/diagnostics"
	pubsub_adapter "github.com/liuxd6825/dapr/pkg/runtime/pubsub"
	"strings"
)

func WithApplog(loggers ...applogger.Logger) Option {
	return func(o *runtimeOpts) {
		o.appLoggers = append(o.appLoggers, loggers...)
	}
}

func (a *DaprRuntime) initAppLogger(c components_v1alpha1.Component) error {
	logger, err := a.applogRegistry.Create(c.Spec.Type, c.Spec.Version)
	if err != nil {
		log.Warnf("error creating applogger %s (%s/%s): %s", &c.ObjectMeta.Name, c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "creation")
		return err
	}

	properties := a.convertMetadataItemsToProperties(c.Spec.Metadata)
	consumerID := strings.TrimSpace(properties["consumerID"])
	if consumerID == "" {
		consumerID = a.runtimeConfig.ID
	}
	properties["consumerID"] = consumerID

	getAdapter := func() pubsub_adapter.Adapter {
		return a.getPublishAdapter()
	}

	err = logger.Init(common.Metadata{
		Properties: properties,
	}, getAdapter)

	a.appLogger = logger

	if err != nil {
		log.Warnf("error initializing pub sub %s/%s: %s", c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "init")
		return err
	}

	diag.DefaultMonitoring.ComponentInitialized(c.Spec.Type)
	return nil
}
