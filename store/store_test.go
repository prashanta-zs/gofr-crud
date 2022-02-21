package store

import (
	"context"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofr-crud/model"
	"reflect"
	"testing"
)

func getTestData(t *testing.T) (*gofr.Context, sqlmock.Sqlmock) {
	app := gofr.New()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error %s was not expected when opening a stub database connection", err)
	}

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	ctx.DB().DB = db

	return ctx, mock
}

func TestCustomer_Get(t *testing.T) {
	ctx, mock := getTestData(t)

	row := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "abc")
	scanRow := sqlmock.NewRows([]string{"id"}).AddRow(1)
	dbErr := errors.Error("error from db")

	tests := []struct {
		desc      string
		expOutput []model.Customer
		expErr    error
		mockQuery *sqlmock.ExpectedQuery
	}{
		{
			desc: "Get all customers ",
			expOutput: []model.Customer{
				{
					ID:   1,
					Name: "abc",
				},
			},
			expErr:    nil,
			mockQuery: mock.ExpectQuery(get).WillReturnRows(row),
		},
		{
			desc:      "DB error",
			expOutput: nil,
			expErr:    errors.DB{Err: dbErr},
			mockQuery: mock.ExpectQuery(get).WillReturnError(dbErr),
		},
		{
			desc:      "error from scan",
			expOutput: nil,
			expErr:    errors.Error("error from scan"),
			mockQuery: mock.ExpectQuery(get).WillReturnRows(scanRow),
		},
	}

	for i, tc := range tests {
		store := New()

		customers, err := store.Get(ctx)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(customers, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, customers)

		}
	}
}

func TestCustomer_GetByID(t *testing.T) {
	ctx, mock := getTestData(t)

	row := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "abc")

	tests := []struct {
		desc      string
		id        int
		expOutput model.Customer
		expErr    error
		mockQuery *sqlmock.ExpectedQuery
	}{
		{
			desc:      "Success",
			id:        1,
			expOutput: model.Customer{ID: 1, Name: "abc"},
			expErr:    nil,
			mockQuery: mock.ExpectQuery(getByID).WillReturnRows(row),
		},
		{
			desc:      "fail",
			id:        2,
			expOutput: model.Customer{},
			expErr:    sql.ErrNoRows,
			mockQuery: mock.ExpectQuery(getByID).WillReturnError(sql.ErrNoRows),
		},
	}

	for i, tc := range tests {
		store := New()

		customer, err := store.GetByID(ctx, tc.id)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(customer, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, customer)

		}
	}
}

func TestCustomer_Create(t *testing.T) {
	ctx, mock := getTestData(t)
	dbErr := errors.Error("error from db")

	tests := []struct {
		desc      string
		input     model.Customer
		expOutput int
		expErr    error
		mockQuery *sqlmock.ExpectedExec
	}{
		{
			desc:      "success",
			input:     model.Customer{Name: "abc"},
			expOutput: 1,
			expErr:    nil,
			mockQuery: mock.ExpectExec(insert).WithArgs("abc").WillReturnResult(sqlmock.NewResult(1, 1)),
		},
		{
			desc:      "DB error",
			input:     model.Customer{Name: "pqr"},
			expOutput: 0,
			expErr:    errors.DB{Err: dbErr},
			mockQuery: mock.ExpectExec(insert).WithArgs("pqr").WillReturnError(dbErr),
		},
		//{
		//	desc:      "error from last insert id",
		//	input:     model.Customer{Name: "xyz"},
		//	expOutput: 0,
		//	expErr:    errors.Error("error from lastInsertID"),
		//	mockQuery: mock.ExpectExec(insert).WithArgs("xyz").WillReturnResult(sqlmock.NewResult(1, 1)),
		//},
	}

	for i, tc := range tests {
		store := New()

		id, err := store.Create(ctx, &tc.input)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(id, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, id)

		}
	}
}

func TestCustomer_Update(t *testing.T) {
	ctx, mock := getTestData(t)
	dbErr := errors.Error("error from db")

	tests := []struct {
		desc      string
		input     model.Customer
		expErr    error
		mockQuery *sqlmock.ExpectedExec
	}{
		{
			desc:      "Successfully update",
			input:     model.Customer{ID: 1, Name: "abc"},
			expErr:    nil,
			mockQuery: mock.ExpectExec(update).WithArgs("abc", 1).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			desc:      "DB error",
			input:     model.Customer{ID: 2, Name: "pqr"},
			expErr:    errors.DB{Err: dbErr},
			mockQuery: mock.ExpectExec(update).WithArgs("pqr", 2).WillReturnError(dbErr),
		},
	}

	for i, tc := range tests {
		store := New()

		err := store.Update(ctx, &tc.input)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}
	}
}

func TestCustomer_Delete(t *testing.T) {
	ctx, mock := getTestData(t)
	dbErr := errors.Error("error from db")

	tests := []struct {
		desc      string
		input     int
		expErr    error
		mockQuery *sqlmock.ExpectedExec
	}{
		{
			desc:      "successfully deleted",
			input:     1,
			expErr:    nil,
			mockQuery: mock.ExpectExec(delete).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1)),
		},
		{
			desc:      "BD error",
			input:     2,
			expErr:    errors.DB{Err: dbErr},
			mockQuery: mock.ExpectExec(delete).WithArgs(2).WillReturnError(dbErr),
		},
	}

	for i, tc := range tests {
		store := New()

		err := store.Delete(ctx, tc.input)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}
	}
}
