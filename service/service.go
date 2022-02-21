package service

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/gofr-crud/model"
	"github.com/gofr-crud/store"
	"strconv"
)

type service struct {
	store store.Store
}

func New(store store.Store) service {
	return service{store: store}
}

func (s service) Get(ctx *gofr.Context) ([]model.Customer, error) {
	customers, err := s.store.Get(ctx)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (s service) GetByID(ctx *gofr.Context, id string) (model.Customer, error) {
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return model.Customer{}, errors.InvalidParam{Param: []string{"id"}}
	}

	if customerID < 0 {
		return model.Customer{}, errors.InvalidParam{Param: []string{"id"}}
	}

	customer, err := s.store.GetByID(ctx, customerID)
	if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}

func (s service) Create(ctx *gofr.Context, customer *model.Customer) (model.Customer, error) {
	if customer.Name == "" {
		return model.Customer{}, errors.InvalidParam{Param: []string{"Name"}}
	}

	id, err := s.store.Create(ctx, customer)
	if err != nil {
		return model.Customer{}, err
	}

	cust, err := s.store.GetByID(ctx, id)
	if err != nil {
		return model.Customer{}, err
	}

	return cust, nil
}

func (s service) Update(ctx *gofr.Context, customer *model.Customer) (model.Customer, error) {
	if customer.Name == "" {
		return model.Customer{}, errors.InvalidParam{Param: []string{"Name"}}
	}

	if customer.ID < 1 {
		return model.Customer{}, errors.InvalidParam{Param: []string{"id"}}
	}

	err := s.store.Update(ctx, customer)
	if err != nil {
		return model.Customer{}, err
	}

	cust, err := s.store.GetByID(ctx, customer.ID)
	if err != nil {
		return model.Customer{}, err
	}

	return cust, nil
}

func (s service) Delete(ctx *gofr.Context, id string) error {
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return errors.InvalidParam{Param: []string{"id"}}
	}
	if customerID < 0 {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	err = s.store.Delete(ctx, customerID)
	if err != nil {
		return err
	}

	return nil
}
