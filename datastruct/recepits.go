package dataStruct

type MinuteDataBlock struct {
	Header       *BlockHeader   `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"`
	DataReceipts []*DataReceipt `protobuf:"bytes,2,rep,name=DataReceipts,proto3" json:"DataReceipts,omitempty"` //元数据
}
type DataReceipt struct{
	CreateTimeStamp string `json:"CreateTimestamp" validate:"required"`
	EntityId  string `json:"entityId"`
	KeyId               string   `json:"keyId" validate:"required"`
	ReceiptValue        float64  `json:"receiptValue"`
	Version             string   `json:"version"`
	UserName            string   `json:"userName"`
	OperationType       string   `json:"operationType"`
	DataType            string   `json:"dataType" validate:"required"`
	ServiceType         string   `json:"serviceType"`
	FileName            string   `json:"fileName"`
	FileSize            float64  `json:"fileSize"`
	FileHash            string   `json:"fileHash"`
	Uri                 string   `json:"uri"`
	ParentKeyId         string   `json:"parentKeyId"`
	AttachmentFileUris  []string `json:"attachmentFileUris"`
	AttachmentTotalHash string   `json:"attachmentTotalHash"`
	BlockID 	  		string   `json:"blockId"`  //所属区块的唯一标识

}
type MockDataReceipt struct {
	CreateTimestamp     string   `protobuf:"bytes,1,opt,name=CreateTimestamp,proto3" json:"CreateTimestamp,omitempty"`
	EntityId            string   `protobuf:"bytes,2,opt,name=EntityId,proto3" json:"EntityId,omitempty"`
	KeyId               string   `protobuf:"bytes,3,opt,name=KeyId,proto3" json:"KeyId,omitempty"`
	ReceiptValue        float64  `protobuf:"fixed64,4,opt,name=ReceiptValue,proto3" json:"ReceiptValue,omitempty"`
	Version             string   `protobuf:"bytes,5,opt,name=Version,proto3" json:"Version,omitempty"`
	UserName            string   `protobuf:"bytes,6,opt,name=UserName,proto3" json:"UserName,omitempty"`
	OperationType       string   `protobuf:"bytes,7,opt,name=OperationType,proto3" json:"OperationType,omitempty"`
	DataType            string   `protobuf:"bytes,8,opt,name=DataType,proto3" json:"DataType,omitempty"`
	ServiceType         string   `protobuf:"bytes,9,opt,name=ServiceType,proto3" json:"ServiceType,omitempty"`
	FileName            string   `protobuf:"bytes,10,opt,name=FileName,proto3" json:"FileName,omitempty"`
	FileSize            float64  `protobuf:"fixed64,11,opt,name=FileSize,proto3" json:"FileSize,omitempty"`
	FileHash            string   `protobuf:"bytes,12,opt,name=FileHash,proto3" json:"FileHash,omitempty"`
	Uri                 string   `protobuf:"bytes,13,opt,name=Uri,proto3" json:"Uri,omitempty"`
	ParentKeyId         string   `protobuf:"bytes,14,opt,name=ParentKeyId,proto3" json:"ParentKeyId,omitempty"`
	AttachmentFileUris  []string `protobuf:"bytes,15,rep,name=AttachmentFileUris,proto3" json:"AttachmentFileUris,omitempty"`
	AttachmentTotalHash string   `protobuf:"bytes,16,opt,name=AttachmentTotalHash,proto3" json:"AttachmentTotalHash,omitempty"`
	BlockID 	  		string  `protobuf:"bytes,17,opt,name=BlockID,proto3" json:"BlockID,omitempty"`  //所属区块的唯一标识
}