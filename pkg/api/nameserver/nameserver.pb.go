// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.4
// source: nameserver.proto

package nameserver

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsDir bool   `protobuf:"varint,1,opt,name=is_dir,json=isDir,proto3" json:"is_dir,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Size  uint64 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetIsDir() bool {
	if x != nil {
		return x.IsDir
	}
	return false
}

func (x *Node) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Node) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

type LookupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *LookupRequest) Reset() {
	*x = LookupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupRequest) ProtoMessage() {}

func (x *LookupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupRequest.ProtoReflect.Descriptor instead.
func (*LookupRequest) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{1}
}

func (x *LookupRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type LookupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes []*Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

func (x *LookupResponse) Reset() {
	*x = LookupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupResponse) ProtoMessage() {}

func (x *LookupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupResponse.ProtoReflect.Descriptor instead.
func (*LookupResponse) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{2}
}

func (x *LookupResponse) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path  string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	IsDir bool   `protobuf:"varint,2,opt,name=is_dir,json=isDir,proto3" json:"is_dir,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{3}
}

func (x *CreateRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *CreateRequest) GetIsDir() bool {
	if x != nil {
		return x.IsDir
	}
	return false
}

type CreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateResponse) Reset() {
	*x = CreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateResponse) ProtoMessage() {}

func (x *CreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateResponse.ProtoReflect.Descriptor instead.
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{4}
}

type RemoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *RemoveRequest) Reset() {
	*x = RemoveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveRequest) ProtoMessage() {}

func (x *RemoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveRequest.ProtoReflect.Descriptor instead.
func (*RemoveRequest) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{5}
}

func (x *RemoveRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type RemoveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveResponse) Reset() {
	*x = RemoveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveResponse) ProtoMessage() {}

func (x *RemoveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveResponse.ProtoReflect.Descriptor instead.
func (*RemoveResponse) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{6}
}

type MapFSRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *MapFSRequest) Reset() {
	*x = MapFSRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapFSRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapFSRequest) ProtoMessage() {}

func (x *MapFSRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapFSRequest.ProtoReflect.Descriptor instead.
func (*MapFSRequest) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{7}
}

func (x *MapFSRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type MapFSResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fsurl string `protobuf:"bytes,1,opt,name=fsurl,proto3" json:"fsurl,omitempty"`
	Inode string `protobuf:"bytes,2,opt,name=inode,proto3" json:"inode,omitempty"`
}

func (x *MapFSResponse) Reset() {
	*x = MapFSResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nameserver_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapFSResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapFSResponse) ProtoMessage() {}

func (x *MapFSResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nameserver_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapFSResponse.ProtoReflect.Descriptor instead.
func (*MapFSResponse) Descriptor() ([]byte, []int) {
	return file_nameserver_proto_rawDescGZIP(), []int{8}
}

func (x *MapFSResponse) GetFsurl() string {
	if x != nil {
		return x.Fsurl
	}
	return ""
}

func (x *MapFSResponse) GetInode() string {
	if x != nil {
		return x.Inode
	}
	return ""
}

var File_nameserver_proto protoreflect.FileDescriptor

var file_nameserver_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x6e, 0x73, 0x61, 0x70, 0x69, 0x22, 0x45, 0x0a, 0x04, 0x4e, 0x6f, 0x64,
	0x65, 0x12, 0x15, 0x0a, 0x06, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x05, 0x69, 0x73, 0x44, 0x69, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x22, 0x23, 0x0a, 0x0d, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x33, 0x0a, 0x0e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x6e, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x4e,
	0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x22, 0x3a, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12,
	0x15, 0x0a, 0x06, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x05, 0x69, 0x73, 0x44, 0x69, 0x72, 0x22, 0x10, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x23, 0x0a, 0x0d, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x10, 0x0a,
	0x0e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x22, 0x0a, 0x0c, 0x4d, 0x61, 0x70, 0x46, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x22, 0x3b, 0x0a, 0x0d, 0x4d, 0x61, 0x70, 0x46, 0x53, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x73, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x73, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e,
	0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65,
	0x32, 0xed, 0x01, 0x0a, 0x0a, 0x4e, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12,
	0x37, 0x0a, 0x06, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x12, 0x14, 0x2e, 0x6e, 0x73, 0x61, 0x70,
	0x69, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x6e, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x14, 0x2e, 0x6e, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6e, 0x73, 0x61, 0x70, 0x69,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x37, 0x0a, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x14, 0x2e, 0x6e, 0x73,
	0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x6e, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x05, 0x4d, 0x61,
	0x70, 0x46, 0x53, 0x12, 0x13, 0x2e, 0x6e, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x61, 0x70, 0x46,
	0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x6e, 0x73, 0x61, 0x70, 0x69,
	0x2e, 0x4d, 0x61, 0x70, 0x46, 0x53, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d,
	0x65, 0x78, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x47, 0x6f, 0x2d, 0x76, 0x6e, 0x6f, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nameserver_proto_rawDescOnce sync.Once
	file_nameserver_proto_rawDescData = file_nameserver_proto_rawDesc
)

func file_nameserver_proto_rawDescGZIP() []byte {
	file_nameserver_proto_rawDescOnce.Do(func() {
		file_nameserver_proto_rawDescData = protoimpl.X.CompressGZIP(file_nameserver_proto_rawDescData)
	})
	return file_nameserver_proto_rawDescData
}

