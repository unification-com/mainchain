// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mainchain/wrkchain/v1/wrkchain.proto

package v040

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

// WrkChain holds metadata about a registered wrkchain
type WrkChain struct {
	WrkchainId   uint64 `protobuf:"varint,1,opt,name=wrkchain_id,json=wrkchainId,proto3" json:"wrkchain_id,omitempty"`
	Moniker      string `protobuf:"bytes,2,opt,name=moniker,proto3" json:"moniker,omitempty"`
	Name         string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Genesis      string `protobuf:"bytes,4,opt,name=genesis,proto3" json:"genesis,omitempty"`
	Type         string `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	Lastblock    uint64 `protobuf:"varint,6,opt,name=lastblock,proto3" json:"lastblock,omitempty"`
	NumBlocks    uint64 `protobuf:"varint,7,opt,name=num_blocks,json=numBlocks,proto3" json:"num_blocks,omitempty"`
	LowestHeight uint64 `protobuf:"varint,8,opt,name=lowest_height,json=lowestHeight,proto3" json:"lowest_height,omitempty"`
	RegTime      uint64 `protobuf:"varint,9,opt,name=reg_time,json=regTime,proto3" json:"reg_time,omitempty"`
	Owner        string `protobuf:"bytes,10,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (m *WrkChain) Reset()         { *m = WrkChain{} }
func (m *WrkChain) String() string { return proto.CompactTextString(m) }
func (*WrkChain) ProtoMessage()    {}
func (*WrkChain) Descriptor() ([]byte, []int) {
	return fileDescriptor_02a2970309ba545c, []int{0}
}
func (m *WrkChain) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WrkChain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WrkChain.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WrkChain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WrkChain.Merge(m, src)
}
func (m *WrkChain) XXX_Size() int {
	return m.Size()
}
func (m *WrkChain) XXX_DiscardUnknown() {
	xxx_messageInfo_WrkChain.DiscardUnknown(m)
}

var xxx_messageInfo_WrkChain proto.InternalMessageInfo

func (m *WrkChain) GetWrkchainId() uint64 {
	if m != nil {
		return m.WrkchainId
	}
	return 0
}

func (m *WrkChain) GetMoniker() string {
	if m != nil {
		return m.Moniker
	}
	return ""
}

func (m *WrkChain) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *WrkChain) GetGenesis() string {
	if m != nil {
		return m.Genesis
	}
	return ""
}

func (m *WrkChain) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *WrkChain) GetLastblock() uint64 {
	if m != nil {
		return m.Lastblock
	}
	return 0
}

func (m *WrkChain) GetNumBlocks() uint64 {
	if m != nil {
		return m.NumBlocks
	}
	return 0
}

func (m *WrkChain) GetLowestHeight() uint64 {
	if m != nil {
		return m.LowestHeight
	}
	return 0
}

func (m *WrkChain) GetRegTime() uint64 {
	if m != nil {
		return m.RegTime
	}
	return 0
}

func (m *WrkChain) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

// WrkChainBlock holds data about a wrkchain's block hash submission
type WrkChainBlock struct {
	Height     uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Blockhash  string `protobuf:"bytes,2,opt,name=blockhash,proto3" json:"blockhash,omitempty"`
	Parenthash string `protobuf:"bytes,3,opt,name=parenthash,proto3" json:"parenthash,omitempty"`
	Hash1      string `protobuf:"bytes,4,opt,name=hash1,proto3" json:"hash1,omitempty"`
	Hash2      string `protobuf:"bytes,5,opt,name=hash2,proto3" json:"hash2,omitempty"`
	Hash3      string `protobuf:"bytes,6,opt,name=hash3,proto3" json:"hash3,omitempty"`
	SubTime    uint64 `protobuf:"varint,7,opt,name=sub_time,json=subTime,proto3" json:"sub_time,omitempty"`
}

func (m *WrkChainBlock) Reset()         { *m = WrkChainBlock{} }
func (m *WrkChainBlock) String() string { return proto.CompactTextString(m) }
func (*WrkChainBlock) ProtoMessage()    {}
func (*WrkChainBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_02a2970309ba545c, []int{1}
}
func (m *WrkChainBlock) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WrkChainBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WrkChainBlock.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WrkChainBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WrkChainBlock.Merge(m, src)
}
func (m *WrkChainBlock) XXX_Size() int {
	return m.Size()
}
func (m *WrkChainBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_WrkChainBlock.DiscardUnknown(m)
}

var xxx_messageInfo_WrkChainBlock proto.InternalMessageInfo

