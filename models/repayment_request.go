// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RepaymentRequest repayment request
//
// swagger:model RepaymentRequest
type RepaymentRequest struct {

	// repayment amount
	// Required: true
	RepaymentAmount *float64 `json:"repaymentAmount"`
}

// Validate validates this repayment request
func (m *RepaymentRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRepaymentAmount(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RepaymentRequest) validateRepaymentAmount(formats strfmt.Registry) error {

	if err := validate.Required("repaymentAmount", "body", m.RepaymentAmount); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this repayment request based on context it is used
func (m *RepaymentRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RepaymentRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RepaymentRequest) UnmarshalBinary(b []byte) error {
	var res RepaymentRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}