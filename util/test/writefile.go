package main

import (
	"blockv2.0/config"
	dataStruct "blockv2.0/datastruct"
	"blockv2.0/rpc"
	"bufio"
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"os"
)

type Person struct {
	Id   string
	Name string
	Age  int
	Sex  int
}
type Transactions struct {
	 transactions []rpc.Transaction
}
func main()  {
	//p1:= Person{
	//	Id:   "11",
	//	Name: "Black",
	//	Age:  10,
	//	Sex:  1,
	//}
	//p2:= Person{
	//	Id:   "12",
	//	Name: "Green",
	//	Age:  15,
	//	Sex:  2,
	//}
	//
	//bb1,_:=json.Marshal(p1)
	//bb2,_:=json.Marshal(p2)
	//json1:=string(bb1)
	//json2:=string(bb2)
	//fmt.Println(json1)
	//fmt.Println(json2)
	//var person1,person2 Person
	//json.Unmarshal([]byte(json1),&person1)
	//json.Unmarshal([]byte(json2),&person2)
	//fmt.Println(person1)
	//fmt.Println(person2)

	//receipts := "{CreateTimestamp:\"2021-10-25 22:51:01.010001\" EntityId:\"设备id\" KeyId:\"lwq|test.mp4|v1\" ReceiptValue:1 Version:\"v1\" UserName:\"lwq\" OperationType:\"create\" DataType:\"video\" ServiceType:\"ADD\" FileName:\"test.mp4\"ileSize:100 FileHash:\"asdfghjkl\" Uri:\"/abc\" ParentKeyId:\"000\" AttachmentTotalHash:\"zxcvbnm\"} "
	//var receipts string
	//receipts := "[{\"CreateTimestamp\":\"2021-06-16 20:39:42.543001\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|1\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543002\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|10\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543003\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|2\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543004\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|3\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543005\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|4\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543006\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|5\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543007\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|6\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543008\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|7\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543009\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|8\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"},{\"CreateTimestamp\":\"2021-06-16 20:39:42.543010\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|9\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"}]"
	//receipts := "{\"CreateTimestamp\":\"2021-06-16 20:39:42.543001\",\"EntityId\":\"A01\",\"TransactionId\":\"A01|2021-06-16|1\",\"Initiator\":\"A01\",\"TxAmount\":500,\"DataType\":\"node_credible\",\"ServiceType\":\"3\",\"Remark\":\"fee\",\"BlockIdentify\":\"node_credible_MINUTE_1239\",\"date\":\"2021-06-16\"}"


	//1.写文件
	blockPath :=  "/root/"+config.Initialize().ChainName+"/2021-10-25/video/MINUTE/1200"
	//logger.Info("目标区块文件目录为：",blockPath)
	//file ,err := os.OpenFile(blockPath,os.O_CREATE | os.O_WRONLY ,0777)
	//if err != nil{
	//	fmt.Println("err=",err)
	//}
	//
	//file.Write([]byte(receipts))

	//2.读文件
	var file *os.File
	var rawBlockdata []byte
	rawBlockdata = make([]byte,10240)
	var dataReceipts []dataStruct.DataReceipt
	file, err := os.OpenFile(blockPath,os.O_RDONLY,0777)
	defer file.Close()
	if err != nil{
		logger.Error("打开文件失败,err=",err)
		return
	}
	reader := bufio.NewReader(file)
	n,err := reader.Read(rawBlockdata)
	if err!=nil{
		panic(err)
	}
	fmt.Println("n:=",n)
	fmt.Println(string(rawBlockdata))
	err = json.Unmarshal(rawBlockdata[:n], &dataReceipts)
	if err != nil{
		logger.Error("从区块文件反序列化失败：",err)
		return
	}
	logger.Info("从区块文件反序列化成功：",dataReceipts)
	//插入tdengine数据库


	//fmt.Println(receipts)
	//var transactions []dataStruct.Transaction
	//err := json.Unmarshal([]byte(receipts),&transactions)
	//if err != nil{
	//	fmt.Println("err=",err)
	//}
	//fmt.Println(transactions)
	//receipts := "{\"CreateTimestamp\":\"2021-10-25 22:59:22.124001\" EntityId:\"设备id\" KeyId:\"lwq|test.mp4|v1\" ReceiptValue:1 Version:\"v1\" UserName:\"lwq\" OperationType:\"create\" DataType:\"video\" ServiceType:\"ADD\" FileName:\"test.mp4\" FileSize:0 FileHash:\"asdfghjkl\" Uri:\"/abc\" ParentKeyId:\"000\" AttachmentTotalHash:\"zxcvbnm\"} DataReceipts:{CreateTimestamp:\"2021-10-25 22:59:23.124001\" EntityId:\"设备id\" KeyId:\"l|test.mp4|v1\" ReceiptValue:1 Version:\"v1\" UserName:\"lwq\" OperationType:\"create\" DataType:\"video\" ServiceType:\"ADD\" FileName:\"test.mp4\" FileSize:100 FileHash:\"asdfghjkl\" Uri:\"/abc\" ParentKeyId:\"000\" AttachmentTotalHash:\"zxcvbnm\"} "
	//json.Unmarshal([]data(receipts),[]rpc.DataReceipt{})
	//blockPath :=  "/root/"+config.Initialize().ChainName+"/2021-10-25/video/MINUTE/1200"
	//logger.Info("目标区块文件目录为：",blockPath)
	//
	//
	//
	//

}