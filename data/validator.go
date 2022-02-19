package data

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"net/http"
)

type Validator struct {
	validator *validator.Validate
}

// NewValidator inits a new validator with an underlying
// OSS struct and value validator
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate func validates input by delegating the struct
// validation work to the underlying validator
func (v *Validator) Validate(i interface{}) error {
	flowLogSlice := i.([]FlowLog)
	for _, log := range flowLogSlice {
		if err := v.validator.Struct(log); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	return nil
}