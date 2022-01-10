// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mainchain/enterprise/v1/genesis.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the enterprise module's genesis state.
type GenesisState struct {
	// params defines all the paramaters of the module.
	Params                  Params                       `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	StartingPurchaseOrderId uint64                       `protobuf:"varint,2,opt,name=starting_purchase_order_id,json=startingPurchaseOrderId,proto3" json:"starting_purchase_order_id,omitempty"`
	PurchaseOrders          []EnterpriseUndPurchaseOrder `protobuf:"bytes,3,rep,name=purchase_orders,json=purchaseOrders,proto3" json:"purchase_orders"`
	LockedUnd               []LockedUnd                  `protobuf:"bytes,4,rep,name=locked_und,json=lockedUnd,proto3" json:"locked_und"`
	TotalLocked             types.Coin                   `protobuf:"bytes,5,opt,name=total_locked,json=totalLocked,proto3" json:"total_locked"`
	Whitelist               []string                     `protobuf:"bytes,6,rep,name=whitelist,proto3" json:"whitelist,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfcf11da3dee12f2, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetStartingPurchaseOrderId() uint64 {
	if m != nil {
		return m.StartingPurchaseOrderId
	}
	return 0
}

func (m *GenesisState) GetPurchaseOrders() []EnterpriseUndPurchaseOrder {
	if m != nil {
		return m.PurchaseOrders
	}
	return nil
}

func (m *GenesisState) GetLockedUnd() []LockedUnd {
	if m != nil {
		return m.LockedUnd
	}
	return nil
}

func (m *GenesisState) GetTotalLocked() types.Coin {
	if m != nil {
		return m.TotalLocked
	}
	return types.Coin{}
}

func (m *GenesisState) GetWhitelist() []string {
	if m != nil {
		return m.Whitelist
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "mainchain.enterprise.v1.GenesisState")
}

func init() {
	proto.RegisterFile("mainchain/enterprise/v1/genesis.proto", fileDescriptor_dfcf11da3dee12f2)
}

