// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gateway.proto

package apis

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// EndpointType defines supported endpoints' types
type EndpointType int32

const (
	EndpointType_UNKNOWN EndpointType = 0
	EndpointType_IPV4    EndpointType = 1
	EndpointType_IPV6    EndpointType = 2
	EndpointType_SID     EndpointType = 3
)

var EndpointType_name = map[int32]string{
	0: "UNKNOWN",
	1: "IPV4",
	2: "IPV6",
	3: "SID",
}

var EndpointType_value = map[string]int32{
	"UNKNOWN": 0,
	"IPV4":    1,
	"IPV6":    2,
	"SID":     3,
}

func (x EndpointType) String() string {
	return proto.EnumName(EndpointType_name, int32(x))
}

func (EndpointType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{0}
}

// Errors defines error codes which a gateway can return to a client
type GatewayErrors int32

const (
	GatewayErrors_OK GatewayErrors = 0
	// Returned when a gateway encountered a generic error
	GatewayErrors_EIO GatewayErrors = 1
	// Returned when request key does not exist
	GatewayErrors_ENOENT GatewayErrors = 2
	// Returned when an operation triggered by the client's requested timed out
	// and was canceled
	GatewayErrors_ETIMEDOUT GatewayErrors = 3
	// Returned when a gateway cannot reach a DB host
	GatewayErrors_EHOSTDOWN GatewayErrors = 4
)

var GatewayErrors_name = map[int32]string{
	0: "OK",
	1: "EIO",
	2: "ENOENT",
	3: "ETIMEDOUT",
	4: "EHOSTDOWN",
}

var GatewayErrors_value = map[string]int32{
	"OK":        0,
	"EIO":       1,
	"ENOENT":    2,
	"ETIMEDOUT": 3,
	"EHOSTDOWN": 4,
}

func (x GatewayErrors) String() string {
	return proto.EnumName(GatewayErrors_name, int32(x))
}

func (GatewayErrors) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{1}
}

// endpoint is an object defining a source or a destination of a communication path.
// The type of the endpoint is defined by EndpointType, it can be ipv4, ipv6 addresses or segment routing SIDs etc.
// The list of supported endpoint types may be extended in future.
// The address of endpoint is stored in a slice of bytes and should be decoded according to the endpoint type.
type Endpoint struct {
	Type                 EndpointType `protobuf:"varint,1,opt,name=type,proto3,enum=apis.EndpointType" json:"type,omitempty"`
	Address              []byte       `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{0}
}

func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (m *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(m, src)
}
func (m *Endpoint) XXX_Size() int {
	return xxx_messageInfo_Endpoint.Size(m)
}
func (m *Endpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoint proto.InternalMessageInfo

func (m *Endpoint) GetType() EndpointType {
	if m != nil {
		return m.Type
	}
	return EndpointType_UNKNOWN
}

func (m *Endpoint) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

// latency defines one of available QoE metrics
// value expresses latency value in milliseconds.
// variation expresses percent of the acceptable variation from value.
type Latency struct {
	Value                int32    `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	Variation            int32    `protobuf:"varint,2,opt,name=variation,proto3" json:"variation,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Latency) Reset()         { *m = Latency{} }
func (m *Latency) String() string { return proto.CompactTextString(m) }
func (*Latency) ProtoMessage()    {}
func (*Latency) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{1}
}

func (m *Latency) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Latency.Unmarshal(m, b)
}
func (m *Latency) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Latency.Marshal(b, m, deterministic)
}
func (m *Latency) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Latency.Merge(m, src)
}
func (m *Latency) XXX_Size() int {
	return xxx_messageInfo_Latency.Size(m)
}
func (m *Latency) XXX_DiscardUnknown() {
	xxx_messageInfo_Latency.DiscardUnknown(m)
}

var xxx_messageInfo_Latency proto.InternalMessageInfo

func (m *Latency) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Latency) GetVariation() int32 {
	if m != nil {
		return m.Variation
	}
	return 0
}

// qoe_parameters defines a list of QoE parameters a client can request.
// Currently only latency is supported.
type QoeParameters struct {
	Latency              *Latency `protobuf:"bytes,1,opt,name=latency,proto3" json:"latency,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QoeParameters) Reset()         { *m = QoeParameters{} }
func (m *QoeParameters) String() string { return proto.CompactTextString(m) }
func (*QoeParameters) ProtoMessage()    {}
func (*QoeParameters) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{2}
}

