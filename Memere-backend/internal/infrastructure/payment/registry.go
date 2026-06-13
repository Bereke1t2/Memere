package payment

import (
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Registry resolves a configured provider by its entity name (spec §10.1). The
// usecase asks for the provider named on a payment row, so adding a provider is
// a wiring change here, never a usecase change.
type Registry struct {
	providers map[entity.PaymentProvider]service.PaymentProvider
}

var _ service.PaymentProviderRegistry = (*Registry)(nil)

// NewRegistry builds a Registry from the given providers, keyed by Name().
func NewRegistry(providers ...service.PaymentProvider) *Registry {
	m := make(map[entity.PaymentProvider]service.PaymentProvider, len(providers))
	for _, p := range providers {
		if p == nil {
			continue
		}
		m[entity.PaymentProvider(p.Name())] = p
	}
	return &Registry{providers: m}
}

// Get returns the provider registered for p, or PROVIDER_NOT_CONFIGURED.
func (r *Registry) Get(p entity.PaymentProvider) (service.PaymentProvider, error) {
	prov, ok := r.providers[p]
	if !ok {
		return nil, apperror.New(400, "PROVIDER_NOT_CONFIGURED",
			"payment provider not configured: "+string(p), nil)
	}
	return prov, nil
}
