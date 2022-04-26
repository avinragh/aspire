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

// Loan loan
//
// swagger:model Loan
type Loan struct {

	// amount
	// Required: true
	Amount *int64 `json:"amount"`

	// currency
	Currency string `json:"currency,omitempty"`

	// id
	// Required: true
	ID *int64 `json:"id"`

	// state
	State string `json:"state,omitempty"`

	// term
	// Required: true
	Term *int64 `json:"term"`

	// user Id
	UserID int64 `json:"userId,omitempty"`
}

// Validate validates this loan
func (m *Loan) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAmount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTerm(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Loan) validateAmount(formats strfmt.Registry) error {

	if err := validate.Required("amount", "body", m.Amount); err != nil {
		return err
	}

	return nil
}

func (m *Loan) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *Loan) validateTerm(formats strfmt.Registry) error {

	if err := validate.Required("term", "body", m.Term); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this loan based on context it is used
func (m *Loan) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Loan) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Loan) UnmarshalBinary(b []byte) error {
	var res Loan
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}