package v4

type V3WrkChain struct {
	// wrkchain_id is the id of the wrkchain
	WrkchainId uint64 `protobuf:"varint,1,opt,name=wrkchain_id,json=wrkchainId,proto3" json:"wrkchain_id,omitempty"`
	// moniker is the readable id of the wrkchain
	Moniker string `protobuf:"bytes,2,opt,name=moniker,proto3" json:"moniker,omitempty"`
	// name is the human friendly name of the wrkchain
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// genesis is an optional hash of the wrkchain's genesis block
	Genesis string `protobuf:"bytes,4,opt,name=genesis,proto3" json:"genesis,omitempty"`
	// type is the wrkchain type, e.g. geth, cosmos etc.
	Type string `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	// lastblock is the current highest recorded height for the wrkchain
	Lastblock uint64 `protobuf:"varint,6,opt,name=lastblock,proto3" json:"lastblock,omitempty"`
	// num_blocks is the current number of block hashes stored in state for the wrkchain
	NumBlocks uint64 `protobuf:"varint,7,opt,name=num_blocks,json=numBlocks,proto3" json:"num_blocks,omitempty"`
	// lowest_height is the lowest recorded height currently held in state for the wrkchain
	LowestHeight uint64 `protobuf:"varint,8,opt,name=lowest_height,json=lowestHeight,proto3" json:"lowest_height,omitempty"`
	// reg_time is the unix epoch of the wrkchain's registration time
	RegTime uint64 `protobuf:"varint,9,opt,name=reg_time,json=regTime,proto3" json:"reg_time,omitempty"`
	// owner is the owner address of the wrkchain
	Owner string `protobuf:"bytes,10,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (v V3WrkChain) Reset() {
}

func (v V3WrkChain) String() string {
	return ""
}

func (v V3WrkChain) ProtoMessage() {
}
