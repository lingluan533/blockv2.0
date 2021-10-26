package dataStruct

type Transaction struct {
	Timestamp     string  `protobuf:"bytes,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	EntityId      string  `protobuf:"bytes,2,opt,name=EntityId,proto3" json:"EntityId,omitempty"`
	TransactionId string  `protobuf:"bytes,3,opt,name=TransactionId,proto3" json:"TransactionId,omitempty"`
	Initiator     string  `protobuf:"bytes,4,opt,name=Initiator,proto3" json:"Initiator,omitempty"`
	Recipient     string  `protobuf:"bytes,5,opt,name=Recipient,proto3" json:"Recipient,omitempty"`
	TxAmount      float64 `protobuf:"fixed64,6,opt,name=TxAmount,proto3" json:"TxAmount,omitempty"`
	DataType      string  `protobuf:"bytes,7,opt,name=DataType,proto3" json:"DataType,omitempty"`
	ServiceType   string  `protobuf:"bytes,8,opt,name=ServiceType,proto3" json:"ServiceType,omitempty"`
	Remark        string  `protobuf:"bytes,9,opt,name=Remark,proto3" json:"Remark,omitempty"`
}