func (m *QoeParameters) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QoeParameters.Unmarshal(m, b)
}
func (m *QoeParameters) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QoeParameters.Marshal(b, m, deterministic)
}
func (m *QoeParameters) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QoeParameters.Merge(m, src)
}
func (m *QoeParameters) XXX_Size() int {
	return xxx_messageInfo_QoeParameters.Size(m)
}
func (m *QoeParameters) XXX_DiscardUnknown() {
	xxx_messageInfo_QoeParameters.DiscardUnknown(m)
}

var xxx_messageInfo_QoeParameters proto.InternalMessageInfo

func (m *QoeParameters) GetLatency() *Latency {
	if m != nil {
		return m.Latency
	}
	return nil
}

type Qoe struct {
	Src                  *Endpoint      `protobuf:"bytes,1,opt,name=src,proto3" json:"src,omitempty"`
	Dst                  *Endpoint      `protobuf:"bytes,2,opt,name=dst,proto3" json:"dst,omitempty"`
	Qoe                  *QoeParameters `protobuf:"bytes,3,opt,name=qoe,proto3" json:"qoe,omitempty"`
	Label                []uint32       `protobuf:"varint,4,rep,packed,name=label,proto3" json:"label,omitempty"`
	Err                  GatewayErrors  `protobuf:"varint,5,opt,name=err,proto3,enum=apis.GatewayErrors" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Qoe) Reset()         { *m = Qoe{} }
func (m *Qoe) String() string { return proto.CompactTextString(m) }
func (*Qoe) ProtoMessage()    {}
func (*Qoe) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{3}
}

func (m *Qoe) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Qoe.Unmarshal(m, b)
}
func (m *Qoe) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Qoe.Marshal(b, m, deterministic)
}
func (m *Qoe) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Qoe.Merge(m, src)
}
func (m *Qoe) XXX_Size() int {
	return xxx_messageInfo_Qoe.Size(m)
}
func (m *Qoe) XXX_DiscardUnknown() {
	xxx_messageInfo_Qoe.DiscardUnknown(m)
}

var xxx_messageInfo_Qoe proto.InternalMessageInfo

func (m *Qoe) GetSrc() *Endpoint {
	if m != nil {
		return m.Src
	}
	return nil
}

func (m *Qoe) GetDst() *Endpoint {
	if m != nil {
		return m.Dst
	}
	return nil
}

func (m *Qoe) GetQoe() *QoeParameters {
	if m != nil {
		return m.Qoe
	}
	return nil
}

func (m *Qoe) GetLabel() []uint32 {
	if m != nil {
		return m.Label
	}
	return nil
}

func (m *Qoe) GetErr() GatewayErrors {
	if m != nil {
		return m.Err
	}
	return GatewayErrors_OK
}

// RequestQoE defines the rpc message sent by the client to the gateway.
// Multiple Src/Dst/QoE are supported in a single request.
type RequestQoE struct {
	Qoes                 map[int32]*Qoe `protobuf:"bytes,1,rep,name=qoes,proto3" json:"qoes,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RequestQoE) Reset()         { *m = RequestQoE{} }
func (m *RequestQoE) String() string { return proto.CompactTextString(m) }
func (*RequestQoE) ProtoMessage()    {}
func (*RequestQoE) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{4}
}

