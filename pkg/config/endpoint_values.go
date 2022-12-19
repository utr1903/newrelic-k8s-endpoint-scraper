package config

import (
	"sync"
)

// Object to store all values of all endpoints
type EndpointValues struct {
	// To avoid multi-thread read/write into the map
	mux *sync.RWMutex

	// Map to store all values according to endpoints
	// -> Key: endpoint itself
	// -> Val: attributes which the endpoint has exposed
	Values map[Endpoint](map[string]string)
}

func NewEndpointValues() *EndpointValues {
	return &EndpointValues{
		mux:    &sync.RWMutex{},
		Values: make(map[Endpoint](map[string]string)),
	}
}

func (evs *EndpointValues) AddEndpointValues(
	endpoint Endpoint,
	values map[string]string,
) {
	evs.mux.Lock()
	evs.Values[endpoint] = values
	evs.mux.Unlock()
}

func (evs *EndpointValues) GetEndpoints() []Endpoint {
	endpoints := make([]Endpoint, len(evs.Values))

	i := 0
	for endpoint := range evs.Values {
		endpoints[i] = endpoint
		i++
	}
	return endpoints
}

func (evs *EndpointValues) GetEndpointValues(
	endpoint Endpoint,
) map[string]string {
	evs.mux.RLock()
	values := evs.Values[endpoint]
	evs.mux.RUnlock()
	return values
}