var fileDescriptor_dfcf11da3dee12f2 = []byte{
	// 397 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0xeb, 0xd3, 0x30,
	0x00, 0xc5, 0x5b, 0x3b, 0x07, 0xcb, 0x86, 0x42, 0x11, 0x56, 0x87, 0x74, 0x65, 0x20, 0xf4, 0x62,
	0xca, 0xb6, 0x83, 0x07, 0xf1, 0x32, 0x91, 0x21, 0x08, 0xce, 0xc9, 0x2e, 0x5e, 0x4a, 0xda, 0xc4,
	0x36, 0xac, 0x4d, 0x4a, 0x92, 0x4e, 0xfd, 0x2f, 0xbc, 0xfa, 0x1f, 0xed, 0xb8, 0xa3, 0x27, 0xf9,
	0xb2, 0xfd, 0x23, 0x5f, 0x9a, 0x76, 0xbf, 0x0e, 0xbd, 0x25, 0xaf, 0x9f, 0xf7, 0x78, 0xf4, 0x05,
	0xbc, 0xce, 0x11, 0x65, 0x71, 0x8a, 0x28, 0x0b, 0x08, 0x53, 0x44, 0x14, 0x82, 0x4a, 0x12, 0xec,
	0xa6, 0x41, 0x42, 0x18, 0x91, 0x54, 0xc2, 0x42, 0x70, 0xc5, 0xed, 0xe1, 0x05, 0x83, 0x57, 0x0c,
	0xee, 0xa6, 0x23, 0xbf, 0xcd, 0x7f, 0x83, 0xe9, 0x88, 0xd1, 0x24, 0xe6, 0x32, 0xe7, 0x32, 0x94,
	0x78, 0x1b, 0x44, 0x48, 0x43, 0x11, 0x51, 0x68, 0x1a, 0xc4, 0x9c, 0xb2, 0x86, 0x79, 0x91, 0xf0,
	0x84, 0xeb, 0x63, 0x50, 0x9d, 0x6a, 0x75, 0xf2, 0xd7, 0x02, 0x83, 0x65, 0x5d, 0xe7, 0x9b, 0x42,
	0x8a, 0xd8, 0xef, 0x41, 0xb7, 0x40, 0x02, 0xe5, 0xd2, 0x31, 0x3d, 0xd3, 0xef, 0xcf, 0xc6, 0xb0,
	0xa5, 0x1e, 0x5c, 0x69, 0x6c, 0xd1, 0xd9, 0xff, 0x1f, 0x1b, 0xeb, 0xc6, 0x64, 0xbf, 0x03, 0x23,
	0xa9, 0x90, 0x50, 0x94, 0x25, 0x61, 0x51, 0x8a, 0x38, 0x45, 0x92, 0x84, 0x5c, 0x60, 0x22, 0x42,
	0x8a, 0x9d, 0x27, 0x9e, 0xe9, 0x77, 0xd6, 0xc3, 0x33, 0xb1, 0x6a, 0x80, 0x2f, 0xd5, 0xf7, 0x4f,
	0xd8, 0x8e, 0xc0, 0xf3, 0x7b, 0x8f, 0x74, 0x2c, 0xcf, 0xf2, 0xfb, 0xb3, 0x79, 0x6b, 0x89, 0x8f,
	0x97, 0xdb, 0x86, 0xe1, 0xbb, 0xbc, 0xa6, 0xd8, 0xb3, 0xe2, 0x56, 0x94, 0xf6, 0x12, 0x80, 0x8c,
	0xc7, 0x5b, 0x82, 0xc3, 0x92, 0x61, 0xa7, 0xa3, 0xe3, 0x27, 0xad, 0xf1, 0x9f, 0x35, 0xba, 0x61,
	0xb8, 0x49, 0xeb, 0x65, 0x67, 0xc1, 0x5e, 0x80, 0x81, 0xe2, 0x0a, 0x65, 0x61, 0x2d, 0x39, 0x4f,
	0xf5, 0xef, 0x7a, 0x09, 0xeb, 0x29, 0x60, 0x35, 0x03, 0x6c, 0x66, 0x80, 0x1f, 0x38, 0x65, 0x4d,
	0x42, 0x5f, 0x9b, 0xea, 0x5c, 0xfb, 0x15, 0xe8, 0xfd, 0x4c, 0xa9, 0x22, 0x19, 0x95, 0xca, 0xe9,
	0x7a, 0x96, 0xdf, 0x5b, 0x5f, 0x85, 0xc5, 0xd7, 0xfd, 0xd1, 0x35, 0x0f, 0x47, 0xd7, 0x7c, 0x38,
	0xba, 0xe6, 0x9f, 0x93, 0x6b, 0x1c, 0x4e, 0xae, 0xf1, 0xef, 0xe4, 0x1a, 0xdf, 0xdf, 0x26, 0x54,
	0xa5, 0x65, 0x04, 0x63, 0x9e, 0x07, 0x25, 0xa3, 0x3f, 0x68, 0x8c, 0x14, 0xe5, 0xec, 0x4d, 0x75,
	0xbf, 0x3e, 0x9a, 0x5f, 0xb7, 0xcf, 0x46, 0xfd, 0x2e, 0x88, 0x8c, 0xba, 0x7a, 0xf5, 0xf9, 0x63,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x2a, 0x5a, 0xbd, 0xb9, 0x9b, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Whitelist) > 0 {
		for iNdEx := len(m.Whitelist) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Whitelist[iNdEx])
			copy(dAtA[i:], m.Whitelist[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.Whitelist[iNdEx])))
			i--
			dAtA[i] = 0x32
		}
	}
	{
		size, err := m.TotalLocked.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.LockedUnd) > 0 {
		for iNdEx := len(m.LockedUnd) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LockedUnd[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.PurchaseOrders) > 0 {
		for iNdEx := len(m.PurchaseOrders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PurchaseOrders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.StartingPurchaseOrderId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.StartingPurchaseOrderId))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.StartingPurchaseOrderId != 0 {
		n += 1 + sovGenesis(uint64(m.StartingPurchaseOrderId))
	}
	if len(m.PurchaseOrders) > 0 {
		for _, e := range m.PurchaseOrders {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.LockedUnd) > 0 {
		for _, e := range m.LockedUnd {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.TotalLocked.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Whitelist) > 0 {
		for _, s := range m.Whitelist {
			l = len(s)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartingPurchaseOrderId", wireType)
			}
			m.StartingPurchaseOrderId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartingPurchaseOrderId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PurchaseOrders", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PurchaseOrders = append(m.PurchaseOrders, EnterpriseUndPurchaseOrder{})
			if err := m.PurchaseOrders[len(m.PurchaseOrders)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockedUnd", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LockedUnd = append(m.LockedUnd, LockedUnd{})
			if err := m.LockedUnd[len(m.LockedUnd)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalLocked", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalLocked.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Whitelist", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Whitelist = append(m.Whitelist, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
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
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)