func (m *RequestQoE) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestQoE.Unmarshal(m, b)
}
func (m *RequestQoE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestQoE.Marshal(b, m, deterministic)
}
func (m *RequestQoE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestQoE.Merge(m, src)
}
func (m *RequestQoE) XXX_Size() int {
	return xxx_messageInfo_RequestQoE.Size(m)
}
func (m *RequestQoE) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestQoE.DiscardUnknown(m)
}

var xxx_messageInfo_RequestQoE proto.InternalMessageInfo

func (m *RequestQoE) GetQoes() map[int32]*Qoe {
	if m != nil {
		return m.Qoes
	}
	return nil
}

// ResponseQoE defines  the rpc message sent as a reply to the client.
// it is the same message as request, but the gateway populates labels and
// err for each qoe.
type ResponseQoE struct {
	Qoes                 map[int32]*Qoe `protobuf:"bytes,1,rep,name=qoes,proto3" json:"qoes,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ResponseQoE) Reset()         { *m = ResponseQoE{} }
func (m *ResponseQoE) String() string { return proto.CompactTextString(m) }
func (*ResponseQoE) ProtoMessage()    {}
func (*ResponseQoE) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{5}
}

func (m *ResponseQoE) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseQoE.Unmarshal(m, b)
}
func (m *ResponseQoE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseQoE.Marshal(b, m, deterministic)
}
func (m *ResponseQoE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseQoE.Merge(m, src)
}
func (m *ResponseQoE) XXX_Size() int {
	return xxx_messageInfo_ResponseQoE.Size(m)
}
func (m *ResponseQoE) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseQoE.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseQoE proto.InternalMessageInfo

func (m *ResponseQoE) GetQoes() map[int32]*Qoe {
	if m != nil {
		return m.Qoes
	}
	return nil
}

// Borrowed from gobgp project
//
// Route Distinguisher is expressed by 2 bytes Autonomous System
type RouteDistinguisherTwoOctetAS struct {
	Admin                uint32   `protobuf:"varint,1,opt,name=admin,proto3" json:"admin,omitempty"`
	Assigned             uint32   `protobuf:"varint,2,opt,name=assigned,proto3" json:"assigned,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RouteDistinguisherTwoOctetAS) Reset()         { *m = RouteDistinguisherTwoOctetAS{} }
func (m *RouteDistinguisherTwoOctetAS) String() string { return proto.CompactTextString(m) }
func (*RouteDistinguisherTwoOctetAS) ProtoMessage()    {}
func (*RouteDistinguisherTwoOctetAS) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{6}
}

func (m *RouteDistinguisherTwoOctetAS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteDistinguisherTwoOctetAS.Unmarshal(m, b)
}
func (m *RouteDistinguisherTwoOctetAS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteDistinguisherTwoOctetAS.Marshal(b, m, deterministic)
}
func (m *RouteDistinguisherTwoOctetAS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteDistinguisherTwoOctetAS.Merge(m, src)
}
func (m *RouteDistinguisherTwoOctetAS) XXX_Size() int {
	return xxx_messageInfo_RouteDistinguisherTwoOctetAS.Size(m)
}
func (m *RouteDistinguisherTwoOctetAS) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteDistinguisherTwoOctetAS.DiscardUnknown(m)
}

var xxx_messageInfo_RouteDistinguisherTwoOctetAS proto.InternalMessageInfo

func (m *RouteDistinguisherTwoOctetAS) GetAdmin() uint32 {
	if m != nil {
		return m.Admin
	}
	return 0
}

func (m *RouteDistinguisherTwoOctetAS) GetAssigned() uint32 {
	if m != nil {
		return m.Assigned
	}
	return 0
}

