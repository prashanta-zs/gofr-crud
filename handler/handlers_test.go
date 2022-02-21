package handler

import (
	"context"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/gofr-crud/model"
	"github.com/gofr-crud/service"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_Get(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockService(ctrl)
	handler := New(mockService)

	tests := []struct {
		desc      string
		expOutput interface{}
		expErr    error
		mockCall  *gomock.Call
	}{
		{
			desc:      "Success case",
			expOutput: []model.Customer{{ID: 1, Name: "abc"}},
			expErr:    nil,
			mockCall:  mockService.EXPECT().Get(ctx).Return([]model.Customer{{ID: 1, Name: "abc"}}, nil),
		},
		{
			desc:      "Failure case",
			expOutput: nil,
			expErr:    errors.Error("error from service layer"),
			mockCall:  mockService.EXPECT().Get(ctx).Return([]model.Customer{}, errors.Error("error from service layer")),
		},
	}

	for i, tc := range tests {
		customer, err := handler.Get(ctx)

		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(customer, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, customer)

		}
	}
}

func TestHandler_GetByID(t *testing.T) {
	app := gofr.New()

	req := httptest.NewRequest(http.MethodGet, "/customer/{id}", nil)
	res := httptest.NewRecorder()

	r := request.NewHTTPRequest(req)
	w := responder.NewContextualResponder(res, req)

	ctx := gofr.NewContext(w, r, app)
	ctx.Context = context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockService(ctrl)
	handler := New(mockService)

	tests := []struct {
		desc      string
		id        string
		expOutput interface{}
		expErr    error
		mockCall  *gomock.Call
	}{
		{
			desc:      "success case",
			id:        "1",
			expOutput: model.Customer{ID: 1, Name: "abc"},
			expErr:    nil,
			mockCall:  mockService.EXPECT().GetByID(ctx, "1").Return(model.Customer{ID: 1, Name: "abc"}, nil),
		},
		{
			desc:      "missing param",
			id:        "",
			expOutput: nil,
			expErr:    errors.MissingParam{Param: []string{"id"}},
		},
		{
			desc:      "error from service layer",
			id:        "2",
			expOutput: nil,
			expErr:    errors.Error("error from service layer"),
			mockCall:  mockService.EXPECT().GetByID(ctx, "2").Return(model.Customer{}, errors.Error("error from service layer")),
		},
	}

	for i, tc := range tests {
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})
		customer, err := handler.GetByID(ctx)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected error: %v\nGot: %v", i, tc.desc, tc.expErr, err)
		}

		if !reflect.DeepEqual(customer, tc.expOutput) {
			t.Errorf("Testcase[%v] failed (%v)\nExpected:\n%v\nGot:\n%v", i, tc.desc, tc.expOutput, customer)

		}
	}
}
