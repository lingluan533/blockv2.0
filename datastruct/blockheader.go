package dataStruct

type BlockHeader struct {
	CreateTimestamp string `protobuf:"bytes,1,opt,name=CreateTimestamp,proto3" json:"CreateTimestamp,omitempty"` //创建时间戳
	KeyId           string `protobuf:"bytes,2,opt,name=keyId,proto3" json:"keyId,omitempty"`
	BlockHeight     int64  `protobuf:"varint,3,opt,name=BlockHeight,proto3" json:"BlockHeight,omitempty"` //通过该字段，获取当前区块 可以使用不同链
	//具体数据结构类型
	DataType         string `protobuf:"bytes,4,opt,name=DataType,proto3" json:"DataType,omitempty"`                   //数据类型
	DataValue        string `protobuf:"bytes,5,opt,name=DataValue,proto3" json:"DataValue,omitempty"`                 //数据价值
	UpdateTimestamp  string `protobuf:"bytes,6,opt,name=UpdateTimestamp,proto3" json:"UpdateTimestamp,omitempty"`     //更新时间戳
	DataHash         string `protobuf:"bytes,7,opt,name=DataHash,proto3" json:"DataHash,omitempty"`                   //数据哈希值
	BlockHash        string `protobuf:"bytes,8,opt,name=BlockHash,proto3" json:"BlockHash,omitempty"`                 //区块哈希值
	PreBlockHash     string `protobuf:"bytes,9,opt,name=PreBlockHash,proto3" json:"PreBlockHash,omitempty"`           //前一个区块hash值
	Nonce            int32  `protobuf:"varint,10,opt,name=Nonce,proto3" json:"Nonce,omitempty"`                       //nonce 值
	Target           int32  `protobuf:"varint,11,opt,name=Target,proto3" json:"Target,omitempty"`                     //目标值
	CurrentDataCount int64  `protobuf:"varint,12,opt,name=CurrentDataCount,proto3" json:"CurrentDataCount,omitempty"` //当前数据记录量
	CurrentDataSize  int64  `protobuf:"varint,13,opt,name=CurrentDataSize,proto3" json:"CurrentDataSize,omitempty"`   //当前数据大小
	Version          string `protobuf:"bytes,14,opt,name=Version,proto3" json:"Version,omitempty"`                    //版本号
	BlockType        string `protobuf:"bytes,15,opt,name=BlockType,proto3" json:"BlockType,omitempty"`                //区块类型
	LedgerType       string `protobuf:"bytes,16,opt,name=LedgerType,proto3" json:"LedgerType,omitempty"`              //账本类型
	Date             string `protobuf:"bytes,17,opt,name=date,proto3" json:"date,omitempty"`                          //创建时间
}