package service_registry

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	ServiceRegistryDriverConsul = "consul"
)

// agrs[0] host
// args[1] port
func NewServiceRegistry(driver string, config *ServiceRegistryConfig) (ServiceRegistry, error) {
	switch driver {
	case ServiceRegistryDriverConsul:
		return NewCustomConsulClient(config), nil
	default:
		errMsg := fmt.Sprintf("unknown driver: %s", driver)
		return nil, errors.New(errMsg)
	}
}
