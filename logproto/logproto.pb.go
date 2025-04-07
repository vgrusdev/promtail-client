// Code generated by protoc-gen-go. DO NOT EDIT.
// source: logproto.proto

package logproto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type PushRequest struct {
	Streams              []*Stream `protobuf:"bytes,1,rep,name=streams,proto3" json:"streams,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PushRequest) Reset()         { *m = PushRequest{} }
func (m *PushRequest) String() string { return proto.CompactTextString(m) }
func (*PushRequest) ProtoMessage()    {}
func (*PushRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7a8976f235a02f79, []int{0}
}

func (m *PushRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushRequest.Unmarshal(m, b)
}
func (m *PushRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushRequest.Marshal(b, m, deterministic)
}
func (m *PushRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushRequest.Merge(m, src)
}
func (m *PushRequest) XXX_Size() int {
	return xxx_messageInfo_PushRequest.Size(m)
}
func (m *PushRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PushRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PushRequest proto.InternalMessageInfo

func (m *PushRequest) GetStreams() []*Stream {
	if m != nil {
		return m.Streams
	}
	return nil
}

type Stream struct {
	//Labels               string   `protobuf:"bytes,1,opt,name=labels,proto3" json:"labels,omitempty"`
	Labels               map[string]string   `protobuf:"bytes,1,opt,name=labels,proto3" json:"labels,omitempty"`
	Entries              []*Entry             `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Stream) Reset()         { *m = Stream{} }
func (m *Stream) String() string { return proto.CompactTextString(m) }
func (*Stream) ProtoMessage()    {}
func (*Stream) Descriptor() ([]byte, []int) {
	return fileDescriptor_7a8976f235a02f79, []int{1}
}

func (m *Stream) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stream.Unmarshal(m, b)
}
func (m *Stream) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stream.Marshal(b, m, deterministic)
}
func (m *Stream) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stream.Merge(m, src)
}
func (m *Stream) XXX_Size() int {
	return xxx_messageInfo_Stream.Size(m)
}
func (m *Stream) XXX_DiscardUnknown() {
	xxx_messageInfo_Stream.DiscardUnknown(m)
}

var xxx_messageInfo_Stream proto.InternalMessageInfo

func (m *Stream) GetLabels() string {
	if m != nil {
		return m.Labels
	}
	return ""
}

func (m *Stream) GetEntries() []*Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type Entry struct {
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Line                 string               `protobuf:"bytes,2,opt,name=line,proto3" json:"line,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_7a8976f235a02f79, []int{2}
}

func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Entry) GetLine() string {
	if m != nil {
		return m.Line
	}
	return ""
}

func init() {
	proto.RegisterType((*PushRequest)(nil), "logproto.PushRequest")
	proto.RegisterType((*Stream)(nil), "logproto.Stream")
	proto.RegisterType((*Entry)(nil), "logproto.Entry")
}

func init() { proto.RegisterFile("logproto.proto", fileDescriptor_7a8976f235a02f79) }

var fileDescriptor_7a8976f235a02f79 = []byte{
	// 208 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0xc9, 0x4f, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x03, 0x93, 0x42, 0x1c, 0x30, 0xbe, 0x94, 0x7c, 0x7a, 0x7e, 0x7e,
	0x7a, 0x4e, 0xaa, 0x3e, 0x98, 0x97, 0x54, 0x9a, 0xa6, 0x5f, 0x92, 0x99, 0x9b, 0x5a, 0x5c, 0x92,
	0x98, 0x5b, 0x00, 0x51, 0xaa, 0x64, 0xc9, 0xc5, 0x1d, 0x50, 0x5a, 0x9c, 0x11, 0x94, 0x5a, 0x58,
	0x9a, 0x5a, 0x5c, 0x22, 0xa4, 0xc5, 0xc5, 0x5e, 0x5c, 0x52, 0x94, 0x9a, 0x98, 0x5b, 0x2c, 0xc1,
	0xa8, 0xc0, 0xac, 0xc1, 0x6d, 0x24, 0xa0, 0x07, 0x37, 0x3b, 0x18, 0x2c, 0x11, 0x04, 0x53, 0xa0,
	0xe4, 0xcd, 0xc5, 0x06, 0x11, 0x12, 0x12, 0xe3, 0x62, 0xcb, 0x49, 0x4c, 0x4a, 0xcd, 0x01, 0x69,
	0x62, 0xd4, 0xe0, 0x0c, 0x82, 0xf2, 0x84, 0x34, 0xb9, 0xd8, 0x53, 0xf3, 0x4a, 0x8a, 0x32, 0x53,
	0x8b, 0x25, 0x98, 0xc0, 0xa6, 0xf1, 0x23, 0x4c, 0x73, 0xcd, 0x2b, 0x29, 0xaa, 0x0c, 0x82, 0xc9,
	0x2b, 0x85, 0x72, 0xb1, 0x82, 0x45, 0x84, 0x2c, 0xb8, 0x38, 0xe1, 0x6e, 0x04, 0x1b, 0xc7, 0x6d,
	0x24, 0xa5, 0x07, 0xf1, 0x85, 0x1e, 0xcc, 0x17, 0x7a, 0x21, 0x30, 0x15, 0x41, 0x08, 0xc5, 0x42,
	0x42, 0x5c, 0x2c, 0x39, 0x99, 0x79, 0xa9, 0x12, 0x4c, 0x60, 0x37, 0x80, 0xd9, 0x49, 0x6c, 0x60,
	0x2d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x63, 0xb9, 0x3c, 0x7e, 0x22, 0x01, 0x00, 0x00,
}
