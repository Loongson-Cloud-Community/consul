// Code generated by protoc-gen-go-binary. DO NOT EDIT.
// source: pbmesh/v1alpha1/address.proto

package meshv1alpha1

import (
	"google.golang.org/protobuf/proto"
)

// MarshalBinary implements encoding.BinaryMarshaler
func (msg *HostPortAddress) MarshalBinary() ([]byte, error) {
	return proto.Marshal(msg)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (msg *HostPortAddress) UnmarshalBinary(b []byte) error {
	return proto.Unmarshal(b, msg)
}
