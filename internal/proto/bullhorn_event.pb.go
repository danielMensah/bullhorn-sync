// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: bullhorn_event.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityId          int32                  `protobuf:"varint,1,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
	EntityName        string                 `protobuf:"bytes,2,opt,name=entity_name,json=entityName,proto3" json:"entity_name,omitempty"`
	EntityEventType   string                 `protobuf:"bytes,3,opt,name=entity_event_type,json=entityEventType,proto3" json:"entity_event_type,omitempty"`
	UpdatedProperties []string               `protobuf:"bytes,4,rep,name=updated_properties,json=updatedProperties,proto3" json:"updated_properties,omitempty"`
	EventTimestamp    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=event_timestamp,json=eventTimestamp,proto3" json:"event_timestamp,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bullhorn_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_bullhorn_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_bullhorn_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetEntityId() int32 {
	if x != nil {
		return x.EntityId
	}
	return 0
}

func (x *Event) GetEntityName() string {
	if x != nil {
		return x.EntityName
	}
	return ""
}

func (x *Event) GetEntityEventType() string {
	if x != nil {
		return x.EntityEventType
	}
	return ""
}

func (x *Event) GetUpdatedProperties() []string {
	if x != nil {
		return x.UpdatedProperties
	}
	return nil
}

func (x *Event) GetEventTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTimestamp
	}
	return nil
}

var File_bullhorn_event_proto protoreflect.FileDescriptor

var file_bullhorn_event_proto_rawDesc = []byte{
	0x0a, 0x14, 0x62, 0x75, 0x6c, 0x6c, 0x68, 0x6f, 0x72, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x75, 0x6c, 0x6c, 0x68, 0x6f, 0x72, 0x6e,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe5, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x2a, 0x0a, 0x11, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2d, 0x0a, 0x12, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x11, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x12, 0x43, 0x0a, 0x0f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42,
	0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61,
	0x6e, 0x69, 0x65, 0x6c, 0x4d, 0x65, 0x6e, 0x73, 0x61, 0x68, 0x2f, 0x62, 0x75, 0x6c, 0x6c, 0x68,
	0x6f, 0x72, 0x6e, 0x2d, 0x73, 0x79, 0x6e, 0x63, 0x2d, 0x70, 0x6f, 0x63, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_bullhorn_event_proto_rawDescOnce sync.Once
	file_bullhorn_event_proto_rawDescData = file_bullhorn_event_proto_rawDesc
)

func file_bullhorn_event_proto_rawDescGZIP() []byte {
	file_bullhorn_event_proto_rawDescOnce.Do(func() {
		file_bullhorn_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_bullhorn_event_proto_rawDescData)
	})
	return file_bullhorn_event_proto_rawDescData
}

var file_bullhorn_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_bullhorn_event_proto_goTypes = []interface{}{
	(*Event)(nil),                 // 0: bullhorn.event.Event
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_bullhorn_event_proto_depIdxs = []int32{
	1, // 0: bullhorn.event.Event.event_timestamp:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_bullhorn_event_proto_init() }
func file_bullhorn_event_proto_init() {
	if File_bullhorn_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bullhorn_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
			RawDescriptor: file_bullhorn_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bullhorn_event_proto_goTypes,
		DependencyIndexes: file_bullhorn_event_proto_depIdxs,
		MessageInfos:      file_bullhorn_event_proto_msgTypes,
	}.Build()
	File_bullhorn_event_proto = out.File
	file_bullhorn_event_proto_rawDesc = nil
	file_bullhorn_event_proto_goTypes = nil
	file_bullhorn_event_proto_depIdxs = nil
}