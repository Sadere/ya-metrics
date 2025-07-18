// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: proto/metrics/v1/metrics.proto

package v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metric_MetricType int32

const (
	Metric_GAUGE   Metric_MetricType = 0
	Metric_COUNTER Metric_MetricType = 1
)

// Enum value maps for Metric_MetricType.
var (
	Metric_MetricType_name = map[int32]string{
		0: "GAUGE",
		1: "COUNTER",
	}
	Metric_MetricType_value = map[string]int32{
		"GAUGE":   0,
		"COUNTER": 1,
	}
)

func (x Metric_MetricType) Enum() *Metric_MetricType {
	p := new(Metric_MetricType)
	*p = x
	return p
}

func (x Metric_MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Metric_MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_metrics_v1_metrics_proto_enumTypes[0].Descriptor()
}

func (Metric_MetricType) Type() protoreflect.EnumType {
	return &file_proto_metrics_v1_metrics_proto_enumTypes[0]
}

func (x Metric_MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Metric_MetricType.Descriptor instead.
func (Metric_MetricType) EnumDescriptor() ([]byte, []int) {
	return file_proto_metrics_v1_metrics_proto_rawDescGZIP(), []int{0, 0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    string            `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	MType Metric_MetricType `protobuf:"varint,2,opt,name=MType,proto3,enum=metrics.v1.Metric_MetricType" json:"MType,omitempty"`
	// Types that are assignable to MetricValue:
	//
	//	*Metric_Value
	//	*Metric_Delta
	MetricValue isMetric_MetricValue `protobuf_oneof:"MetricValue"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_v1_metrics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_v1_metrics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metrics_v1_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Metric) GetMType() Metric_MetricType {
	if x != nil {
		return x.MType
	}
	return Metric_GAUGE
}

func (m *Metric) GetMetricValue() isMetric_MetricValue {
	if m != nil {
		return m.MetricValue
	}
	return nil
}

func (x *Metric) GetValue() float64 {
	if x, ok := x.GetMetricValue().(*Metric_Value); ok {
		return x.Value
	}
	return 0
}

func (x *Metric) GetDelta() int64 {
	if x, ok := x.GetMetricValue().(*Metric_Delta); ok {
		return x.Delta
	}
	return 0
}

type isMetric_MetricValue interface {
	isMetric_MetricValue()
}

type Metric_Value struct {
	Value float64 `protobuf:"fixed64,3,opt,name=Value,proto3,oneof"`
}

type Metric_Delta struct {
	Delta int64 `protobuf:"varint,4,opt,name=Delta,proto3,oneof"`
}

func (*Metric_Value) isMetric_MetricValue() {}

func (*Metric_Delta) isMetric_MetricValue() {}

type SaveMetricsBatchRequestV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *SaveMetricsBatchRequestV1) Reset() {
	*x = SaveMetricsBatchRequestV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_v1_metrics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveMetricsBatchRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveMetricsBatchRequestV1) ProtoMessage() {}

func (x *SaveMetricsBatchRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_v1_metrics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveMetricsBatchRequestV1.ProtoReflect.Descriptor instead.
func (*SaveMetricsBatchRequestV1) Descriptor() ([]byte, []int) {
	return file_proto_metrics_v1_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *SaveMetricsBatchRequestV1) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type GetMetricRequestV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics *Metric `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *GetMetricRequestV1) Reset() {
	*x = GetMetricRequestV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_v1_metrics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequestV1) ProtoMessage() {}

func (x *GetMetricRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_v1_metrics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricRequestV1.ProtoReflect.Descriptor instead.
func (*GetMetricRequestV1) Descriptor() ([]byte, []int) {
	return file_proto_metrics_v1_metrics_proto_rawDescGZIP(), []int{2}
}

func (x *GetMetricRequestV1) GetMetrics() *Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type GetMetricResponseV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *GetMetricResponseV1) Reset() {
	*x = GetMetricResponseV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_v1_metrics_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponseV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponseV1) ProtoMessage() {}

func (x *GetMetricResponseV1) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_v1_metrics_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricResponseV1.ProtoReflect.Descriptor instead.
func (*GetMetricResponseV1) Descriptor() ([]byte, []int) {
	return file_proto_metrics_v1_metrics_proto_rawDescGZIP(), []int{3}
}

func (x *GetMetricResponseV1) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

var File_proto_metrics_v1_metrics_proto protoreflect.FileDescriptor

var file_proto_metrics_v1_metrics_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2f,
	0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbb, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x12, 0x17, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba,
	0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x49, 0x44, 0x12, 0x33, 0x0a, 0x05, 0x4d, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x16, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00,
	0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x22,
	0x24, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a,
	0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x55, 0x4e,
	0x54, 0x45, 0x52, 0x10, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x49, 0x0a, 0x19, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56,
	0x31, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22,
	0x42, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x56, 0x31, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x22, 0x41, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31, 0x12, 0x2a, 0x0a, 0x06, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x32, 0xb7, 0x01, 0x0a, 0x10, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x31, 0x12, 0x53, 0x0a, 0x12, 0x53,
	0x61, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x56,
	0x31, 0x12, 0x25, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x61, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x4e, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x56, 0x31, 0x12,
	0x1e, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x1a,
	0x1f, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31,
	0x42, 0x12, 0x5a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metrics_v1_metrics_proto_rawDescOnce sync.Once
	file_proto_metrics_v1_metrics_proto_rawDescData = file_proto_metrics_v1_metrics_proto_rawDesc
)

func file_proto_metrics_v1_metrics_proto_rawDescGZIP() []byte {
	file_proto_metrics_v1_metrics_proto_rawDescOnce.Do(func() {
		file_proto_metrics_v1_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metrics_v1_metrics_proto_rawDescData)
	})
	return file_proto_metrics_v1_metrics_proto_rawDescData
}

var file_proto_metrics_v1_metrics_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_metrics_v1_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_metrics_v1_metrics_proto_goTypes = []any{
	(Metric_MetricType)(0),            // 0: metrics.v1.Metric.MetricType
	(*Metric)(nil),                    // 1: metrics.v1.Metric
	(*SaveMetricsBatchRequestV1)(nil), // 2: metrics.v1.SaveMetricsBatchRequestV1
	(*GetMetricRequestV1)(nil),        // 3: metrics.v1.GetMetricRequestV1
	(*GetMetricResponseV1)(nil),       // 4: metrics.v1.GetMetricResponseV1
	(*emptypb.Empty)(nil),             // 5: google.protobuf.Empty
}
var file_proto_metrics_v1_metrics_proto_depIdxs = []int32{
	0, // 0: metrics.v1.Metric.MType:type_name -> metrics.v1.Metric.MetricType
	1, // 1: metrics.v1.SaveMetricsBatchRequestV1.metrics:type_name -> metrics.v1.Metric
	1, // 2: metrics.v1.GetMetricRequestV1.metrics:type_name -> metrics.v1.Metric
	1, // 3: metrics.v1.GetMetricResponseV1.metric:type_name -> metrics.v1.Metric
	2, // 4: metrics.v1.MetricsServiceV1.SaveMetricsBatchV1:input_type -> metrics.v1.SaveMetricsBatchRequestV1
	3, // 5: metrics.v1.MetricsServiceV1.GetMetricV1:input_type -> metrics.v1.GetMetricRequestV1
	5, // 6: metrics.v1.MetricsServiceV1.SaveMetricsBatchV1:output_type -> google.protobuf.Empty
	4, // 7: metrics.v1.MetricsServiceV1.GetMetricV1:output_type -> metrics.v1.GetMetricResponseV1
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_metrics_v1_metrics_proto_init() }
func file_proto_metrics_v1_metrics_proto_init() {
	if File_proto_metrics_v1_metrics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metrics_v1_metrics_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Metric); i {
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
		file_proto_metrics_v1_metrics_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SaveMetricsBatchRequestV1); i {
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
		file_proto_metrics_v1_metrics_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricRequestV1); i {
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
		file_proto_metrics_v1_metrics_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricResponseV1); i {
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
	file_proto_metrics_v1_metrics_proto_msgTypes[0].OneofWrappers = []any{
		(*Metric_Value)(nil),
		(*Metric_Delta)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metrics_v1_metrics_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metrics_v1_metrics_proto_goTypes,
		DependencyIndexes: file_proto_metrics_v1_metrics_proto_depIdxs,
		EnumInfos:         file_proto_metrics_v1_metrics_proto_enumTypes,
		MessageInfos:      file_proto_metrics_v1_metrics_proto_msgTypes,
	}.Build()
	File_proto_metrics_v1_metrics_proto = out.File
	file_proto_metrics_v1_metrics_proto_rawDesc = nil
	file_proto_metrics_v1_metrics_proto_goTypes = nil
	file_proto_metrics_v1_metrics_proto_depIdxs = nil
}
