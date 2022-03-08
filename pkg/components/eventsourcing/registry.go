package eventsourcing

import (
	"strings"

	es "github.com/dapr/components-contrib/eventsourcing/v1"
	"github.com/dapr/dapr/pkg/components"
	"github.com/pkg/errors"
)

type (
	EventSourcing struct {
		Name          string
		FactoryMethod func() es.EventSourcing
	}

	Registry interface {
		Register(components ...EventSourcing)
		Create(name, version string) (es.EventSourcing, error)
	}

	eventSourcingRegistry struct {
		messageBuses map[string]func() es.EventSourcing
	}
)

func New(name string, factoryMethod func() es.EventSourcing) EventSourcing {
	return EventSourcing{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

func NewRegistry() Registry {
	return &eventSourcingRegistry{
		messageBuses: map[string]func() es.EventSourcing{},
	}
}

// Register registers one or more new message buses.
func (p *eventSourcingRegistry) Register(components ...EventSourcing) {
	for _, component := range components {
		p.messageBuses[createFullName(component.Name)] = component.FactoryMethod
	}
}

// Create instantiates a pub/sub based on `name`.
func (p *eventSourcingRegistry) Create(name, version string) (es.EventSourcing, error) {
	if method, ok := p.getEventSourcing(name, version); ok {
		return method(), nil
	}
	return nil, errors.Errorf("couldn't find message bus %s/%s", name, version)
}

func (p *eventSourcingRegistry) getEventSourcing(name, version string) (func() es.EventSourcing, bool) {
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
	return strings.ToLower("eventsourcing." + name)
}