var file_nameserver_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_nameserver_proto_goTypes = []interface{}{
	(*Node)(nil),           // 0: nsapi.Node
	(*LookupRequest)(nil),  // 1: nsapi.LookupRequest
	(*LookupResponse)(nil), // 2: nsapi.LookupResponse
	(*CreateRequest)(nil),  // 3: nsapi.CreateRequest
	(*CreateResponse)(nil), // 4: nsapi.CreateResponse
	(*RemoveRequest)(nil),  // 5: nsapi.RemoveRequest
	(*RemoveResponse)(nil), // 6: nsapi.RemoveResponse
	(*MapFSRequest)(nil),   // 7: nsapi.MapFSRequest
	(*MapFSResponse)(nil),  // 8: nsapi.MapFSResponse
}
var file_nameserver_proto_depIdxs = []int32{
	0, // 0: nsapi.LookupResponse.nodes:type_name -> nsapi.Node
	1, // 1: nsapi.NameServer.Lookup:input_type -> nsapi.LookupRequest
	3, // 2: nsapi.NameServer.Create:input_type -> nsapi.CreateRequest
	5, // 3: nsapi.NameServer.Remove:input_type -> nsapi.RemoveRequest
	7, // 4: nsapi.NameServer.MapFS:input_type -> nsapi.MapFSRequest
	2, // 5: nsapi.NameServer.Lookup:output_type -> nsapi.LookupResponse
	4, // 6: nsapi.NameServer.Create:output_type -> nsapi.CreateResponse
	6, // 7: nsapi.NameServer.Remove:output_type -> nsapi.RemoveResponse
	8, // 8: nsapi.NameServer.MapFS:output_type -> nsapi.MapFSResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_nameserver_proto_init() }
func file_nameserver_proto_init() {
	if File_nameserver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nameserver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_nameserver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupRequest); i {
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
		file_nameserver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupResponse); i {
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
		file_nameserver_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_nameserver_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateResponse); i {
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
		file_nameserver_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveRequest); i {
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
		file_nameserver_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveResponse); i {
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
		file_nameserver_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapFSRequest); i {
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
		file_nameserver_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapFSResponse); i {
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
			RawDescriptor: file_nameserver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nameserver_proto_goTypes,
		DependencyIndexes: file_nameserver_proto_depIdxs,
		MessageInfos:      file_nameserver_proto_msgTypes,
	}.Build()
	File_nameserver_proto = out.File
	file_nameserver_proto_rawDesc = nil
	file_nameserver_proto_goTypes = nil
	file_nameserver_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// NameServerClient is the client API for NameServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NameServerClient interface {
	Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveResponse, error)
	MapFS(ctx context.Context, in *MapFSRequest, opts ...grpc.CallOption) (*MapFSResponse, error)
}

type nameServerClient struct {
	cc grpc.ClientConnInterface
}

func NewNameServerClient(cc grpc.ClientConnInterface) NameServerClient {
	return &nameServerClient{cc}
}

func (c *nameServerClient) Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error) {
	out := new(LookupResponse)
	err := c.cc.Invoke(ctx, "/nsapi.NameServer/Lookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameServerClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/nsapi.NameServer/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameServerClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveResponse, error) {
	out := new(RemoveResponse)
	err := c.cc.Invoke(ctx, "/nsapi.NameServer/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameServerClient) MapFS(ctx context.Context, in *MapFSRequest, opts ...grpc.CallOption) (*MapFSResponse, error) {
	out := new(MapFSResponse)
	err := c.cc.Invoke(ctx, "/nsapi.NameServer/MapFS", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NameServerServer is the server API for NameServer service.
type NameServerServer interface {
	Lookup(context.Context, *LookupRequest) (*LookupResponse, error)
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Remove(context.Context, *RemoveRequest) (*RemoveResponse, error)
	MapFS(context.Context, *MapFSRequest) (*MapFSResponse, error)
}

// UnimplementedNameServerServer can be embedded to have forward compatible implementations.
type UnimplementedNameServerServer struct {
}

func (*UnimplementedNameServerServer) Lookup(context.Context, *LookupRequest) (*LookupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
}
func (*UnimplementedNameServerServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedNameServerServer) Remove(context.Context, *RemoveRequest) (*RemoveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedNameServerServer) MapFS(context.Context, *MapFSRequest) (*MapFSResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapFS not implemented")
}

func RegisterNameServerServer(s *grpc.Server, srv NameServerServer) {
	s.RegisterService(&_NameServer_serviceDesc, srv)
}

func _NameServer_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameServerServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nsapi.NameServer/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameServerServer).Lookup(ctx, req.(*LookupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameServer_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameServerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nsapi.NameServer/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameServerServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameServer_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameServerServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nsapi.NameServer/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameServerServer).Remove(ctx, req.(*RemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameServer_MapFS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MapFSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameServerServer).MapFS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nsapi.NameServer/MapFS",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameServerServer).MapFS(ctx, req.(*MapFSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NameServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nsapi.NameServer",
	HandlerType: (*NameServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Lookup",
			Handler:    _NameServer_Lookup_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _NameServer_Create_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _NameServer_Remove_Handler,
		},
		{
			MethodName: "MapFS",
			Handler:    _NameServer_MapFS_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nameserver.proto",
}
