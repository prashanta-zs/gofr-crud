package handler

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/gofr-crud/model"
	"github.com/gofr-crud/service"
)

type handler struct {
	service service.Service
}

func New(service service.Service) handler {
	return handler{service: service}
}

func (h handler) Get(ctx *gofr.Context) (interface{}, error) {
	customers, err := h.service.Get(ctx)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (h handler) GetByID(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	customer, err := h.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (h handler) Create(ctx *gofr.Context) (interface{}, error) {
	var customer model.Customer
	if err := ctx.Bind(&customer); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := h.service.Create(ctx, &customer)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h handler) Update(ctx *gofr.Context) (interface{}, error) {
	var customer model.Customer
	if err := ctx.Bind(&customer); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := h.service.Update(ctx, &customer)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	err := h.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "Deleted Successfully", nil
}