// Route Distinguisher is expressed by IP address
type RouteDistinguisherIPAddress struct {
	Admin                string   `protobuf:"bytes,1,opt,name=admin,proto3" json:"admin,omitempty"`
	Assigned             uint32   `protobuf:"varint,2,opt,name=assigned,proto3" json:"assigned,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RouteDistinguisherIPAddress) Reset()         { *m = RouteDistinguisherIPAddress{} }
func (m *RouteDistinguisherIPAddress) String() string { return proto.CompactTextString(m) }
func (*RouteDistinguisherIPAddress) ProtoMessage()    {}
func (*RouteDistinguisherIPAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{7}
}

func (m *RouteDistinguisherIPAddress) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteDistinguisherIPAddress.Unmarshal(m, b)
}
func (m *RouteDistinguisherIPAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteDistinguisherIPAddress.Marshal(b, m, deterministic)
}
func (m *RouteDistinguisherIPAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteDistinguisherIPAddress.Merge(m, src)
}
func (m *RouteDistinguisherIPAddress) XXX_Size() int {
	return xxx_messageInfo_RouteDistinguisherIPAddress.Size(m)
}
func (m *RouteDistinguisherIPAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteDistinguisherIPAddress.DiscardUnknown(m)
}

var xxx_messageInfo_RouteDistinguisherIPAddress proto.InternalMessageInfo

func (m *RouteDistinguisherIPAddress) GetAdmin() string {
	if m != nil {
		return m.Admin
	}
	return ""
}

func (m *RouteDistinguisherIPAddress) GetAssigned() uint32 {
	if m != nil {
		return m.Assigned
	}
	return 0
}

// Route Distinguisher is expressed by 4 bytes Autonomous System
type RouteDistinguisherFourOctetAS struct {
	Admin                uint32   `protobuf:"varint,1,opt,name=admin,proto3" json:"admin,omitempty"`
	Assigned             uint32   `protobuf:"varint,2,opt,name=assigned,proto3" json:"assigned,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RouteDistinguisherFourOctetAS) Reset()         { *m = RouteDistinguisherFourOctetAS{} }
func (m *RouteDistinguisherFourOctetAS) String() string { return proto.CompactTextString(m) }
func (*RouteDistinguisherFourOctetAS) ProtoMessage()    {}
func (*RouteDistinguisherFourOctetAS) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{8}
}

func (m *RouteDistinguisherFourOctetAS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteDistinguisherFourOctetAS.Unmarshal(m, b)
}
func (m *RouteDistinguisherFourOctetAS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteDistinguisherFourOctetAS.Marshal(b, m, deterministic)
}
func (m *RouteDistinguisherFourOctetAS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteDistinguisherFourOctetAS.Merge(m, src)
}
func (m *RouteDistinguisherFourOctetAS) XXX_Size() int {
	return xxx_messageInfo_RouteDistinguisherFourOctetAS.Size(m)
}
func (m *RouteDistinguisherFourOctetAS) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteDistinguisherFourOctetAS.DiscardUnknown(m)
}

var xxx_messageInfo_RouteDistinguisherFourOctetAS proto.InternalMessageInfo

func (m *RouteDistinguisherFourOctetAS) GetAdmin() uint32 {
	if m != nil {
		return m.Admin
	}
	return 0
}

func (m *RouteDistinguisherFourOctetAS) GetAssigned() uint32 {
	if m != nil {
		return m.Assigned
	}
	return 0
}

