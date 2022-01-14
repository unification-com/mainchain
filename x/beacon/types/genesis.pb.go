// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mainchain/beacon/v1/genesis.proto

package types

import (
	fmt "fmt"
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

// GenesisState defines the beacon module's genesis state.
type GenesisState struct {
	// params defines all the paramaters of the module.
	Params            Params        `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	StartingBeaconId  uint64        `protobuf:"varint,2,opt,name=starting_beacon_id,json=startingBeaconId,proto3" json:"starting_beacon_id,omitempty"`
	RegisteredBeacons BeaconExports `protobuf:"bytes,3,rep,name=registered_beacons,json=registeredBeacons,proto3,castrepeated=BeaconExports" json:"registered_beacons"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c746adb253d3070, []int{0}
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

func (m *GenesisState) GetStartingBeaconId() uint64 {
	if m != nil {
		return m.StartingBeaconId
	}
	return 0
}

func (m *GenesisState) GetRegisteredBeacons() BeaconExports {
	if m != nil {
		return m.RegisteredBeacons
	}
	return nil
}

type BeaconTimestampGenesisExport struct {
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	T  uint64 `protobuf:"varint,2,opt,name=t,proto3" json:"t,omitempty"`
	H  string `protobuf:"bytes,3,opt,name=h,proto3" json:"h,omitempty"`
}

func (m *BeaconTimestampGenesisExport) Reset()         { *m = BeaconTimestampGenesisExport{} }
func (m *BeaconTimestampGenesisExport) String() string { return proto.CompactTextString(m) }
func (*BeaconTimestampGenesisExport) ProtoMessage()    {}
func (*BeaconTimestampGenesisExport) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c746adb253d3070, []int{1}
}
func (m *BeaconTimestampGenesisExport) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BeaconTimestampGenesisExport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BeaconTimestampGenesisExport.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BeaconTimestampGenesisExport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BeaconTimestampGenesisExport.Merge(m, src)
}
func (m *BeaconTimestampGenesisExport) XXX_Size() int {
	return m.Size()
}
func (m *BeaconTimestampGenesisExport) XXX_DiscardUnknown() {
	xxx_messageInfo_BeaconTimestampGenesisExport.DiscardUnknown(m)
}

var xxx_messageInfo_BeaconTimestampGenesisExport proto.InternalMessageInfo

func (m *BeaconTimestampGenesisExport) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *BeaconTimestampGenesisExport) GetT() uint64 {
	if m != nil {
		return m.T
	}
	return 0
}

func (m *BeaconTimestampGenesisExport) GetH() string {
	if m != nil {
		return m.H
	}
	return ""
}

type BeaconExport struct {
	Beacon     Beacon                        `protobuf:"bytes,1,opt,name=beacon,proto3" json:"beacon"`
	Timestamps BeaconTimestampGenesisExports `protobuf:"bytes,2,rep,name=timestamps,proto3,castrepeated=BeaconTimestampGenesisExports" json:"timestamps"`
}

func (m *BeaconExport) Reset()         { *m = BeaconExport{} }
func (m *BeaconExport) String() string { return proto.CompactTextString(m) }
func (*BeaconExport) ProtoMessage()    {}
func (*BeaconExport) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c746adb253d3070, []int{2}
}
func (m *BeaconExport) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BeaconExport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BeaconExport.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BeaconExport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BeaconExport.Merge(m, src)
}
func (m *BeaconExport) XXX_Size() int {
	return m.Size()
}
func (m *BeaconExport) XXX_DiscardUnknown() {
	xxx_messageInfo_BeaconExport.DiscardUnknown(m)
}

var xxx_messageInfo_BeaconExport proto.InternalMessageInfo

func (m *BeaconExport) GetBeacon() Beacon {
	if m != nil {
		return m.Beacon
	}
	return Beacon{}
}

