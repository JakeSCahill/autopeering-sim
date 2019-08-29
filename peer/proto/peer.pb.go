// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peer/proto/peer.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Minimal encoding of a peer
type Peer struct {
	// public key used for signing
	PublicKey []byte `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	// address of autopeering protocol (e.g. "192.0.2.1:25", "[2001:db8::1]:80")
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_155860cd2f47eba7, []int{0}
}

func (m *Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer.Unmarshal(m, b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
}
func (m *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(m, src)
}
func (m *Peer) XXX_Size() int {
	return xxx_messageInfo_Peer.Size(m)
}
func (m *Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Peer proto.InternalMessageInfo

func (m *Peer) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *Peer) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*Peer)(nil), "proto.Peer")
}

func init() { proto.RegisterFile("peer/proto/peer.proto", fileDescriptor_155860cd2f47eba7) }

var fileDescriptor_155860cd2f47eba7 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x48, 0x4d, 0x2d,
	0xd2, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x07, 0x31, 0xf5, 0xc0, 0x4c, 0x21, 0x56, 0x30, 0xa5,
	0x64, 0xcf, 0xc5, 0x12, 0x90, 0x9a, 0x5a, 0x24, 0x24, 0xcb, 0xc5, 0x55, 0x50, 0x9a, 0x94, 0x93,
	0x99, 0x1c, 0x9f, 0x9d, 0x5a, 0x29, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x13, 0xc4, 0x09, 0x11, 0xf1,
	0x4e, 0xad, 0x14, 0x92, 0xe0, 0x62, 0x4f, 0x4c, 0x49, 0x29, 0x4a, 0x2d, 0x2e, 0x96, 0x60, 0x52,
	0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x71, 0x9d, 0xb4, 0xa2, 0x34, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93,
	0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xcb, 0xf3, 0x73, 0x72, 0x12, 0x93, 0xf5, 0x13, 0x4b, 0x4b, 0xf2,
	0x41, 0x76, 0x65, 0xe6, 0xa5, 0xeb, 0x23, 0xac, 0x4f, 0x62, 0x03, 0x53, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xf4, 0x93, 0x65, 0xc5, 0x93, 0x00, 0x00, 0x00,
}