// RequestVPNLabel defines the rpc message which is sent by the client to the Gateway
// to request Layer3 VPN label for the corresponding Route Distinguisher.
type RequestVPNLabel struct {
	// Route Distinguisher must be one of
	// RouteDistinguisherTwoOctetAS,
	// RouteDistinguisherIPAddressAS,
	// or RouteDistinguisherFourOctetAS.
	Rd                   *any.Any `protobuf:"bytes,1,opt,name=rd,proto3" json:"rd,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVPNLabel) Reset()         { *m = RequestVPNLabel{} }
func (m *RequestVPNLabel) String() string { return proto.CompactTextString(m) }
func (*RequestVPNLabel) ProtoMessage()    {}
func (*RequestVPNLabel) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{9}
}

func (m *RequestVPNLabel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVPNLabel.Unmarshal(m, b)
}
func (m *RequestVPNLabel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVPNLabel.Marshal(b, m, deterministic)
}
func (m *RequestVPNLabel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVPNLabel.Merge(m, src)
}
func (m *RequestVPNLabel) XXX_Size() int {
	return xxx_messageInfo_RequestVPNLabel.Size(m)
}
func (m *RequestVPNLabel) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVPNLabel.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVPNLabel proto.InternalMessageInfo

func (m *RequestVPNLabel) GetRd() *any.Any {
	if m != nil {
		return m.Rd
	}
	return nil
}

// ResponseVPNLabel defines the rpc message which is sent by the Gateway to the client
// in response to RequestVPNLabel. If VPN label is not found, the error will be returned.
type ResponseVPNLabel struct {
	Label                []uint32      `protobuf:"varint,1,rep,packed,name=label,proto3" json:"label,omitempty"`
	Err                  GatewayErrors `protobuf:"varint,2,opt,name=err,proto3,enum=apis.GatewayErrors" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ResponseVPNLabel) Reset()         { *m = ResponseVPNLabel{} }
func (m *ResponseVPNLabel) String() string { return proto.CompactTextString(m) }
func (*ResponseVPNLabel) ProtoMessage()    {}
func (*ResponseVPNLabel) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{10}
}

func (m *ResponseVPNLabel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseVPNLabel.Unmarshal(m, b)
}
func (m *ResponseVPNLabel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseVPNLabel.Marshal(b, m, deterministic)
}
func (m *ResponseVPNLabel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseVPNLabel.Merge(m, src)
}
func (m *ResponseVPNLabel) XXX_Size() int {
	return xxx_messageInfo_ResponseVPNLabel.Size(m)
}
func (m *ResponseVPNLabel) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseVPNLabel.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseVPNLabel proto.InternalMessageInfo

func (m *ResponseVPNLabel) GetLabel() []uint32 {
	if m != nil {
		return m.Label
	}
	return nil
}

func (m *ResponseVPNLabel) GetErr() GatewayErrors {
	if m != nil {
		return m.Err
	}
	return GatewayErrors_OK
}

func init() {
	proto.RegisterEnum("apis.EndpointType", EndpointType_name, EndpointType_value)
	proto.RegisterEnum("apis.GatewayErrors", GatewayErrors_name, GatewayErrors_value)
	proto.RegisterType((*Endpoint)(nil), "apis.endpoint")
	proto.RegisterType((*Latency)(nil), "apis.latency")
	proto.RegisterType((*QoeParameters)(nil), "apis.qoe_parameters")
	proto.RegisterType((*Qoe)(nil), "apis.qoe")
	proto.RegisterType((*RequestQoE)(nil), "apis.RequestQoE")
	proto.RegisterMapType((map[int32]*Qoe)(nil), "apis.RequestQoE.QoesEntry")
	proto.RegisterType((*ResponseQoE)(nil), "apis.ResponseQoE")
	proto.RegisterMapType((map[int32]*Qoe)(nil), "apis.ResponseQoE.QoesEntry")
	proto.RegisterType((*RouteDistinguisherTwoOctetAS)(nil), "apis.RouteDistinguisherTwoOctetAS")
	proto.RegisterType((*RouteDistinguisherIPAddress)(nil), "apis.RouteDistinguisherIPAddress")
	proto.RegisterType((*RouteDistinguisherFourOctetAS)(nil), "apis.RouteDistinguisherFourOctetAS")
	proto.RegisterType((*RequestVPNLabel)(nil), "apis.RequestVPNLabel")
	proto.RegisterType((*ResponseVPNLabel)(nil), "apis.ResponseVPNLabel")
}

func init() { proto.RegisterFile("gateway.proto", fileDescriptor_f1a937782ebbded5) }

