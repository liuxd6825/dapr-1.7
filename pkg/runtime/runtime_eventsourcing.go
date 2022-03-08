package runtime

import (
	eventSourcing "github.com/dapr/components-contrib/eventsourcing/v1"
	components_v1alpha1 "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	diag "github.com/dapr/dapr/pkg/diagnostics"
	pubsub_adapter "github.com/dapr/dapr/pkg/runtime/pubsub"
	"strings"
)

func (a *DaprRuntime) initEventSourcing(c components_v1alpha1.Component) error {
	es, err := a.eventSourcingRegistry.Create(c.Spec.Type, c.Spec.Version)
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

	err = es.Init(eventSourcing.Metadata{
		Properties: properties,
	}, getAdapter)

	a.eventSourcing = es

	if err != nil {
		log.Warnf("error initializing pub sub %s/%s: %s", c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "init")
		return err
	}

	diag.DefaultMonitoring.ComponentInitialized(c.Spec.Type)
	return nil
}
