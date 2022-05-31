package applogger

import (
	"github.com/liuxd6825/components-contrib/liuxd/applog"
	"github.com/liuxd6825/dapr/pkg/components"
	"github.com/pkg/errors"
	"strings"
)

type (
	Logger struct {
		Name          string
		FactoryMethod func() applog.Logger
	}

	Registry interface {
		Register(components ...Logger)
		Create(name, version string) (applog.Logger, error)
	}

	applogRegistry struct {
		messageBuses map[string]func() applog.Logger
	}
)

func New(name string, factoryMethod func() applog.Logger) Logger {
	return Logger{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

func NewRegistry() Registry {
	return &applogRegistry{
		messageBuses: map[string]func() applog.Logger{},
	}
}

// Register registers one or more new message buses.
func (p *applogRegistry) Register(components ...Logger) {
	for _, component := range components {
		p.messageBuses[createFullName(component.Name)] = component.FactoryMethod
	}
}

// Create instantiates a applogger on `name`.
func (p *applogRegistry) Create(name, version string) (applog.Logger, error) {
	if method, ok := p.getAppLoggerSourcing(name, version); ok {
		return method(), nil
	}
	return nil, errors.Errorf("couldn't find applogger %s/%s", name, version)
}

func (p *applogRegistry) getAppLoggerSourcing(name, version string) (func() applog.Logger, bool) {
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
	return strings.ToLower("applogger." + name)
}
