package runtime

import (
	"github.com/dapr/components-contrib/liuxd/common"
	components_v1alpha1 "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	es "github.com/dapr/dapr/pkg/components/liuxd/eventstorage"
	diag "github.com/dapr/dapr/pkg/diagnostics"
	pubsub_adapter "github.com/dapr/dapr/pkg/runtime/pubsub"
	"strings"
)

func WithEventStorage(eventsourdings ...es.EventStorage) Option {
	return func(o *runtimeOpts) {
		o.eventStorages = append(o.eventStorages, eventsourdings...)
	}
}

func (a *DaprRuntime) initEventStorage(c components_v1alpha1.Component) error {
	es, err := a.eventStorageRegistry.Create(c.Spec.Type, c.Spec.Version)
	if err != nil {
		log.Warnf("error creating pub sub %s (%s/%s): %s", &c.ObjectMeta.Name, c.Spec.Type, c.Spec.Version, err)
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

	err = es.Init(common.Metadata{
		Properties: properties,
	}, getAdapter)

	a.eventStorage = es

	if err != nil {
		log.Warnf("error initializing pub sub %s/%s: %s", c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "init")
		return err
	}

	diag.DefaultMonitoring.ComponentInitialized(c.Spec.Type)
	return nil
}
