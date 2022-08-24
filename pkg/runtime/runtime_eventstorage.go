package runtime

import (
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/gorm_impl"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/mongo_impl"
	components_v1alpha1 "github.com/liuxd6825/dapr/pkg/apis/components/v1alpha1"
	es "github.com/liuxd6825/dapr/pkg/components/liuxd/eventstorage"
	diag "github.com/liuxd6825/dapr/pkg/diagnostics"
	pubsub_adapter "github.com/liuxd6825/dapr/pkg/runtime/pubsub"
	"github.com/pkg/errors"
	"strings"
)

func WithEventStorage(eventsourdings ...es.EventStorage) Option {
	return func(o *runtimeOpts) {
		o.eventStorages = append(o.eventStorages, eventsourdings...)
	}
}

func (a *DaprRuntime) initEventStorage(c components_v1alpha1.Component) error {
	eventStorage, err := a.eventStorageRegistry.Create(c.Spec.Type, c.Spec.Version)
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

	var opts *eventstorage.Options
	switch c.Spec.Type {
	case mongo_impl.ComponentSpecMongo:
		opts, err = mongo_impl.NewMongoOptions(eventStorage.GetLogger(), common.Metadata{Properties: properties}, getAdapter)
	case gorm_impl.ComponentSpecMySql:
		opts, err = gorm_impl.NewMySqlOptions(eventStorage.GetLogger(), common.Metadata{Properties: properties}, getAdapter)
	default:
		err = errors.Errorf("%v 不支持的配置类型", c.Spec.Type)
	}

	if err != nil {
		log.Warnf("error initializing pub sub %s/%s: %s", c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "init")
		return err
	}
	if err = eventStorage.Init(opts); err != nil {
		return err
	}

	if err != nil {
		log.Warnf("error initializing pub sub %s/%s: %s", c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "init")
		return err
	}
	a.eventStorage = eventStorage

	diag.DefaultMonitoring.ComponentInitialized(c.Spec.Type)
	return nil
}