func (m *WrkChainBlock) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *WrkChainBlock) GetBlockhash() string {
	if m != nil {
		return m.Blockhash
	}
	return ""
}

func (m *WrkChainBlock) GetParenthash() string {
	if m != nil {
		return m.Parenthash
	}
	return ""
}

func (m *WrkChainBlock) GetHash1() string {
	if m != nil {
		return m.Hash1
	}
	return ""
}

func (m *WrkChainBlock) GetHash2() string {
	if m != nil {
		return m.Hash2
	}
	return ""
}

func (m *WrkChainBlock) GetHash3() string {
	if m != nil {
		return m.Hash3
	}
	return ""
}

func (m *WrkChainBlock) GetSubTime() uint64 {
	if m != nil {
		return m.SubTime
	}
	return 0
}

// Params defines the parameters for the wrkchain module.
type Params struct {
	FeeRegister uint64 `protobuf:"varint,1,opt,name=fee_register,json=feeRegister,proto3" json:"fee_register,omitempty"`
	FeeRecord   uint64 `protobuf:"varint,2,opt,name=fee_record,json=feeRecord,proto3" json:"fee_record,omitempty"`
	Denom       string `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_02a2970309ba545c, []int{2}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetFeeRegister() uint64 {
	if m != nil {
		return m.FeeRegister
	}
	return 0
}

func (m *Params) GetFeeRecord() uint64 {
	if m != nil {
		return m.FeeRecord
	}
	return 0
}

func (m *Params) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func init() {
	proto.RegisterType((*WrkChain)(nil), "mainchain.wrkchain.v1.WrkChain")
	proto.RegisterType((*WrkChainBlock)(nil), "mainchain.wrkchain.v1.WrkChainBlock")
	proto.RegisterType((*Params)(nil), "mainchain.wrkchain.v1.Params")
}

func init() {
	proto.RegisterFile("mainchain/wrkchain/v1/wrkchain.proto", fileDescriptor_02a2970309ba545c)
}

var fileDescriptor_02a2970309ba545c = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x92, 0xbf, 0x6e, 0x13, 0x41,
	0x10, 0xc6, 0x7d, 0xc1, 0xf1, 0x9f, 0x49, 0xd2, 0xac, 0x02, 0x5a, 0x10, 0x1c, 0xc1, 0x50, 0xa4,
	0xc1, 0x8e, 0x31, 0x1d, 0x5d, 0x68, 0xa0, 0x43, 0x16, 0x08, 0x89, 0xe6, 0xd8, 0x3b, 0x8f, 0xef,
	0x56, 0xf6, 0xee, 0x5a, 0xbb, 0x7b, 0x36, 0x79, 0x0b, 0x6a, 0x1e, 0x87, 0x8a, 0xd2, 0x25, 0x25,
	0xb2, 0x5f, 0x04, 0xed, 0x9f, 0xb3, 0x53, 0xdd, 0x7c, 0xbf, 0xf9, 0x4e, 0xb3, 0xf3, 0x69, 0xe0,
	0x95, 0x60, 0x5c, 0x16, 0x15, 0xe3, 0x72, 0xb4, 0xd1, 0x8b, 0x50, 0xac, 0xc7, 0x87, 0x7a, 0xb8,
	0xd2, 0xca, 0x2a, 0xf2, 0xf0, 0xe0, 0x1a, 0x1e, 0x3a, 0xeb, 0xf1, 0x93, 0xcb, 0x52, 0x95, 0xca,
	0x3b, 0x46, 0xae, 0x0a, 0xe6, 0xc1, 0xaf, 0x13, 0xe8, 0x7d, 0xd5, 0x8b, 0xf7, 0xce, 0x45, 0x9e,
	0xc3, 0x59, 0xf3, 0x47, 0xc6, 0x67, 0x34, 0xb9, 0x4a, 0xae, 0xdb, 0x53, 0x68, 0xd0, 0xc7, 0x19,
	0xa1, 0xd0, 0x15, 0x4a, 0xf2, 0x05, 0x6a, 0x7a, 0x72, 0x95, 0x5c, 0xf7, 0xa7, 0x8d, 0x24, 0x04,
	0xda, 0x92, 0x09, 0xa4, 0x0f, 0x3c, 0xf6, 0xb5, 0x73, 0x97, 0x28, 0xd1, 0x70, 0x43, 0xdb, 0xc1,
	0x1d, 0xa5, 0x73, 0xdb, 0xbb, 0x15, 0xd2, 0xd3, 0xe0, 0x76, 0x35, 0x79, 0x0a, 0xfd, 0x25, 0x33,
	0x36, 0x5f, 0xaa, 0x62, 0x41, 0x3b, 0x7e, 0xf4, 0x11, 0x90, 0x67, 0x00, 0xb2, 0x16, 0x99, 0x17,
	0x86, 0x76, 0x43, 0x5b, 0xd6, 0xe2, 0xd6, 0x03, 0xf2, 0x12, 0x2e, 0x96, 0x6a, 0x83, 0xc6, 0x66,
	0x15, 0xf2, 0xb2, 0xb2, 0xb4, 0xe7, 0x1d, 0xe7, 0x01, 0x7e, 0xf0, 0x8c, 0x3c, 0x86, 0x9e, 0xc6,
	0x32, 0xb3, 0x5c, 0x20, 0xed, 0xfb, 0x7e, 0x57, 0x63, 0xf9, 0x99, 0x0b, 0x24, 0x97, 0x70, 0xaa,
	0x36, 0x12, 0x35, 0x05, 0xff, 0xa2, 0x20, 0x06, 0xbf, 0x13, 0xb8, 0x68, 0xc2, 0xf1, 0x83, 0xc8,
	0x23, 0xe8, 0xc4, 0x01, 0x21, 0x9c, 0xa8, 0xdc, 0xe3, 0xfd, 0xd3, 0x2a, 0x66, 0xaa, 0x18, 0xcd,
	0x11, 0x90, 0x14, 0x60, 0xc5, 0x34, 0x4a, 0xeb, 0xdb, 0x21, 0xa2, 0x7b, 0xc4, 0x4d, 0x77, 0xdf,
	0x71, 0x8c, 0x29, 0x88, 0x86, 0xbe, 0x89, 0x29, 0x05, 0xd1, 0xd0, 0x89, 0x8f, 0x28, 0xd2, 0x89,
	0x5b, 0xcd, 0xd4, 0x79, 0x58, 0x2d, 0x84, 0xd3, 0x35, 0x75, 0xee, 0x56, 0x1b, 0x7c, 0x87, 0xce,
	0x27, 0xa6, 0x99, 0x30, 0xe4, 0x05, 0x9c, 0xcf, 0x11, 0x33, 0x8d, 0x25, 0x37, 0x16, 0x75, 0x5c,
	0xe1, 0x6c, 0x8e, 0x38, 0x8d, 0xc8, 0xc5, 0x1c, 0x2c, 0x85, 0xd2, 0x33, 0xbf, 0x48, 0x7b, 0xda,
	0xf7, 0x06, 0x07, 0xdc, 0xf0, 0x19, 0x4a, 0x25, 0xe2, 0x0e, 0x41, 0xdc, 0x7e, 0xf9, 0xb3, 0x4b,
	0x93, 0xed, 0x2e, 0x4d, 0xfe, 0xed, 0xd2, 0xe4, 0xe7, 0x3e, 0x6d, 0x6d, 0xf7, 0x69, 0xeb, 0xef,
	0x3e, 0x6d, 0x7d, 0x7b, 0x57, 0x72, 0x5b, 0xd5, 0xf9, 0xb0, 0x50, 0x62, 0x54, 0x4b, 0x3e, 0xe7,
	0x05, 0xb3, 0x5c, 0xc9, 0xd7, 0x4e, 0x1f, 0x6f, 0xf9, 0xc7, 0xf1, 0x9a, 0x97, 0x58, 0xb2, 0xe2,
	0x6e, 0xb4, 0xbe, 0x79, 0x7b, 0x93, 0x77, 0xfc, 0x85, 0x4e, 0xfe, 0x07, 0x00, 0x00, 0xff, 0xff,
	0xdd, 0x19, 0xb1, 0x9f, 0xf6, 0x02, 0x00, 0x00,
}

func (m *WrkChain) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WrkChain) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WrkChain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x52
	}
	if m.RegTime != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.RegTime))
		i--
		dAtA[i] = 0x48
	}
	if m.LowestHeight != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.LowestHeight))
		i--
		dAtA[i] = 0x40
	}
	if m.NumBlocks != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.NumBlocks))
		i--
		dAtA[i] = 0x38
	}
	if m.Lastblock != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.Lastblock))
		i--
		dAtA[i] = 0x30
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Genesis) > 0 {
		i -= len(m.Genesis)
		copy(dAtA[i:], m.Genesis)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Genesis)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Moniker) > 0 {
		i -= len(m.Moniker)
		copy(dAtA[i:], m.Moniker)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Moniker)))
		i--
		dAtA[i] = 0x12
	}
	if m.WrkchainId != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.WrkchainId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *WrkChainBlock) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WrkChainBlock) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WrkChainBlock) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SubTime != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.SubTime))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Hash3) > 0 {
		i -= len(m.Hash3)
		copy(dAtA[i:], m.Hash3)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Hash3)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Hash2) > 0 {
		i -= len(m.Hash2)
		copy(dAtA[i:], m.Hash2)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Hash2)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Hash1) > 0 {
		i -= len(m.Hash1)
		copy(dAtA[i:], m.Hash1)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Hash1)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Parenthash) > 0 {
		i -= len(m.Parenthash)
		copy(dAtA[i:], m.Parenthash)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Parenthash)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Blockhash) > 0 {
		i -= len(m.Blockhash)
		copy(dAtA[i:], m.Blockhash)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Blockhash)))
		i--
		dAtA[i] = 0x12
	}
	if m.Height != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintWrkchain(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if m.FeeRecord != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.FeeRecord))
		i--
		dAtA[i] = 0x10
	}
	if m.FeeRegister != 0 {
		i = encodeVarintWrkchain(dAtA, i, uint64(m.FeeRegister))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintWrkchain(dAtA []byte, offset int, v uint64) int {
	offset -= sovWrkchain(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *WrkChain) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.WrkchainId != 0 {
		n += 1 + sovWrkchain(uint64(m.WrkchainId))
	}
	l = len(m.Moniker)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	l = len(m.Genesis)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	if m.Lastblock != 0 {
		n += 1 + sovWrkchain(uint64(m.Lastblock))
	}
	if m.NumBlocks != 0 {
		n += 1 + sovWrkchain(uint64(m.NumBlocks))
	}
	if m.LowestHeight != 0 {
		n += 1 + sovWrkchain(uint64(m.LowestHeight))
	}
	if m.RegTime != 0 {
		n += 1 + sovWrkchain(uint64(m.RegTime))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	return n
}

func (m *WrkChainBlock) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovWrkchain(uint64(m.Height))
	}
	l = len(m.Blockhash)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	l = len(m.Parenthash)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	l = len(m.Hash1)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	l = len(m.Hash2)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	l = len(m.Hash3)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	if m.SubTime != 0 {
		n += 1 + sovWrkchain(uint64(m.SubTime))
	}
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.FeeRegister != 0 {
		n += 1 + sovWrkchain(uint64(m.FeeRegister))
	}
	if m.FeeRecord != 0 {
		n += 1 + sovWrkchain(uint64(m.FeeRecord))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovWrkchain(uint64(l))
	}
	return n
}

func sovWrkchain(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozWrkchain(x uint64) (n int) {
	return sovWrkchain(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WrkChain) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWrkchain
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
			return fmt.Errorf("proto: WrkChain: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WrkChain: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WrkchainId", wireType)
			}
			m.WrkchainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.WrkchainId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Moniker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Moniker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Genesis", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Genesis = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Lastblock", wireType)
			}
			m.Lastblock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Lastblock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumBlocks", wireType)
			}
			m.NumBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumBlocks |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LowestHeight", wireType)
			}
			m.LowestHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LowestHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegTime", wireType)
			}
			m.RegTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RegTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipWrkchain(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthWrkchain
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
func (m *WrkChainBlock) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWrkchain
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
			return fmt.Errorf("proto: WrkChainBlock: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WrkChainBlock: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Blockhash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Blockhash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Parenthash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Parenthash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash1", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash1 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash2", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash2 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash3", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash3 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubTime", wireType)
			}
			m.SubTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SubTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipWrkchain(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthWrkchain
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
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWrkchain
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeRegister", wireType)
			}
			m.FeeRegister = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeeRegister |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeRecord", wireType)
			}
			m.FeeRecord = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeeRecord |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWrkchain
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
				return ErrInvalidLengthWrkchain
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthWrkchain
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipWrkchain(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthWrkchain
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
func skipWrkchain(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowWrkchain
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
					return 0, ErrIntOverflowWrkchain
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
					return 0, ErrIntOverflowWrkchain
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
				return 0, ErrInvalidLengthWrkchain
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupWrkchain
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthWrkchain
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthWrkchain        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowWrkchain          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupWrkchain = fmt.Errorf("proto: unexpected end of group")
)