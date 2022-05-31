package eventstorage

import (
	"strings"

	es "github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/dapr/pkg/components"
	"github.com/pkg/errors"
)

type (
	EventStorage struct {
		Name          string
		FactoryMethod func() es.EventStorage
	}

	Registry interface {
		Register(components ...EventStorage)
		Create(name, version string) (es.EventStorage, error)
	}

	eventStorageRegistry struct {
		messageBuses map[string]func() es.EventStorage
	}
)

func New(name string, factoryMethod func() es.EventStorage) EventStorage {
	return EventStorage{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

func NewRegistry() Registry {
	return &eventStorageRegistry{
		messageBuses: map[string]func() es.EventStorage{},
	}
}

// Register registers one or more new message buses.
func (p *eventStorageRegistry) Register(components ...EventStorage) {
	for _, component := range components {
		p.messageBuses[createFullName(component.Name)] = component.FactoryMethod
	}
}

// Create instantiates a pub/sub based on `name`.
func (p *eventStorageRegistry) Create(name, version string) (es.EventStorage, error) {
	if method, ok := p.getEventStorage(name, version); ok {
		return method(), nil
	}
	return nil, errors.Errorf("couldn't find message bus %s/%s", name, version)
}

func (p *eventStorageRegistry) getEventStorage(name, version string) (func() es.EventStorage, bool) {
	nameLower := strings.ToLower(name)
	versionLower := strings.ToLower(version)
	pubSubFn, ok := p.messageBuses[nameLower+"/"+versionLower]
	if ok {
		return pubSubFn, true
	}
	if components.IsInitialVersion(versionLower) {
		pubSubFn, ok = p.messageBuses[nameLower]
	}
	return pubSubFn, ok
}

func createFullName(name string) string {
	return strings.ToLower("eventstorage." + name)
}
