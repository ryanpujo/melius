package registry

import (
	"database/sql"

	"github.com/ryanpujo/melius/internal/adapter"
)

type Registry struct {
	db *sql.DB
}

func NewRegistry(db *sql.DB) *Registry {
	return &Registry{
		db: db,
	}
}

func (r *Registry) NewAppControllers() *adapter.Adapter {
	return &adapter.Adapter{
		CredentialController: r.GetCredentialController(),
	}
}
