package store

//go:generate mockgen -destination=mock_interfaces.go -package=store -source=interface.go

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/gofr-crud/model"
)

type Store interface {
	Get(ctx *gofr.Context) ([]model.Customer, error)
	GetByID(ctx *gofr.Context, id int) (model.Customer, error)
	Create(ctx *gofr.Context, customer *model.Customer) (int, error)
	Update(ctx *gofr.Context, customer *model.Customer) error
	Delete(ctx *gofr.Context, id int) error
}