func (m *BeaconExport) GetTimestamps() BeaconTimestampGenesisExports {
	if m != nil {
		return m.Timestamps
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "mainchain.beacon.v1.GenesisState")
	proto.RegisterType((*BeaconTimestampGenesisExport)(nil), "mainchain.beacon.v1.BeaconTimestampGenesisExport")
	proto.RegisterType((*BeaconExport)(nil), "mainchain.beacon.v1.BeaconExport")
}

func init() { proto.RegisterFile("mainchain/beacon/v1/genesis.proto", fileDescriptor_9c746adb253d3070) }

var fileDescriptor_9c746adb253d3070 = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0xb3, 0x69, 0x29, 0xb8, 0xad, 0xa2, 0xab, 0x42, 0xa8, 0x9a, 0xa6, 0x05, 0x21, 0x07,
	0x4d, 0x68, 0x7b, 0xf2, 0x1a, 0x10, 0x51, 0x10, 0x24, 0x7a, 0xf2, 0x52, 0xb6, 0xc9, 0x9a, 0xec,
	0x21, 0xd9, 0x90, 0xdd, 0x96, 0xfa, 0x16, 0x3e, 0x87, 0xaf, 0xe1, 0xa5, 0xc7, 0x1e, 0x3d, 0x69,
	0x69, 0x5f, 0x44, 0x92, 0x4d, 0x6c, 0x0f, 0xa1, 0xb7, 0xdd, 0x99, 0xef, 0x9f, 0x7f, 0x66, 0x18,
	0xd8, 0x8d, 0x30, 0x8d, 0xbd, 0x10, 0xd3, 0xd8, 0x1e, 0x13, 0xec, 0xb1, 0xd8, 0x9e, 0xf6, 0xed,
	0x80, 0xc4, 0x84, 0x53, 0x6e, 0x25, 0x29, 0x13, 0x0c, 0x1d, 0xff, 0x23, 0x96, 0x44, 0xac, 0x69,
	0xbf, 0x7d, 0x12, 0xb0, 0x80, 0xe5, 0x79, 0x3b, 0x7b, 0x49, 0xb4, 0x6d, 0x54, 0x55, 0x2b, 0x44,
	0x39, 0xd1, 0x5b, 0x02, 0xd8, 0xba, 0x93, 0xe5, 0x9f, 0x05, 0x16, 0x04, 0xdd, 0xc0, 0x46, 0x82,
	0x53, 0x1c, 0x71, 0x0d, 0x18, 0xc0, 0x6c, 0x0e, 0xce, 0xac, 0x0a, 0x3b, 0xeb, 0x29, 0x47, 0x9c,
	0xfa, 0xfc, 0xa7, 0xa3, 0xb8, 0x85, 0x00, 0x5d, 0x41, 0xc4, 0x05, 0x4e, 0x05, 0x8d, 0x83, 0x91,
	0x44, 0x47, 0xd4, 0xd7, 0x54, 0x03, 0x98, 0x75, 0xf7, 0xb0, 0xcc, 0x38, 0x79, 0xe2, 0xde, 0x47,
	0x01, 0x44, 0x29, 0x09, 0x28, 0x17, 0x24, 0x25, 0x7e, 0xc1, 0x73, 0xad, 0x66, 0xd4, 0xcc, 0xe6,
	0xa0, 0x5b, 0x69, 0x2a, 0xa5, 0xb7, 0xb3, 0x84, 0xa5, 0xc2, 0x39, 0xcd, 0xac, 0x3f, 0x7f, 0x3b,
	0xfb, 0xdb, 0x51, 0xee, 0x1e, 0x6d, 0x6a, 0xca, 0x04, 0xef, 0x3d, 0xc0, 0x73, 0xf9, 0x7c, 0xa1,
	0x11, 0xe1, 0x02, 0x47, 0x49, 0x31, 0xb0, 0xd4, 0xa0, 0x03, 0xa8, 0x52, 0x3f, 0x9f, 0xb6, 0xee,
	0xaa, 0xd4, 0x47, 0x2d, 0x08, 0x44, 0xd1, 0x35, 0x10, 0xd9, 0x2f, 0xd4, 0x6a, 0x06, 0x30, 0xf7,
	0x5c, 0x10, 0xf6, 0xbe, 0x00, 0x6c, 0x6d, 0x1b, 0x66, 0xeb, 0x92, 0x0d, 0xee, 0x5c, 0x97, 0x94,
	0x94, 0xeb, 0x92, 0x71, 0x34, 0x83, 0x50, 0x94, 0x1d, 0x71, 0x4d, 0xcd, 0x07, 0xef, 0xef, 0x90,
	0x57, 0xb7, 0xef, 0x5c, 0x16, 0x8b, 0xb8, 0xd8, 0x45, 0x71, 0x77, 0xcb, 0xcb, 0x79, 0x9c, 0xaf,
	0x74, 0xb0, 0x58, 0xe9, 0x60, 0xb9, 0xd2, 0xc1, 0xc7, 0x5a, 0x57, 0x16, 0x6b, 0x5d, 0xf9, 0x5e,
	0xeb, 0xca, 0xeb, 0x30, 0xa0, 0x22, 0x9c, 0x8c, 0x2d, 0x8f, 0x45, 0xf6, 0x24, 0xa6, 0x6f, 0xd4,
	0xc3, 0x82, 0xb2, 0xf8, 0x3a, 0xfb, 0x6f, 0x6e, 0x69, 0x56, 0x5e, 0x93, 0x78, 0x4f, 0x08, 0x1f,
	0x37, 0xf2, 0x53, 0x1a, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x2d, 0x7e, 0x28, 0x9e, 0xbc, 0x02,
	0x00, 0x00,
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
	if len(m.RegisteredBeacons) > 0 {
		for iNdEx := len(m.RegisteredBeacons) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RegisteredBeacons[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if m.StartingBeaconId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.StartingBeaconId))
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

