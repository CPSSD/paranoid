// Code generated by protoc-gen-go.
// source: paranoidnetwork/paranoid.proto
// DO NOT EDIT!

/*
Package paranoid is a generated protocol buffer package.

It is generated from these files:
	paranoidnetwork/paranoid.proto

It has these top-level messages:
	EmptyMessage
	PingRequest
	CreatRequest
	WriteRequest
	WriteResponse
	LinkRequest
	UnlinkRequest
	RenameRequest
	TruncateRequest
	UtimesRequest
	ChmodRequest
	MkdirRequest
	RmdirRequest
	KeyPiece
*/
package paranoid

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EmptyMessage struct {
}

func (m *EmptyMessage) Reset()                    { *m = EmptyMessage{} }
func (m *EmptyMessage) String() string            { return proto.CompactTextString(m) }
func (*EmptyMessage) ProtoMessage()               {}
func (*EmptyMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type PingRequest struct {
	Ip         string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Port       string `protobuf:"bytes,2,opt,name=port" json:"port,omitempty"`
	CommonName string `protobuf:"bytes,3,opt,name=common_name" json:"common_name,omitempty"`
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type CreatRequest struct {
	Path        string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Permissions uint32 `protobuf:"varint,2,opt,name=permissions" json:"permissions,omitempty"`
}

func (m *CreatRequest) Reset()                    { *m = CreatRequest{} }
func (m *CreatRequest) String() string            { return proto.CompactTextString(m) }
func (*CreatRequest) ProtoMessage()               {}
func (*CreatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type WriteRequest struct {
	Path   string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Data   []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Offset uint64 `protobuf:"varint,3,opt,name=offset" json:"offset,omitempty"`
	Length uint64 `protobuf:"varint,4,opt,name=length" json:"length,omitempty"`
}

func (m *WriteRequest) Reset()                    { *m = WriteRequest{} }
func (m *WriteRequest) String() string            { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()               {}
func (*WriteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type WriteResponse struct {
	BytesWritten uint64 `protobuf:"varint,1,opt,name=bytes_written" json:"bytes_written,omitempty"`
}

func (m *WriteResponse) Reset()                    { *m = WriteResponse{} }
func (m *WriteResponse) String() string            { return proto.CompactTextString(m) }
func (*WriteResponse) ProtoMessage()               {}
func (*WriteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type LinkRequest struct {
	OldPath string `protobuf:"bytes,1,opt,name=old_path" json:"old_path,omitempty"`
	NewPath string `protobuf:"bytes,2,opt,name=new_path" json:"new_path,omitempty"`
}

func (m *LinkRequest) Reset()                    { *m = LinkRequest{} }
func (m *LinkRequest) String() string            { return proto.CompactTextString(m) }
func (*LinkRequest) ProtoMessage()               {}
func (*LinkRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type UnlinkRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *UnlinkRequest) Reset()                    { *m = UnlinkRequest{} }
func (m *UnlinkRequest) String() string            { return proto.CompactTextString(m) }
func (*UnlinkRequest) ProtoMessage()               {}
func (*UnlinkRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type RenameRequest struct {
	OldPath string `protobuf:"bytes,1,opt,name=old_path" json:"old_path,omitempty"`
	NewPath string `protobuf:"bytes,2,opt,name=new_path" json:"new_path,omitempty"`
}

func (m *RenameRequest) Reset()                    { *m = RenameRequest{} }
func (m *RenameRequest) String() string            { return proto.CompactTextString(m) }
func (*RenameRequest) ProtoMessage()               {}
func (*RenameRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type TruncateRequest struct {
	Path   string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Length uint64 `protobuf:"varint,2,opt,name=length" json:"length,omitempty"`
}

func (m *TruncateRequest) Reset()                    { *m = TruncateRequest{} }
func (m *TruncateRequest) String() string            { return proto.CompactTextString(m) }
func (*TruncateRequest) ProtoMessage()               {}
func (*TruncateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type UtimesRequest struct {
	Path              string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	AccessSeconds     int64  `protobuf:"varint,2,opt,name=access_seconds" json:"access_seconds,omitempty"`
	AccessNanoseconds int64  `protobuf:"varint,3,opt,name=access_nanoseconds" json:"access_nanoseconds,omitempty"`
	ModifySeconds     int64  `protobuf:"varint,4,opt,name=modify_seconds" json:"modify_seconds,omitempty"`
	ModifyNanoseconds int64  `protobuf:"varint,5,opt,name=modify_nanoseconds" json:"modify_nanoseconds,omitempty"`
}

func (m *UtimesRequest) Reset()                    { *m = UtimesRequest{} }
func (m *UtimesRequest) String() string            { return proto.CompactTextString(m) }
func (*UtimesRequest) ProtoMessage()               {}
func (*UtimesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type ChmodRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Mode uint32 `protobuf:"varint,2,opt,name=mode" json:"mode,omitempty"`
}

func (m *ChmodRequest) Reset()                    { *m = ChmodRequest{} }
func (m *ChmodRequest) String() string            { return proto.CompactTextString(m) }
func (*ChmodRequest) ProtoMessage()               {}
func (*ChmodRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type MkdirRequest struct {
	Directory string `protobuf:"bytes,1,opt,name=directory" json:"directory,omitempty"`
	Mode      uint32 `protobuf:"varint,2,opt,name=mode" json:"mode,omitempty"`
}

func (m *MkdirRequest) Reset()                    { *m = MkdirRequest{} }
func (m *MkdirRequest) String() string            { return proto.CompactTextString(m) }
func (*MkdirRequest) ProtoMessage()               {}
func (*MkdirRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type RmdirRequest struct {
	Directory string `protobuf:"bytes,1,opt,name=directory" json:"directory,omitempty"`
}

func (m *RmdirRequest) Reset()                    { *m = RmdirRequest{} }
func (m *RmdirRequest) String() string            { return proto.CompactTextString(m) }
func (*RmdirRequest) ProtoMessage()               {}
func (*RmdirRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type KeyPiece struct {
	Data              []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	ParentFingerprint []byte `protobuf:"bytes,2,opt,name=parent_fingerprint,proto3" json:"parent_fingerprint,omitempty"`
	Prime             []byte `protobuf:"bytes,3,opt,name=prime,proto3" json:"prime,omitempty"`
	Seq               int64  `protobuf:"varint,4,opt,name=seq" json:"seq,omitempty"`
	// The Node data for the node who owns this KeyPiece
	OwnerNode *PingRequest `protobuf:"bytes,5,opt,name=owner_node" json:"owner_node,omitempty"`
}

func (m *KeyPiece) Reset()                    { *m = KeyPiece{} }
func (m *KeyPiece) String() string            { return proto.CompactTextString(m) }
func (*KeyPiece) ProtoMessage()               {}
func (*KeyPiece) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *KeyPiece) GetOwnerNode() *PingRequest {
	if m != nil {
		return m.OwnerNode
	}
	return nil
}

func init() {
	proto.RegisterType((*EmptyMessage)(nil), "paranoid.EmptyMessage")
	proto.RegisterType((*PingRequest)(nil), "paranoid.PingRequest")
	proto.RegisterType((*CreatRequest)(nil), "paranoid.CreatRequest")
	proto.RegisterType((*WriteRequest)(nil), "paranoid.WriteRequest")
	proto.RegisterType((*WriteResponse)(nil), "paranoid.WriteResponse")
	proto.RegisterType((*LinkRequest)(nil), "paranoid.LinkRequest")
	proto.RegisterType((*UnlinkRequest)(nil), "paranoid.UnlinkRequest")
	proto.RegisterType((*RenameRequest)(nil), "paranoid.RenameRequest")
	proto.RegisterType((*TruncateRequest)(nil), "paranoid.TruncateRequest")
	proto.RegisterType((*UtimesRequest)(nil), "paranoid.UtimesRequest")
	proto.RegisterType((*ChmodRequest)(nil), "paranoid.ChmodRequest")
	proto.RegisterType((*MkdirRequest)(nil), "paranoid.MkdirRequest")
	proto.RegisterType((*RmdirRequest)(nil), "paranoid.RmdirRequest")
	proto.RegisterType((*KeyPiece)(nil), "paranoid.KeyPiece")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for ParanoidNetwork service

type ParanoidNetworkClient interface {
	// Used for health checking and discovery. Sends the IP and port of the
	// PFSD instance running on the client.
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	// Filesystem calls
	Creat(ctx context.Context, in *CreatRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Link(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Symlink(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Unlink(ctx context.Context, in *UnlinkRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Rename(ctx context.Context, in *RenameRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Truncate(ctx context.Context, in *TruncateRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Utimes(ctx context.Context, in *UtimesRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Chmod(ctx context.Context, in *ChmodRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Mkdir(ctx context.Context, in *MkdirRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	Rmdir(ctx context.Context, in *RmdirRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
	// Cryptography calls
	Lock(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error)
	Unlock(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error)
	SendKeyPiece(ctx context.Context, in *KeyPiece, opts ...grpc.CallOption) (*EmptyMessage, error)
	RequestKeyPiece(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*KeyPiece, error)
}

type paranoidNetworkClient struct {
	cc *grpc.ClientConn
}

func NewParanoidNetworkClient(cc *grpc.ClientConn) ParanoidNetworkClient {
	return &paranoidNetworkClient{cc}
}

func (c *paranoidNetworkClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Creat(ctx context.Context, in *CreatRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Creat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Write", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Link(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Link", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Symlink(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Symlink", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Unlink(ctx context.Context, in *UnlinkRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Unlink", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Rename(ctx context.Context, in *RenameRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Rename", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Truncate(ctx context.Context, in *TruncateRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Truncate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Utimes(ctx context.Context, in *UtimesRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Utimes", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Chmod(ctx context.Context, in *ChmodRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Chmod", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Mkdir(ctx context.Context, in *MkdirRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Mkdir", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Rmdir(ctx context.Context, in *RmdirRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Rmdir", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Lock(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Lock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) Unlock(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/Unlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) SendKeyPiece(ctx context.Context, in *KeyPiece, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/SendKeyPiece", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paranoidNetworkClient) RequestKeyPiece(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*KeyPiece, error) {
	out := new(KeyPiece)
	err := grpc.Invoke(ctx, "/paranoid.ParanoidNetwork/RequestKeyPiece", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ParanoidNetwork service

type ParanoidNetworkServer interface {
	// Used for health checking and discovery. Sends the IP and port of the
	// PFSD instance running on the client.
	Ping(context.Context, *PingRequest) (*EmptyMessage, error)
	// Filesystem calls
	Creat(context.Context, *CreatRequest) (*EmptyMessage, error)
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
	Link(context.Context, *LinkRequest) (*EmptyMessage, error)
	Symlink(context.Context, *LinkRequest) (*EmptyMessage, error)
	Unlink(context.Context, *UnlinkRequest) (*EmptyMessage, error)
	Rename(context.Context, *RenameRequest) (*EmptyMessage, error)
	Truncate(context.Context, *TruncateRequest) (*EmptyMessage, error)
	Utimes(context.Context, *UtimesRequest) (*EmptyMessage, error)
	Chmod(context.Context, *ChmodRequest) (*EmptyMessage, error)
	Mkdir(context.Context, *MkdirRequest) (*EmptyMessage, error)
	Rmdir(context.Context, *RmdirRequest) (*EmptyMessage, error)
	// Cryptography calls
	Lock(context.Context, *EmptyMessage) (*EmptyMessage, error)
	Unlock(context.Context, *EmptyMessage) (*EmptyMessage, error)
	SendKeyPiece(context.Context, *KeyPiece) (*EmptyMessage, error)
	RequestKeyPiece(context.Context, *PingRequest) (*KeyPiece, error)
}

func RegisterParanoidNetworkServer(s *grpc.Server, srv ParanoidNetworkServer) {
	s.RegisterService(&_ParanoidNetwork_serviceDesc, srv)
}

func _ParanoidNetwork_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Ping(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Creat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CreatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Creat(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Write(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Link_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(LinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Link(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Symlink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(LinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Symlink(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Unlink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(UnlinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Unlink(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Rename_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(RenameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Rename(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Truncate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(TruncateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Truncate(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Utimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(UtimesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Utimes(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Chmod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(ChmodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Chmod(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Mkdir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(MkdirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Mkdir(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Rmdir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(RmdirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Rmdir(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Lock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Lock(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_Unlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).Unlock(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_SendKeyPiece_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(KeyPiece)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).SendKeyPiece(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ParanoidNetwork_RequestKeyPiece_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ParanoidNetworkServer).RequestKeyPiece(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _ParanoidNetwork_serviceDesc = grpc.ServiceDesc{
	ServiceName: "paranoid.ParanoidNetwork",
	HandlerType: (*ParanoidNetworkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _ParanoidNetwork_Ping_Handler,
		},
		{
			MethodName: "Creat",
			Handler:    _ParanoidNetwork_Creat_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _ParanoidNetwork_Write_Handler,
		},
		{
			MethodName: "Link",
			Handler:    _ParanoidNetwork_Link_Handler,
		},
		{
			MethodName: "Symlink",
			Handler:    _ParanoidNetwork_Symlink_Handler,
		},
		{
			MethodName: "Unlink",
			Handler:    _ParanoidNetwork_Unlink_Handler,
		},
		{
			MethodName: "Rename",
			Handler:    _ParanoidNetwork_Rename_Handler,
		},
		{
			MethodName: "Truncate",
			Handler:    _ParanoidNetwork_Truncate_Handler,
		},
		{
			MethodName: "Utimes",
			Handler:    _ParanoidNetwork_Utimes_Handler,
		},
		{
			MethodName: "Chmod",
			Handler:    _ParanoidNetwork_Chmod_Handler,
		},
		{
			MethodName: "Mkdir",
			Handler:    _ParanoidNetwork_Mkdir_Handler,
		},
		{
			MethodName: "Rmdir",
			Handler:    _ParanoidNetwork_Rmdir_Handler,
		},
		{
			MethodName: "Lock",
			Handler:    _ParanoidNetwork_Lock_Handler,
		},
		{
			MethodName: "Unlock",
			Handler:    _ParanoidNetwork_Unlock_Handler,
		},
		{
			MethodName: "SendKeyPiece",
			Handler:    _ParanoidNetwork_SendKeyPiece_Handler,
		},
		{
			MethodName: "RequestKeyPiece",
			Handler:    _ParanoidNetwork_RequestKeyPiece_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 636 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x95, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0xc7, 0x9b, 0xc6, 0x69, 0xd3, 0x89, 0xdd, 0x3e, 0xcf, 0xa2, 0x96, 0x52, 0x09, 0x04, 0x3e,
	0x20, 0xe0, 0xd0, 0x8a, 0xf6, 0xc0, 0x5b, 0x05, 0x07, 0xc4, 0x05, 0x28, 0xaa, 0x5a, 0x10, 0x47,
	0xcb, 0xb5, 0x27, 0xe9, 0x2a, 0xf1, 0xae, 0xbb, 0xbb, 0x55, 0x64, 0x71, 0xe6, 0x13, 0xf2, 0x85,
	0xd8, 0x5d, 0xdb, 0xc9, 0x1a, 0x12, 0x57, 0xf4, 0x38, 0x2f, 0xff, 0xd9, 0x19, 0xcf, 0xfc, 0x64,
	0x78, 0x90, 0xc7, 0x22, 0x66, 0x9c, 0xa6, 0x0c, 0xd5, 0x94, 0x8b, 0xf1, 0x41, 0x6d, 0xef, 0xe7,
	0x82, 0x2b, 0x4e, 0xfa, 0xb5, 0x1d, 0x6e, 0x82, 0xff, 0x21, 0xcb, 0x55, 0x71, 0x82, 0x52, 0xc6,
	0x23, 0x0c, 0x8f, 0x61, 0x70, 0x4a, 0xd9, 0xe8, 0x0c, 0xaf, 0xae, 0x51, 0x2a, 0x02, 0xb0, 0x4a,
	0xf3, 0xdd, 0xce, 0xc3, 0xce, 0x93, 0x0d, 0xe2, 0x83, 0x97, 0x73, 0xa1, 0x76, 0x57, 0xad, 0x75,
	0x07, 0x06, 0x09, 0xcf, 0x32, 0xce, 0x22, 0x16, 0x67, 0xb8, 0xdb, 0x35, 0xce, 0xf0, 0x39, 0xf8,
	0xef, 0x05, 0xc6, 0xaa, 0x96, 0x1b, 0x49, 0xac, 0x2e, 0xab, 0x02, 0x5a, 0x92, 0xa3, 0xc8, 0xa8,
	0x94, 0x94, 0x33, 0x69, 0xeb, 0x04, 0xe1, 0x47, 0xf0, 0xbf, 0x0b, 0xaa, 0x70, 0xb1, 0x44, 0x5b,
	0x69, 0xac, 0x62, 0x9b, 0xeb, 0x93, 0x4d, 0x58, 0xe3, 0xc3, 0xa1, 0x44, 0x65, 0x9f, 0xf3, 0x8c,
	0x3d, 0x41, 0x36, 0xd2, 0xd9, 0x9e, 0xb1, 0xc3, 0xc7, 0x10, 0x54, 0xb5, 0x64, 0xae, 0x5f, 0x40,
	0xb2, 0x0d, 0xc1, 0x45, 0xa1, 0x50, 0x46, 0x53, 0xed, 0x56, 0xc8, 0x6c, 0x55, 0x4f, 0xb7, 0x39,
	0xf8, 0x4c, 0xd9, 0xb8, 0x7e, 0xf2, 0x3f, 0xe8, 0xf3, 0x49, 0x1a, 0x39, 0xcf, 0x6a, 0x0f, 0xc3,
	0x69, 0xe9, 0xb1, 0xe3, 0x86, 0xf7, 0x21, 0xf8, 0xc6, 0x26, 0x8e, 0xa8, 0xd1, 0x67, 0x78, 0x04,
	0xc1, 0x19, 0x9a, 0x0f, 0xf1, 0x2f, 0x35, 0x0f, 0x60, 0xeb, 0xab, 0xb8, 0x66, 0x49, 0xbc, 0x6c,
	0xfa, 0xf9, 0x7c, 0xab, 0xb6, 0xef, 0x9f, 0x1d, 0xdd, 0x85, 0xa2, 0x19, 0xca, 0xc5, 0xf9, 0x3b,
	0xb0, 0x19, 0x27, 0x89, 0xde, 0x64, 0x24, 0x31, 0xe1, 0x2c, 0x2d, 0xbf, 0x71, 0x97, 0xec, 0x01,
	0xa9, 0xfc, 0x4c, 0x6f, 0xbd, 0x8e, 0x75, 0x6d, 0x4c, 0x6b, 0x32, 0x9e, 0xd2, 0x61, 0x31, 0xd3,
	0x78, 0xb5, 0xa6, 0xf2, 0xbb, 0x9a, 0x9e, 0x89, 0x85, 0xcf, 0xf4, 0x9a, 0x2f, 0x75, 0x74, 0xe9,
	0xce, 0x74, 0x0c, 0xab, 0xfd, 0x1e, 0x80, 0x7f, 0x32, 0x4e, 0xa9, 0xa8, 0x73, 0xff, 0x87, 0x0d,
	0x6d, 0x61, 0xa2, 0xb8, 0x28, 0x16, 0x0a, 0x1e, 0x81, 0x7f, 0x96, 0xb5, 0x0a, 0xc2, 0x1f, 0xd0,
	0xff, 0x84, 0xc5, 0x29, 0xc5, 0x04, 0x67, 0x17, 0xd2, 0xb1, 0x17, 0xa2, 0xbb, 0xd6, 0xa7, 0x8d,
	0x4c, 0x45, 0x43, 0x7d, 0xc5, 0x28, 0x72, 0x41, 0x99, 0xaa, 0xae, 0x27, 0x80, 0x9e, 0x36, 0xab,
	0x5b, 0xf5, 0xc9, 0x00, 0xba, 0x12, 0xaf, 0xaa, 0x69, 0x9f, 0x02, 0xf0, 0x29, 0x43, 0x11, 0x31,
	0xd3, 0x88, 0x99, 0x72, 0x70, 0xb8, 0xbd, 0x3f, 0xa3, 0xc6, 0x41, 0xe2, 0xf0, 0xd7, 0x3a, 0x6c,
	0x9d, 0x56, 0x81, 0x2f, 0x25, 0x5e, 0xe4, 0x05, 0x78, 0x26, 0x85, 0x2c, 0x96, 0xec, 0xed, 0xcc,
	0xdd, 0x0d, 0xd8, 0x56, 0xc8, 0x2b, 0xe8, 0x59, 0x60, 0x88, 0x93, 0xe2, 0x12, 0xd4, 0x22, 0x7d,
	0x0d, 0x3d, 0x7b, 0xec, 0xae, 0xd4, 0x25, 0x69, 0xef, 0xee, 0x5f, 0xfe, 0x92, 0x0a, 0xad, 0xd5,
	0xfd, 0x1a, 0x00, 0xdc, 0x7e, 0x1d, 0x20, 0x5a, 0x1f, 0x5d, 0x3f, 0x2f, 0xb2, 0xc9, 0xad, 0xb4,
	0x6f, 0x60, 0xad, 0x44, 0x88, 0x38, 0x9d, 0x35, 0xa0, 0x6a, 0x17, 0x97, 0x80, 0xb9, 0xe2, 0x06,
	0x72, 0x2d, 0xe2, 0x77, 0xd0, 0xaf, 0x41, 0x23, 0xf7, 0xe6, 0x59, 0x7f, 0xc0, 0x77, 0x43, 0xeb,
	0x96, 0xbb, 0x46, 0xeb, 0x2e, 0x89, 0x37, 0xec, 0xd8, 0xd0, 0xd2, 0xd8, 0xb1, 0x83, 0x4f, 0xbb,
	0xd4, 0xc2, 0xe3, 0x4a, 0x5d, 0x9a, 0xda, 0xa5, 0x16, 0x23, 0x57, 0xea, 0x72, 0xd5, 0x22, 0x7d,
	0xa9, 0xaf, 0x83, 0x27, 0x63, 0xb2, 0x24, 0xa3, 0xf5, 0x3c, 0xcc, 0x8a, 0x6f, 0xa7, 0x3d, 0x06,
	0xff, 0x1c, 0x59, 0x3a, 0x03, 0x9b, 0xcc, 0x33, 0x6b, 0x5f, 0x8b, 0xfa, 0x2d, 0x6c, 0x55, 0x83,
	0xcd, 0x0a, 0x2c, 0x81, 0x71, 0x41, 0xdd, 0x70, 0xe5, 0x62, 0xcd, 0xfe, 0x18, 0x8f, 0x7e, 0x07,
	0x00, 0x00, 0xff, 0xff, 0xe6, 0x8c, 0x16, 0xa0, 0x3a, 0x07, 0x00, 0x00,
}
