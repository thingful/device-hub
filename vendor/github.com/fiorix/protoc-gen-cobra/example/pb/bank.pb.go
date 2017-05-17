// Code generated by protoc-gen-gogo.
// source: bank.proto
// DO NOT EDIT!

/*
	Package pb is a generated protocol buffer package.

	It is generated from these files:
		bank.proto
		cache.proto
		timer.proto

	It has these top-level messages:
		DepositRequest
		DepositReply
		SetRequest
		SetResponse
		GetRequest
		GetResponse
		TickRequest
		TickResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DepositRequest struct {
	Account string  `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	Amount  float64 `protobuf:"fixed64,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *DepositRequest) Reset()                    { *m = DepositRequest{} }
func (m *DepositRequest) String() string            { return proto.CompactTextString(m) }
func (*DepositRequest) ProtoMessage()               {}
func (*DepositRequest) Descriptor() ([]byte, []int) { return fileDescriptorBank, []int{0} }

type DepositReply struct {
	Account string  `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	Balance float64 `protobuf:"fixed64,2,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (m *DepositReply) Reset()                    { *m = DepositReply{} }
func (m *DepositReply) String() string            { return proto.CompactTextString(m) }
func (*DepositReply) ProtoMessage()               {}
func (*DepositReply) Descriptor() ([]byte, []int) { return fileDescriptorBank, []int{1} }

func init() {
	proto.RegisterType((*DepositRequest)(nil), "pb.DepositRequest")
	proto.RegisterType((*DepositReply)(nil), "pb.DepositReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Bank service

type BankClient interface {
	Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositReply, error)
}

type bankClient struct {
	cc *grpc.ClientConn
}

func NewBankClient(cc *grpc.ClientConn) BankClient {
	return &bankClient{cc}
}

func (c *bankClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositReply, error) {
	out := new(DepositReply)
	err := grpc.Invoke(ctx, "/pb.Bank/Deposit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Bank service

type BankServer interface {
	Deposit(context.Context, *DepositRequest) (*DepositReply, error)
}

func RegisterBankServer(s *grpc.Server, srv BankServer) {
	s.RegisterService(&_Bank_serviceDesc, srv)
}

func _Bank_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Bank/Deposit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Bank_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Bank",
	HandlerType: (*BankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deposit",
			Handler:    _Bank_Deposit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bank.proto",
}

func (m *DepositRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DepositRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Account) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintBank(dAtA, i, uint64(len(m.Account)))
		i += copy(dAtA[i:], m.Account)
	}
	if m.Amount != 0 {
		dAtA[i] = 0x11
		i++
		i = encodeFixed64Bank(dAtA, i, uint64(math.Float64bits(float64(m.Amount))))
	}
	return i, nil
}

func (m *DepositReply) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DepositReply) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Account) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintBank(dAtA, i, uint64(len(m.Account)))
		i += copy(dAtA[i:], m.Account)
	}
	if m.Balance != 0 {
		dAtA[i] = 0x11
		i++
		i = encodeFixed64Bank(dAtA, i, uint64(math.Float64bits(float64(m.Balance))))
	}
	return i, nil
}

func encodeFixed64Bank(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Bank(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintBank(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *DepositRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sovBank(uint64(l))
	}
	if m.Amount != 0 {
		n += 9
	}
	return n
}

func (m *DepositReply) Size() (n int) {
	var l int
	_ = l
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sovBank(uint64(l))
	}
	if m.Balance != 0 {
		n += 9
	}
	return n
}

func sovBank(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozBank(x uint64) (n int) {
	return sovBank(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DepositRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBank
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DepositRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DepositRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBank
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBank
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Account = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += 8
			v = uint64(dAtA[iNdEx-8])
			v |= uint64(dAtA[iNdEx-7]) << 8
			v |= uint64(dAtA[iNdEx-6]) << 16
			v |= uint64(dAtA[iNdEx-5]) << 24
			v |= uint64(dAtA[iNdEx-4]) << 32
			v |= uint64(dAtA[iNdEx-3]) << 40
			v |= uint64(dAtA[iNdEx-2]) << 48
			v |= uint64(dAtA[iNdEx-1]) << 56
			m.Amount = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipBank(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBank
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DepositReply) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBank
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DepositReply: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DepositReply: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBank
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBank
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Account = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += 8
			v = uint64(dAtA[iNdEx-8])
			v |= uint64(dAtA[iNdEx-7]) << 8
			v |= uint64(dAtA[iNdEx-6]) << 16
			v |= uint64(dAtA[iNdEx-5]) << 24
			v |= uint64(dAtA[iNdEx-4]) << 32
			v |= uint64(dAtA[iNdEx-3]) << 40
			v |= uint64(dAtA[iNdEx-2]) << 48
			v |= uint64(dAtA[iNdEx-1]) << 56
			m.Balance = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipBank(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBank
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBank(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBank
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBank
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBank
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthBank
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowBank
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipBank(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthBank = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBank   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("bank.proto", fileDescriptorBank) }

var fileDescriptorBank = []byte{
	// 173 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4a, 0xcc, 0xcb,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x72, 0xe2, 0xe2, 0x73, 0x49,
	0x2d, 0xc8, 0x2f, 0xce, 0x2c, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x92, 0xe0, 0x62,
	0x4f, 0x4c, 0x4e, 0xce, 0x2f, 0xcd, 0x2b, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x71,
	0x85, 0xc4, 0xb8, 0xd8, 0x12, 0x73, 0xc1, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x8c, 0x41, 0x50, 0x9e,
	0x92, 0x13, 0x17, 0x0f, 0xdc, 0x8c, 0x82, 0x9c, 0x4a, 0x3c, 0x26, 0x48, 0x70, 0xb1, 0x27, 0x25,
	0xe6, 0x24, 0xe6, 0x25, 0xa7, 0x42, 0x8d, 0x80, 0x71, 0x8d, 0xcc, 0xb9, 0x58, 0x9c, 0x12, 0xf3,
	0xb2, 0x85, 0xf4, 0xb9, 0xd8, 0xa1, 0x66, 0x09, 0x09, 0xe9, 0x15, 0x24, 0xe9, 0xa1, 0x3a, 0x4e,
	0x4a, 0x00, 0x45, 0xac, 0x20, 0xa7, 0xd2, 0x49, 0xe0, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4,
	0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf1, 0x58, 0x8e, 0x21, 0x89, 0x0d, 0xec, 0x3b, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xa5, 0x6a, 0x33, 0xe8, 0xeb, 0x00, 0x00, 0x00,
}