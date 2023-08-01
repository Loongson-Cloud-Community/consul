// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	cloud "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
)

// HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest hashicorp cloud global network manager 20220215 agent push server state request
//
// swagger:model hashicorp.cloud.global_network_manager_20220215.AgentPushServerStateRequest
type HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest struct {

	// id is the name of the cluster
	ID string `json:"id,omitempty"`

	// location is the project and organization of the cluster with an optional provider and region
	Location *cloud.HashicorpCloudLocationLocation `json:"location,omitempty"`

	// server_state is the internal consul node information
	ServerState *HashicorpCloudGlobalNetworkManager20220215ServerState `json:"server_state,omitempty"`
}

// Validate validates this hashicorp cloud global network manager 20220215 agent push server state request
func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateServerState(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) validateLocation(formats strfmt.Registry) error {
	if swag.IsZero(m.Location) { // not required
		return nil
	}

	if m.Location != nil {
		if err := m.Location.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("location")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("location")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) validateServerState(formats strfmt.Registry) error {
	if swag.IsZero(m.ServerState) { // not required
		return nil
	}

	if m.ServerState != nil {
		if err := m.ServerState.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("server_state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("server_state")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp cloud global network manager 20220215 agent push server state request based on the context it is used
func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLocation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateServerState(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) contextValidateLocation(ctx context.Context, formats strfmt.Registry) error {

	if m.Location != nil {
		if err := m.Location.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("location")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("location")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) contextValidateServerState(ctx context.Context, formats strfmt.Registry) error {

	if m.ServerState != nil {
		if err := m.ServerState.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("server_state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("server_state")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudGlobalNetworkManager20220215AgentPushServerStateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}