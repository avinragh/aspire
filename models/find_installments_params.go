// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// FindInstallmentsParams find installments params
//
// swagger:model FindInstallmentsParams
type FindInstallmentsParams struct {

	// limit
	Limit string `json:"limit,omitempty"`

	// loan Id
	LoanID string `json:"loanId,omitempty"`

	// page
	Page string `json:"page,omitempty"`

	// sort
	Sort string `json:"sort,omitempty"`

	// state
	State string `json:"state,omitempty"`
}

// Validate validates this find installments params
func (m *FindInstallmentsParams) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this find installments params based on context it is used
func (m *FindInstallmentsParams) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FindInstallmentsParams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FindInstallmentsParams) UnmarshalBinary(b []byte) error {
	var res FindInstallmentsParams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}