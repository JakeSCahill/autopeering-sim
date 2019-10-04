// Code generated by protoc-gen-go. DO NOT EDIT.
// source: neighborhood/proto/message.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	proto1 "github.com/wollac/autopeering/salt/proto"
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

type PeeringRequest struct {
	// string form of the recipient address
	To string `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	// unix time
	Timestamp int64 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// salt of the requester
	Salt                 *proto1.Salt `protobuf:"bytes,3,opt,name=salt,proto3" json:"salt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *PeeringRequest) Reset()         { *m = PeeringRequest{} }
func (m *PeeringRequest) String() string { return proto.CompactTextString(m) }
func (*PeeringRequest) ProtoMessage()    {}
func (*PeeringRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_150a5debc6b6dd9c, []int{0}
}

func (m *PeeringRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeeringRequest.Unmarshal(m, b)
}
func (m *PeeringRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeeringRequest.Marshal(b, m, deterministic)
}
func (m *PeeringRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeeringRequest.Merge(m, src)
}
func (m *PeeringRequest) XXX_Size() int {
	return xxx_messageInfo_PeeringRequest.Size(m)
}
func (m *PeeringRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PeeringRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PeeringRequest proto.InternalMessageInfo

func (m *PeeringRequest) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *PeeringRequest) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *PeeringRequest) GetSalt() *proto1.Salt {
	if m != nil {
		return m.Salt
	}
	return nil
}

type PeeringResponse struct {
	// hash of the corresponding request
	ReqHash []byte `protobuf:"bytes,1,opt,name=req_hash,json=reqHash,proto3" json:"req_hash,omitempty"`
	// response of a peering request
	Status               bool     `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeeringResponse) Reset()         { *m = PeeringResponse{} }
func (m *PeeringResponse) String() string { return proto.CompactTextString(m) }
func (*PeeringResponse) ProtoMessage()    {}
func (*PeeringResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_150a5debc6b6dd9c, []int{1}
}

func (m *PeeringResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeeringResponse.Unmarshal(m, b)
}
func (m *PeeringResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeeringResponse.Marshal(b, m, deterministic)
}
func (m *PeeringResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeeringResponse.Merge(m, src)
}
func (m *PeeringResponse) XXX_Size() int {
	return xxx_messageInfo_PeeringResponse.Size(m)
}
func (m *PeeringResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PeeringResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PeeringResponse proto.InternalMessageInfo

func (m *PeeringResponse) GetReqHash() []byte {
	if m != nil {
		return m.ReqHash
	}
	return nil
}

func (m *PeeringResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

type PeeringDrop struct {
	// string form of the recipient address
	To string `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	// unix time
	Timestamp            int64    `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeeringDrop) Reset()         { *m = PeeringDrop{} }
func (m *PeeringDrop) String() string { return proto.CompactTextString(m) }
func (*PeeringDrop) ProtoMessage()    {}
func (*PeeringDrop) Descriptor() ([]byte, []int) {
	return fileDescriptor_150a5debc6b6dd9c, []int{2}
}

func (m *PeeringDrop) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeeringDrop.Unmarshal(m, b)
}
func (m *PeeringDrop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeeringDrop.Marshal(b, m, deterministic)
}
func (m *PeeringDrop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeeringDrop.Merge(m, src)
}
func (m *PeeringDrop) XXX_Size() int {
	return xxx_messageInfo_PeeringDrop.Size(m)
}
func (m *PeeringDrop) XXX_DiscardUnknown() {
	xxx_messageInfo_PeeringDrop.DiscardUnknown(m)
}

var xxx_messageInfo_PeeringDrop proto.InternalMessageInfo

func (m *PeeringDrop) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *PeeringDrop) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*PeeringRequest)(nil), "proto.PeeringRequest")
	proto.RegisterType((*PeeringResponse)(nil), "proto.PeeringResponse")
	proto.RegisterType((*PeeringDrop)(nil), "proto.PeeringDrop")
}

func init() { proto.RegisterFile("neighborhood/proto/message.proto", fileDescriptor_150a5debc6b6dd9c) }

var fileDescriptor_150a5debc6b6dd9c = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0xd5, 0x16, 0x4a, 0x7b, 0x41, 0x45, 0xb2, 0x04, 0x0a, 0x08, 0x89, 0x2a, 0x13, 0x53,
	0x8c, 0xca, 0xc8, 0x86, 0x3a, 0x30, 0x22, 0xb3, 0xb1, 0x54, 0x4e, 0x39, 0xc5, 0x91, 0x92, 0x9c,
	0x6b, 0x9f, 0xc5, 0xdf, 0xc7, 0x75, 0x22, 0x18, 0x98, 0x3a, 0xdd, 0xbd, 0xf7, 0xac, 0xef, 0x59,
	0x07, 0xeb, 0x1e, 0x9b, 0xda, 0x54, 0xe4, 0x0c, 0xd1, 0x97, 0xb4, 0x8e, 0x98, 0x64, 0x87, 0xde,
	0xeb, 0x1a, 0xcb, 0xa4, 0xc4, 0x79, 0x1a, 0x77, 0xd7, 0x5e, 0xb7, 0x3c, 0x3e, 0x38, 0xae, 0x43,
	0x5a, 0xec, 0x60, 0xf5, 0x8e, 0xe8, 0x9a, 0xbe, 0x56, 0x78, 0x08, 0xe8, 0x59, 0xac, 0x60, 0xca,
	0x94, 0x4f, 0xd6, 0x93, 0xc7, 0xa5, 0x8a, 0x9b, 0xb8, 0x87, 0x25, 0x37, 0x11, 0xc9, 0xba, 0xb3,
	0xf9, 0x34, 0xda, 0x33, 0xf5, 0x67, 0x88, 0x07, 0x38, 0x3b, 0xd2, 0xf2, 0x59, 0x0c, 0xb2, 0x4d,
	0x36, 0x50, 0xcb, 0x8f, 0x68, 0xa9, 0x14, 0x14, 0x5b, 0xb8, 0xfa, 0x2d, 0xf0, 0x96, 0x7a, 0x8f,
	0xe2, 0x16, 0x16, 0x0e, 0x0f, 0x3b, 0xa3, 0xbd, 0x49, 0x3d, 0x97, 0xea, 0x22, 0xea, 0xb7, 0x28,
	0xc5, 0x0d, 0xcc, 0x23, 0x97, 0x83, 0x4f, 0x4d, 0x0b, 0x35, 0xaa, 0xe2, 0x05, 0xb2, 0x91, 0xb2,
	0x75, 0x64, 0x4f, 0xfb, 0xe3, 0xeb, 0xe6, 0xf3, 0xa9, 0x6e, 0xd8, 0x84, 0xaa, 0xdc, 0x53, 0x27,
	0xbf, 0xa9, 0x6d, 0xf5, 0x5e, 0xea, 0xc0, 0x64, 0x07, 0xa4, 0xfc, 0x7f, 0xc3, 0x6a, 0x9e, 0xc6,
	0xf3, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8f, 0x29, 0xf8, 0x7c, 0x60, 0x01, 0x00, 0x00,
}
