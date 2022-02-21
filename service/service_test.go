package service

import (
	"context"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/gofr-crud/model"
	"github.com/gofr-crud/store"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestService_Create(t *testing.T) {
	app := gofr.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	service := New(mockStore)

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	tests := []struct {
		desc      string
		input     model.Customer
		expOutput model.Customer
		expErr    error
		mockCall  []*gomock.Call
	}{
		{
			desc:      "success",
			input:     model.Customer{Name: "abc"},
			expOutput: model.Customer{ID: 1, Name: "abc"},
			expErr:    nil,
			mockCall: []*gomock.Call{
				mockStore.EXPECT().Create(ctx, &model.Customer{Name: "abc"}).Return(1, nil),
				mockStore.EXPECT().GetByID(ctx, 1).Return(model.Customer{ID: 1, Name: "abc"}, nil),
			},
		},
		{
			desc:      "Invalid param",
			input:     model.Customer{Name: ""},
			expOutput: model.Customer{},
			expErr:    errors.InvalidParam{Param: []string{"Name"}},
			mockCall:  []*gomock.Call{},
		},
		{
			desc:      "error from store layer",
			input:     model.Customer{Name: "pqr"},
			expOutput: model.Customer{},
			expErr:    errors.Error("error from create"),
			mockCall: []*gomock.Call{
				mockStore.EXPECT().Create(ctx, &model.Customer{Name: "pqr"}).Return(0, errors.Error("error from create")),
			},
		},
		{
			desc:      "error from getByID",
			input:     model.Customer{Name: "xyz"},
			expOutput: model.Customer{},
			expErr:    errors.Error("error from getByID"),
			mockCall: []*gomock.Call{
				mockStore.EXPECT().Create(ctx, &model.Customer{Name: "xyz"}).Return(3, nil),
				mockStore.EXPECT().GetByID(ctx, 3).Return(model.Customer{}, errors.Error("error from getByID")),
			},
		},
	}

	for i, tc := range tests {

		cust, err := service.Create(ctx, &tc.input)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(cust, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, cust)

		}
	}
}

func TestService_Get(t *testing.T) {
	app := gofr.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	service := New(mockStore)

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	tests := []struct {
		desc      string
		expOutput []model.Customer
		expErr    error
		mockCall  *gomock.Call
	}{
		{
			desc:      "Fetch all customers",
			expOutput: []model.Customer{{ID: 1, Name: "abc"}},
			expErr:    nil,
			mockCall:  mockStore.EXPECT().Get(ctx).Return([]model.Customer{{ID: 1, Name: "abc"}}, nil),
		},
		{
			desc:      "error from store layer",
			expOutput: nil,
			expErr:    errors.Error("error from store layer"),
			mockCall:  mockStore.EXPECT().Get(ctx).Return(nil, errors.Error("error from store layer")),
		},
	}

	for i, tc := range tests {
		customers, err := service.Get(ctx)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(customers, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, customers)

		}

	}
}

func TestService_GetByID(t *testing.T) {
	app := gofr.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	service := New(mockStore)

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	tests := []struct {
		desc      string
		input     string
		expOutput model.Customer
		expErr    error
		mockCall  *gomock.Call
	}{
		{
			desc:      "Success fro id 1",
			input:     "1",
			expOutput: model.Customer{ID: 1, Name: "abc"},
			expErr:    nil,
			mockCall:  mockStore.EXPECT().GetByID(ctx, 1).Return(model.Customer{ID: 1, Name: "abc"}, nil),
		},
		{
			desc:      "Invalid param",
			input:     "a",
			expOutput: model.Customer{},
			expErr:    errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:      "Invalid param",
			input:     "-1",
			expOutput: model.Customer{},
			expErr:    errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:      "error from store layer",
			input:     "2",
			expOutput: model.Customer{},
			expErr:    errors.Error("entity not found"),
			mockCall:  mockStore.EXPECT().GetByID(ctx, 2).Return(model.Customer{}, errors.Error("entity not found")),
		},
	}

	for i, tc := range tests {
		customers, err := service.GetByID(ctx, tc.input)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(customers, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, customers)

		}
	}
}

func TestService_Update(t *testing.T) {
	app := gofr.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	service := New(mockStore)

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	tests := []struct {
		desc      string
		input     model.Customer
		expOutput model.Customer
		expErr    error
		mockCall  []*gomock.Call
	}{
		{
			desc:      "successfully update",
			input:     model.Customer{ID: 1, Name: "abc"},
			expOutput: model.Customer{ID: 1, Name: "abc"},
			expErr:    nil,
			mockCall: []*gomock.Call{
				mockStore.EXPECT().Update(ctx, &model.Customer{ID: 1, Name: "abc"}).Return(nil),
				mockStore.EXPECT().GetByID(ctx, 1).Return(model.Customer{ID: 1, Name: "abc"}, nil),
			},
		},
		{
			desc:      "Invalid param",
			input:     model.Customer{ID: 2, Name: ""},
			expOutput: model.Customer{},
			expErr:    errors.InvalidParam{Param: []string{"Name"}},
			mockCall:  []*gomock.Call{},
		},
		{
			desc:      "Invalid param",
			input:     model.Customer{ID: -2, Name: "abc"},
			expOutput: model.Customer{},
			expErr:    errors.InvalidParam{Param: []string{"id"}},
			mockCall:  []*gomock.Call{},
		},
		{
			desc:      "error from store layer",
			input:     model.Customer{ID: 3, Name: "pqr"},
			expOutput: model.Customer{},
			expErr:    errors.Error("error from update"),
			mockCall: []*gomock.Call{
				mockStore.EXPECT().Update(ctx, &model.Customer{ID: 3, Name: "pqr"}).Return(errors.Error("error from update")),
			},
		},
		{
			desc:      "error from getByID",
			input:     model.Customer{ID: 4, Name: "xyz"},
			expOutput: model.Customer{},
			expErr:    errors.Error("error from getByID"),
			mockCall: []*gomock.Call{
				mockStore.EXPECT().Update(ctx, &model.Customer{ID: 4, Name: "xyz"}).Return(nil),
				mockStore.EXPECT().GetByID(ctx, 4).Return(model.Customer{}, errors.Error("error from getByID")),
			},
		},
	}

	for i, tc := range tests {
		customers, err := service.Update(ctx, &tc.input)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(customers, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, customers)

		}
	}
}

func TestService_Delete(t *testing.T) {
	app := gofr.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	service := New(mockStore)

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	tests := []struct {
		desc     string
		input    string
		expErr   error
		mockCall *gomock.Call
	}{
		{
			desc:     "Successfully deleted",
			input:    "1",
			expErr:   nil,
			mockCall: mockStore.EXPECT().Delete(ctx, 1).Return(nil),
		},
		{
			desc:   "Invalid param",
			input:  "a",
			expErr: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:   "Invalid param",
			input:  "-1",
			expErr: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:     "error from store layer",
			input:    "2",
			expErr:   errors.Error("error from store layer"),
			mockCall: mockStore.EXPECT().Delete(ctx, 2).Return(errors.Error("error from store layer")),
		},
	}

	for i, tc := range tests {
		err := service.Delete(ctx, tc.input)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}
	}
}
