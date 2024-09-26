// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: custom_option/custom_option.proto

package custom_option

import (
	api_errors "github.com/averak/gamebox/protobuf/api/api_errors"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MethodErrorDefinition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     api_errors.ErrorCode_Method `protobuf:"varint,1,opt,name=code,proto3,enum=api.api_errors.ErrorCode_Method" json:"code,omitempty"`
	Severity api_errors.ErrorSeverity    `protobuf:"varint,2,opt,name=severity,proto3,enum=api.api_errors.ErrorSeverity" json:"severity,omitempty"`
	Message  string                      `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MethodErrorDefinition) Reset() {
	*x = MethodErrorDefinition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_custom_option_custom_option_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MethodErrorDefinition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MethodErrorDefinition) ProtoMessage() {}

func (x *MethodErrorDefinition) ProtoReflect() protoreflect.Message {
	mi := &file_custom_option_custom_option_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MethodErrorDefinition.ProtoReflect.Descriptor instead.
func (*MethodErrorDefinition) Descriptor() ([]byte, []int) {
	return file_custom_option_custom_option_proto_rawDescGZIP(), []int{0}
}

func (x *MethodErrorDefinition) GetCode() api_errors.ErrorCode_Method {
	if x != nil {
		return x.Code
	}
	return api_errors.ErrorCode_Method(0)
}

func (x *MethodErrorDefinition) GetSeverity() api_errors.ErrorSeverity {
	if x != nil {
		return x.Severity
	}
	return api_errors.ErrorSeverity(0)
}

func (x *MethodErrorDefinition) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type MethodOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MethodErrorDefinitions []*MethodErrorDefinition `protobuf:"bytes,1,rep,name=method_error_definitions,json=methodErrorDefinitions,proto3" json:"method_error_definitions,omitempty"`
	SkipAuthenticate       bool                     `protobuf:"varint,2,opt,name=skip_authenticate,json=skipAuthenticate,proto3" json:"skip_authenticate,omitempty"`
}

func (x *MethodOption) Reset() {
	*x = MethodOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_custom_option_custom_option_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MethodOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MethodOption) ProtoMessage() {}

func (x *MethodOption) ProtoReflect() protoreflect.Message {
	mi := &file_custom_option_custom_option_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MethodOption.ProtoReflect.Descriptor instead.
func (*MethodOption) Descriptor() ([]byte, []int) {
	return file_custom_option_custom_option_proto_rawDescGZIP(), []int{1}
}

func (x *MethodOption) GetMethodErrorDefinitions() []*MethodErrorDefinition {
	if x != nil {
		return x.MethodErrorDefinitions
	}
	return nil
}

func (x *MethodOption) GetSkipAuthenticate() bool {
	if x != nil {
		return x.SkipAuthenticate
	}
	return false
}

var file_custom_option_custom_option_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*MethodOption)(nil),
		Field:         50000,
		Name:          "custom_option.method_option",
		Tag:           "bytes,50000,opt,name=method_option",
		Filename:      "custom_option/custom_option.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional custom_option.MethodOption method_option = 50000;
	E_MethodOption = &file_custom_option_custom_option_proto_extTypes[0]
)

var File_custom_option_custom_option_proto protoreflect.FileDescriptor

var file_custom_option_custom_option_proto_rawDesc = []byte{
	0x0a, 0x21, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x01, 0x0a, 0x15, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x34, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x39, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70,
	0x69, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x65,
	0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x9b, 0x01, 0x0a, 0x0c, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x5e, 0x0a, 0x18, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x64, 0x65, 0x66, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x16, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x73,
	0x6b, 0x69, 0x70, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x73, 0x6b, 0x69, 0x70, 0x41, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x3a, 0x62, 0x0a, 0x0d, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0xa8, 0x01, 0x0a,
	0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x11, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x76, 0x65, 0x72, 0x61, 0x6b, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x62,
	0x6f, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa,
	0x02, 0x0c, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0xca, 0x02,
	0x0c, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0xe2, 0x02, 0x18,
	0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0c, 0x43, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_custom_option_custom_option_proto_rawDescOnce sync.Once
	file_custom_option_custom_option_proto_rawDescData = file_custom_option_custom_option_proto_rawDesc
)

func file_custom_option_custom_option_proto_rawDescGZIP() []byte {
	file_custom_option_custom_option_proto_rawDescOnce.Do(func() {
		file_custom_option_custom_option_proto_rawDescData = protoimpl.X.CompressGZIP(file_custom_option_custom_option_proto_rawDescData)
	})
	return file_custom_option_custom_option_proto_rawDescData
}

var file_custom_option_custom_option_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_custom_option_custom_option_proto_goTypes = []any{
	(*MethodErrorDefinition)(nil),      // 0: custom_option.MethodErrorDefinition
	(*MethodOption)(nil),               // 1: custom_option.MethodOption
	(api_errors.ErrorCode_Method)(0),   // 2: api.api_errors.ErrorCode.Method
	(api_errors.ErrorSeverity)(0),      // 3: api.api_errors.ErrorSeverity
	(*descriptorpb.MethodOptions)(nil), // 4: google.protobuf.MethodOptions
}
var file_custom_option_custom_option_proto_depIdxs = []int32{
	2, // 0: custom_option.MethodErrorDefinition.code:type_name -> api.api_errors.ErrorCode.Method
	3, // 1: custom_option.MethodErrorDefinition.severity:type_name -> api.api_errors.ErrorSeverity
	0, // 2: custom_option.MethodOption.method_error_definitions:type_name -> custom_option.MethodErrorDefinition
	4, // 3: custom_option.method_option:extendee -> google.protobuf.MethodOptions
	1, // 4: custom_option.method_option:type_name -> custom_option.MethodOption
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	4, // [4:5] is the sub-list for extension type_name
	3, // [3:4] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_custom_option_custom_option_proto_init() }
func file_custom_option_custom_option_proto_init() {
	if File_custom_option_custom_option_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_custom_option_custom_option_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*MethodErrorDefinition); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_custom_option_custom_option_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*MethodOption); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_custom_option_custom_option_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_custom_option_custom_option_proto_goTypes,
		DependencyIndexes: file_custom_option_custom_option_proto_depIdxs,
		MessageInfos:      file_custom_option_custom_option_proto_msgTypes,
		ExtensionInfos:    file_custom_option_custom_option_proto_extTypes,
	}.Build()
	File_custom_option_custom_option_proto = out.File
	file_custom_option_custom_option_proto_rawDesc = nil
	file_custom_option_custom_option_proto_goTypes = nil
	file_custom_option_custom_option_proto_depIdxs = nil
}