func (m *BeaconTimestampGenesisExport) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BeaconTimestampGenesisExport) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BeaconTimestampGenesisExport) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.H) > 0 {
		i -= len(m.H)
		copy(dAtA[i:], m.H)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.H)))
		i--
		dAtA[i] = 0x1a
	}
	if m.T != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.T))
		i--
		dAtA[i] = 0x10
	}
	if m.Id != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BeaconExport) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BeaconExport) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BeaconExport) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Timestamps) > 0 {
		for iNdEx := len(m.Timestamps) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Timestamps[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Beacon.MarshalToSizedBuffer(dAtA[:i])
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
	if m.StartingBeaconId != 0 {
		n += 1 + sovGenesis(uint64(m.StartingBeaconId))
	}
	if len(m.RegisteredBeacons) > 0 {
		for _, e := range m.RegisteredBeacons {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *BeaconTimestampGenesisExport) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovGenesis(uint64(m.Id))
	}
	if m.T != 0 {
		n += 1 + sovGenesis(uint64(m.T))
	}
	l = len(m.H)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func (m *BeaconExport) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Beacon.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Timestamps) > 0 {
		for _, e := range m.Timestamps {
			l = e.Size()
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
				return fmt.Errorf("proto: wrong wireType = %d for field StartingBeaconId", wireType)
			}
			m.StartingBeaconId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartingBeaconId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegisteredBeacons", wireType)
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
			m.RegisteredBeacons = append(m.RegisteredBeacons, BeaconExport{})
			if err := m.RegisteredBeacons[len(m.RegisteredBeacons)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *BeaconTimestampGenesisExport) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: BeaconTimestampGenesisExport: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BeaconTimestampGenesisExport: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field T", wireType)
			}
			m.T = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.T |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field H", wireType)
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
			m.H = string(dAtA[iNdEx:postIndex])
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
func (m *BeaconExport) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: BeaconExport: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BeaconExport: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Beacon", wireType)
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
			if err := m.Beacon.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamps", wireType)
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
			m.Timestamps = append(m.Timestamps, BeaconTimestampGenesisExport{})
			if err := m.Timestamps[len(m.Timestamps)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
