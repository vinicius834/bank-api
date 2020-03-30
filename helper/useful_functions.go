package helper

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-playground/validator"
)

func ErrorsExist(errs []error) bool {
	return errs != nil && len(errs) > 0
}

func ProcessRequest(r http.Handler, method, url, bodyData string) *httptest.ResponseRecorder {
	var req *http.Request
	if bodyData == "" {
		req, _ = http.NewRequest(method, url, nil)
	} else {
		req, _ = http.NewRequest(method, url, bytes.NewBufferString(bodyData))
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ValidateFields(item interface{}) []error {
	var invalidFields []error
	invalidFieldMsg := "invalid"
	validate := validator.New()
	problems := validate.Struct(item)
	if problems != nil {
		for _, err := range problems.(validator.ValidationErrors) {
			msg := fmt.Sprintf("%v field %v", err.Field(), invalidFieldMsg)
			invalidFields = append(invalidFields, errors.New(msg))
		}
	}
	return invalidFields
}