var fileDescriptor_f1a937782ebbded5 = []byte{
	// 646 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0x5b, 0x6f, 0xd3, 0x4c,
	0x10, 0xad, 0xed, 0xa4, 0x69, 0x26, 0x4d, 0xeb, 0x6f, 0xbf, 0x80, 0x4c, 0x5a, 0x44, 0x64, 0x71,
	0x89, 0xfa, 0xe0, 0x4a, 0xa1, 0x82, 0x82, 0xc4, 0x43, 0x51, 0x4c, 0x09, 0x2d, 0x76, 0xba, 0x49,
	0xcb, 0x23, 0xda, 0xc6, 0x4b, 0xb0, 0x48, 0xbd, 0xce, 0xee, 0xa6, 0x95, 0xdf, 0x10, 0x3f, 0x88,
	0xdf, 0x88, 0xbc, 0xbe, 0x94, 0xd0, 0x0a, 0x21, 0x24, 0xde, 0x76, 0x66, 0xce, 0x9c, 0x9d, 0xcb,
	0x19, 0x68, 0x4e, 0x89, 0xa4, 0x57, 0x24, 0x71, 0x62, 0xce, 0x24, 0x43, 0x15, 0x12, 0x87, 0xa2,
	0x7d, 0x6f, 0xca, 0xd8, 0x74, 0x46, 0x77, 0x95, 0xef, 0x7c, 0xf1, 0x69, 0x97, 0x44, 0x39, 0xc0,
	0x3e, 0x86, 0x35, 0x1a, 0x05, 0x31, 0x0b, 0x23, 0x89, 0x1e, 0x43, 0x45, 0x26, 0x31, 0xb5, 0xb4,
	0x8e, 0xd6, 0xdd, 0xe8, 0x21, 0x27, 0xcd, 0x75, 0xdc, 0x3c, 0x3a, 0x4e, 0x62, 0x8a, 0x55, 0x1c,
	0x59, 0x50, 0x23, 0x41, 0xc0, 0xa9, 0x10, 0x96, 0xde, 0xd1, 0xba, 0xeb, 0xb8, 0x30, 0xed, 0x57,
	0x50, 0x9b, 0x11, 0x49, 0xa3, 0x49, 0x82, 0x5a, 0x50, 0xbd, 0x24, 0xb3, 0x45, 0xc6, 0x56, 0xc5,
	0x99, 0x81, 0xb6, 0xa1, 0x7e, 0x49, 0x78, 0x48, 0x64, 0xc8, 0x22, 0x95, 0x5c, 0xc5, 0xd7, 0x0e,
	0xfb, 0x05, 0x6c, 0xcc, 0x19, 0xfd, 0x18, 0x13, 0x4e, 0x2e, 0xa8, 0xa4, 0x5c, 0xa0, 0x27, 0x25,
	0xa1, 0xe2, 0x69, 0xf4, 0x9a, 0x59, 0x55, 0xb9, 0x13, 0x17, 0x51, 0xfb, 0xbb, 0x06, 0xc6, 0x9c,
	0x51, 0xd4, 0x01, 0x43, 0xf0, 0x49, 0x0e, 0xde, 0xc8, 0xc0, 0x45, 0x83, 0x38, 0x0d, 0xa5, 0x88,
	0x40, 0x48, 0xf5, 0xf9, 0x2d, 0x88, 0x40, 0xa4, 0x73, 0x48, 0xa9, 0x2c, 0x43, 0x21, 0x5a, 0x19,
	0x62, 0xb9, 0x2e, 0xac, 0xfe, 0x6a, 0x41, 0x75, 0x46, 0xce, 0xe9, 0xcc, 0xaa, 0x74, 0x8c, 0x6e,
	0x13, 0x67, 0x06, 0x7a, 0x04, 0x06, 0xe5, 0xdc, 0xaa, 0xaa, 0x21, 0xfe, 0x9f, 0x65, 0x1f, 0x66,
	0x4b, 0x71, 0x39, 0x67, 0x69, 0x32, 0xe5, 0xdc, 0xfe, 0xaa, 0x01, 0x60, 0x3a, 0x5f, 0x50, 0x21,
	0x4f, 0x98, 0x8b, 0x1c, 0xa8, 0xcc, 0x19, 0x15, 0x96, 0xd6, 0x31, 0xba, 0x8d, 0x5e, 0x3b, 0x4b,
	0xbb, 0x8e, 0x3b, 0x27, 0x8c, 0x0a, 0x37, 0x92, 0x3c, 0xc1, 0x0a, 0xd7, 0x7e, 0x0d, 0xf5, 0xd2,
	0x85, 0x4c, 0x30, 0xbe, 0xd0, 0x24, 0x9f, 0x74, 0xfa, 0x44, 0x0f, 0x8a, 0xe9, 0x67, 0x6d, 0xd6,
	0xcb, 0x26, 0xf2, 0x45, 0xbc, 0xd4, 0xf7, 0x35, 0xfb, 0x9b, 0x06, 0x0d, 0x4c, 0x45, 0xcc, 0x22,
	0x41, 0xd3, 0x1a, 0x76, 0x97, 0x6a, 0xd8, 0x2a, 0x6a, 0x28, 0x01, 0xff, 0xa4, 0x88, 0x21, 0x6c,
	0x63, 0xb6, 0x90, 0xb4, 0x1f, 0x0a, 0x19, 0x46, 0xd3, 0x45, 0x28, 0x3e, 0x53, 0x3e, 0xbe, 0x62,
	0xfe, 0x44, 0x52, 0x79, 0x30, 0x4a, 0x87, 0x4c, 0x82, 0x8b, 0x30, 0x52, 0xc4, 0x4d, 0x9c, 0x19,
	0xa8, 0x0d, 0x6b, 0x44, 0x88, 0x70, 0x1a, 0xd1, 0x40, 0xb1, 0x37, 0x71, 0x69, 0xdb, 0x3e, 0x6c,
	0xdd, 0x64, 0x1c, 0x0c, 0x0f, 0x32, 0x8d, 0x2e, 0x13, 0xd6, 0xff, 0x84, 0xf0, 0x04, 0xee, 0xdf,
	0x24, 0x7c, 0xc3, 0x16, 0xfc, 0xef, 0x6b, 0x7c, 0x0e, 0x9b, 0xf9, 0x72, 0xcf, 0x86, 0xde, 0xb1,
	0xd2, 0xcd, 0x43, 0xd0, 0x79, 0x90, 0x0b, 0xb7, 0xe5, 0x64, 0x17, 0xeb, 0x14, 0x17, 0xeb, 0x1c,
	0x44, 0x09, 0xd6, 0x79, 0xda, 0x9c, 0x59, 0x6c, 0xa4, 0xcc, 0x2c, 0x75, 0xa8, 0xdd, 0xa2, 0x43,
	0xfd, 0xf7, 0x3a, 0xdc, 0xd9, 0x87, 0xf5, 0x9f, 0x4f, 0x1c, 0x35, 0xa0, 0x76, 0xea, 0x1d, 0x79,
	0xfe, 0x07, 0xcf, 0x5c, 0x41, 0x6b, 0x50, 0x19, 0x0c, 0xcf, 0xf6, 0x4c, 0x2d, 0x7f, 0x3d, 0x33,
	0x75, 0x54, 0x03, 0x63, 0x34, 0xe8, 0x9b, 0xc6, 0xce, 0x3b, 0x68, 0x2e, 0xf1, 0xa1, 0x55, 0xd0,
	0xfd, 0x23, 0x73, 0x25, 0x45, 0xb8, 0x03, 0xdf, 0xd4, 0x10, 0xc0, 0xaa, 0xeb, 0xf9, 0xae, 0x37,
	0x36, 0x75, 0xd4, 0x84, 0xba, 0x3b, 0x1e, 0xbc, 0x77, 0xfb, 0xfe, 0xe9, 0xd8, 0x34, 0x94, 0xf9,
	0xd6, 0x1f, 0x8d, 0xfb, 0xe9, 0x47, 0x95, 0x9e, 0x80, 0xcd, 0x43, 0x22, 0xc9, 0x15, 0x49, 0x46,
	0x94, 0x5f, 0x86, 0x13, 0x2a, 0xd0, 0x0e, 0x18, 0xa9, 0x28, 0xcd, 0x5f, 0x4f, 0xa1, 0xfd, 0xdf,
	0x0d, 0x61, 0xa2, 0x3d, 0x30, 0xce, 0x86, 0x1e, 0xba, 0xb3, 0x84, 0x2d, 0xe6, 0xd3, 0xbe, 0xbb,
	0x9c, 0x50, 0xf8, 0xcf, 0x57, 0xd5, 0x74, 0x9f, 0xfe, 0x08, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x31,
	0xdb, 0x24, 0x34, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GatawayServicesClient is the client API for GatawayServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GatawayServicesClient interface {
	// API to request specified Quality of Experience
	QoE(ctx context.Context, in *RequestQoE, opts ...grpc.CallOption) (*ResponseQoE, error)
	// API to request L3 VPN label for specified VRF, identified by Route Distinguisher
	VPN(ctx context.Context, in *RequestVPNLabel, opts ...grpc.CallOption) (*ResponseVPNLabel, error)
}

type gatawayServicesClient struct {
	cc *grpc.ClientConn
}

func NewGatawayServicesClient(cc *grpc.ClientConn) GatawayServicesClient {
	return &gatawayServicesClient{cc}
}

func (c *gatawayServicesClient) QoE(ctx context.Context, in *RequestQoE, opts ...grpc.CallOption) (*ResponseQoE, error) {
	out := new(ResponseQoE)
	err := c.cc.Invoke(ctx, "/apis.GatawayServices/QoE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatawayServicesClient) VPN(ctx context.Context, in *RequestVPNLabel, opts ...grpc.CallOption) (*ResponseVPNLabel, error) {
	out := new(ResponseVPNLabel)
	err := c.cc.Invoke(ctx, "/apis.GatawayServices/VPN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatawayServicesServer is the server API for GatawayServices service.
type GatawayServicesServer interface {
	// API to request specified Quality of Experience
	QoE(context.Context, *RequestQoE) (*ResponseQoE, error)
	// API to request L3 VPN label for specified VRF, identified by Route Distinguisher
	VPN(context.Context, *RequestVPNLabel) (*ResponseVPNLabel, error)
}

// UnimplementedGatawayServicesServer can be embedded to have forward compatible implementations.
type UnimplementedGatawayServicesServer struct {
}

func (*UnimplementedGatawayServicesServer) QoE(ctx context.Context, req *RequestQoE) (*ResponseQoE, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QoE not implemented")
}
func (*UnimplementedGatawayServicesServer) VPN(ctx context.Context, req *RequestVPNLabel) (*ResponseVPNLabel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VPN not implemented")
}

func RegisterGatawayServicesServer(s *grpc.Server, srv GatawayServicesServer) {
	s.RegisterService(&_GatawayServices_serviceDesc, srv)
}

func _GatawayServices_QoE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestQoE)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatawayServicesServer).QoE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apis.GatawayServices/QoE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatawayServicesServer).QoE(ctx, req.(*RequestQoE))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatawayServices_VPN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVPNLabel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatawayServicesServer).VPN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apis.GatawayServices/VPN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatawayServicesServer).VPN(ctx, req.(*RequestVPNLabel))
	}
	return interceptor(ctx, in, info, handler)
}

var _GatawayServices_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apis.GatawayServices",
	HandlerType: (*GatawayServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QoE",
			Handler:    _GatawayServices_QoE_Handler,
		},
		{
			MethodName: "VPN",
			Handler:    _GatawayServices_VPN_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway.proto",
}
