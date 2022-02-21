package service

//go:generate mockgen -destination=mock_interfaces.go -package=service -source=interface.go
import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/gofr-crud/model"
)

type Service interface {
	Get(ctx *gofr.Context) ([]model.Customer, error)
	GetByID(ctx *gofr.Context, id string) (model.Customer, error)
	Create(ctx *gofr.Context, customer *model.Customer) (model.Customer, error)
	Update(ctx *gofr.Context, customer *model.Customer) (model.Customer, error)
	Delete(ctx *gofr.Context, id string) error
}
