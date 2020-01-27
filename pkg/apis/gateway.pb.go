// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gateway.proto

package apis

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
	// Types that are valid to be assigned to RdType:
	//	*RequestVPNLabel_Rd2As
	//	*RequestVPNLabel_RdIPAddr
	//	*RequestVPNLabel_Rd4As
	RdType               isRequestVPNLabel_RdType `protobuf_oneof:"rd_type"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
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

type isRequestVPNLabel_RdType interface {
	isRequestVPNLabel_RdType()
}

type RequestVPNLabel_Rd2As struct {
	Rd2As *RouteDistinguisherTwoOctetAS `protobuf:"bytes,1,opt,name=rd2as,proto3,oneof"`
}

type RequestVPNLabel_RdIPAddr struct {
	RdIPAddr *RouteDistinguisherIPAddress `protobuf:"bytes,2,opt,name=rdIPAddr,proto3,oneof"`
}

type RequestVPNLabel_Rd4As struct {
	Rd4As *RouteDistinguisherFourOctetAS `protobuf:"bytes,3,opt,name=rd4as,proto3,oneof"`
}

func (*RequestVPNLabel_Rd2As) isRequestVPNLabel_RdType() {}

func (*RequestVPNLabel_RdIPAddr) isRequestVPNLabel_RdType() {}

func (*RequestVPNLabel_Rd4As) isRequestVPNLabel_RdType() {}

func (m *RequestVPNLabel) GetRdType() isRequestVPNLabel_RdType {
	if m != nil {
		return m.RdType
	}
	return nil
}

func (m *RequestVPNLabel) GetRd2As() *RouteDistinguisherTwoOctetAS {
	if x, ok := m.GetRdType().(*RequestVPNLabel_Rd2As); ok {
		return x.Rd2As
	}
	return nil
}

func (m *RequestVPNLabel) GetRdIPAddr() *RouteDistinguisherIPAddress {
	if x, ok := m.GetRdType().(*RequestVPNLabel_RdIPAddr); ok {
		return x.RdIPAddr
	}
	return nil
}

func (m *RequestVPNLabel) GetRd4As() *RouteDistinguisherFourOctetAS {
	if x, ok := m.GetRdType().(*RequestVPNLabel_Rd4As); ok {
		return x.Rd4As
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*RequestVPNLabel) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*RequestVPNLabel_Rd2As)(nil),
		(*RequestVPNLabel_RdIPAddr)(nil),
		(*RequestVPNLabel_Rd4As)(nil),
	}
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
	// 677 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x5d, 0x6f, 0xd3, 0x4a,
	0x10, 0x8d, 0xe3, 0xa4, 0x49, 0x26, 0x4d, 0xae, 0xef, 0xde, 0xde, 0xab, 0x28, 0xed, 0x15, 0xc1,
	0x08, 0x88, 0xfa, 0x10, 0xa4, 0x10, 0xa1, 0x52, 0x84, 0x50, 0xab, 0x18, 0x1a, 0x5a, 0xec, 0x74,
	0x93, 0x96, 0xc7, 0x6a, 0x1b, 0xaf, 0x8a, 0x45, 0xea, 0x75, 0x76, 0x37, 0xad, 0xf2, 0x86, 0xf8,
	0x41, 0xfc, 0x1c, 0x7e, 0x0f, 0x5a, 0xaf, 0xed, 0x36, 0xf4, 0x43, 0x08, 0x89, 0x37, 0xcf, 0xee,
	0x39, 0x67, 0x66, 0x8e, 0x67, 0x16, 0x6a, 0x67, 0x44, 0xd2, 0x4b, 0xb2, 0xe8, 0x44, 0x9c, 0x49,
	0x86, 0x0a, 0x24, 0x0a, 0x84, 0x7d, 0x00, 0x65, 0x1a, 0xfa, 0x11, 0x0b, 0x42, 0x89, 0x9e, 0x40,
	0x41, 0x2e, 0x22, 0xda, 0x30, 0x5a, 0x46, 0xbb, 0xde, 0x45, 0x1d, 0x05, 0xe8, 0x38, 0xc9, 0xed,
	0x78, 0x11, 0x51, 0x1c, 0xdf, 0xa3, 0x06, 0x94, 0x88, 0xef, 0x73, 0x2a, 0x44, 0x23, 0xdf, 0x32,
	0xda, 0xab, 0x38, 0x0d, 0xed, 0xd7, 0x50, 0x9a, 0x12, 0x49, 0xc3, 0xc9, 0x02, 0xad, 0x41, 0xf1,
	0x82, 0x4c, 0xe7, 0x5a, 0xad, 0x88, 0x75, 0x80, 0x36, 0xa0, 0x72, 0x41, 0x78, 0x40, 0x64, 0xc0,
	0xc2, 0x98, 0x5c, 0xc4, 0x57, 0x07, 0xf6, 0x4b, 0xa8, 0xcf, 0x18, 0x3d, 0x89, 0x08, 0x27, 0xe7,
	0x54, 0x52, 0x2e, 0xd0, 0xd3, 0x4c, 0x30, 0xd6, 0xa9, 0x76, 0x6b, 0xba, 0xaa, 0xe4, 0x10, 0xa7,
	0xb7, 0xf6, 0x37, 0x03, 0xcc, 0x19, 0xa3, 0xa8, 0x05, 0xa6, 0xe0, 0x93, 0x04, 0x5c, 0xd7, 0xe0,
	0xb4, 0x41, 0xac, 0xae, 0x14, 0xc2, 0x17, 0x32, 0x4e, 0x7e, 0x0b, 0xc2, 0x17, 0xca, 0x07, 0x25,
	0xd5, 0x30, 0x63, 0xc4, 0x9a, 0x46, 0x2c, 0xd7, 0x85, 0xe3, 0x5c, 0x6b, 0x50, 0x9c, 0x92, 0x53,
	0x3a, 0x6d, 0x14, 0x5a, 0x66, 0xbb, 0x86, 0x75, 0x80, 0x1e, 0x83, 0x49, 0x39, 0x6f, 0x14, 0x63,
	0x13, 0xff, 0xd1, 0xec, 0x77, 0xda, 0x79, 0x87, 0x73, 0xa6, 0xc8, 0x94, 0x73, 0xfb, 0x8b, 0x01,
	0x80, 0xe9, 0x6c, 0x4e, 0x85, 0x3c, 0x64, 0x0e, 0xea, 0x40, 0x61, 0xc6, 0xa8, 0x68, 0x18, 0x2d,
	0xb3, 0x5d, 0xed, 0x36, 0x35, 0xed, 0xea, 0xbe, 0x73, 0xc8, 0xa8, 0x70, 0x42, 0xc9, 0x17, 0x38,
	0xc6, 0x35, 0x77, 0xa1, 0x92, 0x1d, 0x21, 0x0b, 0xcc, 0xcf, 0x74, 0x91, 0x38, 0xad, 0x3e, 0xd1,
	0x83, 0xd4, 0x7d, 0xdd, 0x66, 0x25, 0x6b, 0x22, 0xf9, 0x11, 0xdb, 0xf9, 0x2d, 0xc3, 0xfe, 0x6a,
	0x40, 0x15, 0x53, 0x11, 0xb1, 0x50, 0x50, 0x55, 0xc3, 0xb3, 0xa5, 0x1a, 0xd6, 0xd3, 0x1a, 0x32,
	0xc0, 0x1f, 0x29, 0x62, 0x08, 0x1b, 0x98, 0xcd, 0x25, 0xed, 0x07, 0x42, 0x06, 0xe1, 0xd9, 0x3c,
	0x10, 0x9f, 0x28, 0x1f, 0x5f, 0x32, 0x6f, 0x22, 0xa9, 0xdc, 0x19, 0x29, 0x93, 0x89, 0x7f, 0x1e,
	0x84, 0xb1, 0x70, 0x0d, 0xeb, 0x00, 0x35, 0xa1, 0x4c, 0x84, 0x08, 0xce, 0x42, 0xea, 0xc7, 0xea,
	0x35, 0x9c, 0xc5, 0xb6, 0x07, 0xeb, 0x37, 0x15, 0x07, 0xc3, 0x1d, 0x3d, 0xa3, 0xcb, 0x82, 0x95,
	0x5f, 0x11, 0x3c, 0x84, 0xff, 0x6f, 0x0a, 0xbe, 0x65, 0x73, 0xfe, 0xfb, 0x35, 0x7e, 0x37, 0xe0,
	0xaf, 0xe4, 0xef, 0x1e, 0x0f, 0xdd, 0x83, 0x78, 0x70, 0xb6, 0xa1, 0xc8, 0xfd, 0x2e, 0x11, 0xc9,
	0xf0, 0xda, 0x89, 0xff, 0xf7, 0x98, 0xb3, 0x97, 0xc3, 0x9a, 0x82, 0xde, 0x40, 0x99, 0xfb, 0xba,
	0xc7, 0xc4, 0xed, 0x87, 0x77, 0xd1, 0x33, 0x27, 0xf6, 0x72, 0x38, 0x23, 0xa1, 0x57, 0x2a, 0x79,
	0x8f, 0x88, 0x64, 0xea, 0x1f, 0xdd, 0xc5, 0xbe, 0xd6, 0xb6, 0xce, 0xde, 0x23, 0x62, 0xb7, 0x02,
	0x25, 0xee, 0x9f, 0xa8, 0xb7, 0xc1, 0xf6, 0xc0, 0x4a, 0x27, 0x26, 0x6b, 0x2c, 0xdb, 0x13, 0xe3,
	0x96, 0x3d, 0xc9, 0xdf, 0xbf, 0x27, 0x9b, 0x5b, 0xb0, 0x7a, 0xfd, 0x09, 0x42, 0x55, 0x28, 0x1d,
	0xb9, 0xfb, 0xae, 0xf7, 0xd1, 0xb5, 0x72, 0xa8, 0x0c, 0x85, 0xc1, 0xf0, 0xb8, 0x67, 0x19, 0xc9,
	0xd7, 0x0b, 0x2b, 0x8f, 0x4a, 0x60, 0x8e, 0x06, 0x7d, 0xcb, 0xdc, 0x7c, 0x0f, 0xb5, 0x25, 0x3d,
	0xb4, 0x02, 0x79, 0x6f, 0xdf, 0xca, 0x29, 0x84, 0x33, 0xf0, 0x2c, 0x03, 0x01, 0xac, 0x38, 0xae,
	0xe7, 0xb8, 0x63, 0x2b, 0x8f, 0x6a, 0x50, 0x71, 0xc6, 0x83, 0x0f, 0x4e, 0xdf, 0x3b, 0x1a, 0x5b,
	0x66, 0x1c, 0xee, 0x79, 0xa3, 0x71, 0x5f, 0x25, 0x2a, 0x74, 0x39, 0xd4, 0x13, 0xad, 0x11, 0xe5,
	0x17, 0xc1, 0x84, 0xa2, 0x4d, 0x30, 0xd5, 0xce, 0x58, 0x3f, 0x6f, 0x6a, 0xf3, 0xef, 0x1b, 0x7b,
	0x83, 0x7a, 0x60, 0x1e, 0x0f, 0x5d, 0xf4, 0xef, 0x12, 0x36, 0xb5, 0xa7, 0xf9, 0xdf, 0x32, 0x21,
	0x3d, 0x3f, 0x5d, 0x89, 0xdf, 0xe9, 0xe7, 0x3f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xf2, 0xf4, 0x15,
	0x5d, 0xb8, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GatewayServiceClient is the client API for GatewayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GatewayServiceClient interface {
	// API to request specified Quality of Experience
	QoE(ctx context.Context, in *RequestQoE, opts ...grpc.CallOption) (*ResponseQoE, error)
	// API to request L3 VPN label for specified VRF, identified by Route Distinguisher
	VPN(ctx context.Context, in *RequestVPNLabel, opts ...grpc.CallOption) (*ResponseVPNLabel, error)
}

type gatewayServiceClient struct {
	cc *grpc.ClientConn
}

func NewGatewayServiceClient(cc *grpc.ClientConn) GatewayServiceClient {
	return &gatewayServiceClient{cc}
}

func (c *gatewayServiceClient) QoE(ctx context.Context, in *RequestQoE, opts ...grpc.CallOption) (*ResponseQoE, error) {
	out := new(ResponseQoE)
	err := c.cc.Invoke(ctx, "/apis.GatewayService/QoE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayServiceClient) VPN(ctx context.Context, in *RequestVPNLabel, opts ...grpc.CallOption) (*ResponseVPNLabel, error) {
	out := new(ResponseVPNLabel)
	err := c.cc.Invoke(ctx, "/apis.GatewayService/VPN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayServiceServer is the server API for GatewayService service.
type GatewayServiceServer interface {
	// API to request specified Quality of Experience
	QoE(context.Context, *RequestQoE) (*ResponseQoE, error)
	// API to request L3 VPN label for specified VRF, identified by Route Distinguisher
	VPN(context.Context, *RequestVPNLabel) (*ResponseVPNLabel, error)
}

// UnimplementedGatewayServiceServer can be embedded to have forward compatible implementations.
type UnimplementedGatewayServiceServer struct {
}

func (*UnimplementedGatewayServiceServer) QoE(ctx context.Context, req *RequestQoE) (*ResponseQoE, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QoE not implemented")
}
func (*UnimplementedGatewayServiceServer) VPN(ctx context.Context, req *RequestVPNLabel) (*ResponseVPNLabel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VPN not implemented")
}

func RegisterGatewayServiceServer(s *grpc.Server, srv GatewayServiceServer) {
	s.RegisterService(&_GatewayService_serviceDesc, srv)
}

func _GatewayService_QoE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestQoE)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServiceServer).QoE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apis.GatewayService/QoE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServiceServer).QoE(ctx, req.(*RequestQoE))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayService_VPN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVPNLabel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServiceServer).VPN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apis.GatewayService/VPN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServiceServer).VPN(ctx, req.(*RequestVPNLabel))
	}
	return interceptor(ctx, in, info, handler)
}

var _GatewayService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apis.GatewayService",
	HandlerType: (*GatewayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QoE",
			Handler:    _GatewayService_QoE_Handler,
		},
		{
			MethodName: "VPN",
			Handler:    _GatewayService_VPN_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway.proto",
}
