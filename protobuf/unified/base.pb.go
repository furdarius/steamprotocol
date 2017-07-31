// Code generated by protoc-gen-go.
// source: steammessages_unified_base.steamclient.proto
// DO NOT EDIT!

package unified

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EProtoExecutionSite int32

const (
	EProtoExecutionSite_k_EProtoExecutionSiteUnknown     EProtoExecutionSite = 0
	EProtoExecutionSite_k_EProtoExecutionSiteSteamClient EProtoExecutionSite = 2
)

var EProtoExecutionSite_name = map[int32]string{
	0: "k_EProtoExecutionSiteUnknown",
	2: "k_EProtoExecutionSiteSteamClient",
}
var EProtoExecutionSite_value = map[string]int32{
	"k_EProtoExecutionSiteUnknown":     0,
	"k_EProtoExecutionSiteSteamClient": 2,
}

func (x EProtoExecutionSite) Enum() *EProtoExecutionSite {
	p := new(EProtoExecutionSite)
	*p = x
	return p
}
func (x EProtoExecutionSite) String() string {
	return proto.EnumName(EProtoExecutionSite_name, int32(x))
}
func (x *EProtoExecutionSite) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EProtoExecutionSite_value, data, "EProtoExecutionSite")
	if err != nil {
		return err
	}
	*x = EProtoExecutionSite(value)
	return nil
}
func (EProtoExecutionSite) EnumDescriptor() ([]byte, []int) { return base_fileDescriptor0, []int{0} }

type NoResponse struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *NoResponse) Reset()                    { *m = NoResponse{} }
func (m *NoResponse) String() string            { return proto.CompactTextString(m) }
func (*NoResponse) ProtoMessage()               {}
func (*NoResponse) Descriptor() ([]byte, []int) { return base_fileDescriptor0, []int{0} }

var E_Description = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50000,
	Name:          "description",
	Tag:           "bytes,50000,opt,name=description",
}

var E_ServiceDescription = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.ServiceOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50000,
	Name:          "service_description",
	Tag:           "bytes,50000,opt,name=service_description",
}

var E_ServiceExecutionSite = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.ServiceOptions)(nil),
	ExtensionType: (*EProtoExecutionSite)(nil),
	Field:         50008,
	Name:          "service_execution_site",
	Tag:           "varint,50008,opt,name=service_execution_site,enum=EProtoExecutionSite,def=0",
}

var E_MethodDescription = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50000,
	Name:          "method_description",
	Tag:           "bytes,50000,opt,name=method_description",
}

var E_EnumDescription = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.EnumOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50000,
	Name:          "enum_description",
	Tag:           "bytes,50000,opt,name=enum_description",
}

var E_EnumValueDescription = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.EnumValueOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50000,
	Name:          "enum_value_description",
	Tag:           "bytes,50000,opt,name=enum_value_description",
}

func init() {
	proto.RegisterType((*NoResponse)(nil), "NoResponse")
	proto.RegisterEnum("EProtoExecutionSite", EProtoExecutionSite_name, EProtoExecutionSite_value)
	proto.RegisterExtension(E_Description)
	proto.RegisterExtension(E_ServiceDescription)
	proto.RegisterExtension(E_ServiceExecutionSite)
	proto.RegisterExtension(E_MethodDescription)
	proto.RegisterExtension(E_EnumDescription)
	proto.RegisterExtension(E_EnumValueDescription)
}

var base_fileDescriptor0 = []byte{
	// 306 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x90, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x86, 0x1b, 0xc5, 0x83, 0xa3, 0x48, 0x48, 0xa5, 0x88, 0x54, 0x8d, 0xe2, 0x41, 0x44, 0xb6,
	0x20, 0x1e, 0x24, 0x88, 0x07, 0x4b, 0xc4, 0x8b, 0x1f, 0x18, 0xf4, 0x26, 0x21, 0x4d, 0xa6, 0x71,
	0x69, 0xb2, 0x1b, 0xb2, 0xbb, 0xd5, 0xa3, 0x27, 0x7f, 0x9f, 0x47, 0x7f, 0x8e, 0xcd, 0x86, 0x80,
	0xf9, 0x40, 0x8f, 0xc9, 0xfb, 0x3e, 0xb3, 0xcf, 0x0c, 0x9c, 0x08, 0x89, 0x41, 0x9a, 0xa2, 0x10,
	0x41, 0x8c, 0xc2, 0x57, 0x8c, 0x4e, 0x29, 0x46, 0xfe, 0x24, 0x10, 0x48, 0x74, 0x14, 0x26, 0x14,
	0x99, 0x24, 0x59, 0xce, 0x25, 0xdf, 0xb6, 0x63, 0xce, 0xe3, 0x04, 0x47, 0xfa, 0x6b, 0xa2, 0xa6,
	0xa3, 0x08, 0x45, 0x98, 0xd3, 0x4c, 0xf2, 0xbc, 0x6c, 0x1c, 0xac, 0x03, 0xdc, 0xf1, 0x47, 0x14,
	0x19, 0x67, 0x02, 0x8f, 0x5f, 0xa0, 0xef, 0x3e, 0x14, 0xff, 0xdd, 0x77, 0x0c, 0x95, 0xa4, 0x9c,
	0x79, 0x54, 0xa2, 0x65, 0xc3, 0x70, 0xe6, 0x77, 0x04, 0x4f, 0x6c, 0xc6, 0xf8, 0x1b, 0x33, 0x7b,
	0xd6, 0x21, 0xd8, 0x9d, 0x0d, 0xaf, 0x50, 0x1a, 0x6b, 0x25, 0x73, 0xc9, 0x39, 0x83, 0xb5, 0x4a,
	0x60, 0x91, 0x5b, 0x3b, 0xa4, 0xd4, 0x23, 0x95, 0x1e, 0xb9, 0xa6, 0x98, 0x44, 0xf7, 0x3a, 0x15,
	0x5b, 0x5f, 0x9f, 0xcb, 0xb6, 0x71, 0xb4, 0xea, 0x5c, 0x42, 0x5f, 0x60, 0x3e, 0xa7, 0x21, 0xfa,
	0xbf, 0xe9, 0xbd, 0x16, 0xed, 0x95, 0xad, 0x26, 0xaf, 0x60, 0x50, 0xf1, 0x58, 0xb9, 0xf9, 0xa2,
	0xd8, 0xeb, 0xdf, 0x11, 0xdf, 0x7a, 0xc4, 0xc6, 0xe9, 0x26, 0xe9, 0xd8, 0xcd, 0xf9, 0xf3, 0x28,
	0xce, 0x05, 0x58, 0x29, 0xca, 0x57, 0x1e, 0xd5, 0xac, 0x77, 0x5b, 0x4f, 0xde, 0xea, 0x52, 0x53,
	0xfa, 0x1c, 0x4c, 0x64, 0x2a, 0xad, 0xb1, 0xc3, 0x16, 0xeb, 0x2e, 0x2a, 0x4d, 0x72, 0x0c, 0x03,
	0x4d, 0xce, 0x83, 0x44, 0xd5, 0x2f, 0xb6, 0xdf, 0xc9, 0x3f, 0x17, 0xbd, 0xc6, 0x90, 0xab, 0x95,
	0x1b, 0xe3, 0xc3, 0xe8, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5c, 0xf6, 0x07, 0xbb, 0x6e, 0x02,
	0x00, 0x00,
}
