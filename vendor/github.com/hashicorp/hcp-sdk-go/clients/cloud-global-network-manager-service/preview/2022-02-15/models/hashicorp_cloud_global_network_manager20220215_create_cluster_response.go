// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse hashicorp cloud global network manager 20220215 create cluster response
//
// swagger:model hashicorp.cloud.global_network_manager_20220215.CreateClusterResponse
type HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse struct {

	// client_id is the oauth client_id used to authenticate consul to global-network-manager. This _may_ move into a separate secrets RPC later on.
	ClientID string `json:"client_id,omitempty"`

	// client_secret is the oauth client_secret used to authenticate consul to global-network-manager. This _may_ move into a separate secrets RPC later on.
	ClientSecret string `json:"client_secret,omitempty"`

	// cluster is the viewable representation of the cluster.
	Cluster *HashicorpCloudGlobalNetworkManager20220215Cluster `json:"cluster,omitempty"`
}

// Validate validates this hashicorp cloud global network manager 20220215 create cluster response
func (m *HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCluster(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse) validateCluster(formats strfmt.Registry) error {
	if swag.IsZero(m.Cluster) { // not required
		return nil
	}

	if m.Cluster != nil {
		if err := m.Cluster.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cluster")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("cluster")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp cloud global network manager 20220215 create cluster response based on the context it is used
func (m *HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCluster(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse) contextValidateCluster(ctx context.Context, formats strfmt.Registry) error {

	if m.Cluster != nil {
		if err := m.Cluster.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cluster")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("cluster")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpCloudGlobalNetworkManager20220215CreateClusterResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}