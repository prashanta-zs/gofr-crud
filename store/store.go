package store

import (
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/gofr-crud/model"
)

type customer struct{}

func New() Store {
	return customer{}
}

func (c customer) Get(ctx *gofr.Context) ([]model.Customer, error) {
	rows, err := ctx.DB().QueryContext(ctx, get)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	customers := make([]model.Customer, 0)

	for rows.Next() {
		var c model.Customer

		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, errors.Error("error from scan")
		}

		customers = append(customers, c)
	}

	return customers, nil
}

func (c customer) GetByID(ctx *gofr.Context, id int) (model.Customer, error) {
	row := ctx.DB().QueryRowContext(ctx, getByID, id)

	var customer model.Customer

	err := row.Scan(&customer.ID, &customer.Name)
	if err != nil || err == sql.ErrNoRows {
		return model.Customer{}, sql.ErrNoRows
	}

	return customer, nil
}

func (c customer) Create(ctx *gofr.Context, customer *model.Customer) (int, error) {
	res, err := ctx.DB().ExecContext(ctx, insert, customer.Name)
	if err != nil {
		return 0, errors.DB{Err: err}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Error("error from lastInsertID")
	}
	return int(id), nil
}

func (c customer) Update(ctx *gofr.Context, customer *model.Customer) error {
	_, err := ctx.DB().ExecContext(ctx, update, customer.Name, customer.ID)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

func (c customer) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.DB().ExecContext(ctx, delete, id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
