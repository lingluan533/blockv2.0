package domain

import (
	"blockv2.0/config"
	dataStruct "blockv2.0/datastruct"
	"database/sql"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"time"
)
func AddBlockHeaderToTdengine(date string,header dataStruct.BlockHeader){
	start := time.Now().UnixMilli()
	//var db *sql.DB
	//url := "root:taosdata@/tcp(" + "localhost" + ":" + "6030" + ")/"
	//db, err := sql.Open("taosSql", url)
	db := config.NewDBTDengine()
	header.Date = date
	str := "insert into "+globalconfig.Block.TDengineConfig.DBName+".blockheaders values "+"('"+header.CreateTimestamp +"','"+header.KeyId+"','"+strconv.FormatInt(header.BlockHeight,10)+"','"+header.DataType+"','"+header.DataValue+
		"','"+header.UpdateTimestamp+"','"+header.DataHash+"','"+header.BlockHash+"','"+header.PreBlockHash+"','"+strconv.FormatInt(int64(header.Nonce),10)+"','"+strconv.FormatInt(int64(header.Target),10)+"','"+strconv.FormatInt(header.CurrentDataCount,10)+
		"','"+strconv.FormatInt(header.CurrentDataSize,10)+"','"+header.Version+"','"+header.BlockType+"','"+header.LedgerType+"','"+header.Date+"')"

	res , err := db.DB.Exec(str)

	if err!= nil{
		logger.Warn("插入数据错误：",err)
	}else{
		fmt.Println("插入区块头到tdengine成功")
	}
	numOfaffected , err := res.RowsAffected()
	if err!= nil{
		logger.Warn("获取插入数目错误：",err)

	}else {
		fmt.Println("插入成功数量：",numOfaffected)
	}
	end := time.Now().UnixMilli()
	logger.Errorf("time need:=%d毫秒",end - start)
}
func AddTransactionBlockToTdengine(minuteTransactionBlock dataStruct.MinuteTransactionBlock)(int64,error) {
	start := time.Now().UnixMilli()
	logger.Error(minuteTransactionBlock.Transactions[0])
	str := "insert into "+globalconfig.Block.TDengineConfig.DBName+".transactions values "
	logger.Error("本次要插入数据条数：",len(minuteTransactionBlock.Transactions))
	//logger.Error("sql语句：",str)
	for _,data:=range minuteTransactionBlock.Transactions{
		data.BlockID = minuteTransactionBlock.Header.KeyId
		values := "('"+data.CreateTimestamp +"','"+data.EntityId+"','"+data.TransactionId+"','"+data.Initiator+"','"+data.Recipient+
			"','"+strconv.FormatFloat(data.TxAmount,'E',-1,64)+"','"+data.DataType+"','"+data.ServiceType+"','"+data.Remark+"','"+data.BlockID+"')"
		str = fmt.Sprintf("%s %s",str,values)
	}
	logger.Error("sql语句：",str)
	//res , err := db.Exec(str)
	db := config.NewDBTDengine()
	res , err := db.DB.Exec(str)

	if err!= nil{
		logger.Warn("插入数据错误：",err)
	}else{
		fmt.Println("插入成功")
	}
	numOfaffected , err := res.RowsAffected()
	if err!= nil{
		logger.Warn("获取插入数目错误：",err)
		return -1,err
	}else {
		fmt.Println("插入成功数量：",numOfaffected)
	}
	end := time.Now().UnixMilli()
	logger.Errorf("time need:=%d毫秒",end - start)
	return numOfaffected, nil
}
func AddDataBlockToTdengine(minuteDataBlock dataStruct.MinuteDataBlock)(int64,error){
	start := time.Now().UnixMilli()
	logger.Error(minuteDataBlock.DataReceipts[0])
	globalconfig := config.Initialize()

	str := "insert into "+globalconfig.Block.TDengineConfig.DBName+".datareceipts values "
	logger.Error("本次要插入数据条数：",len(minuteDataBlock.DataReceipts))
	//logger.Error("sql语句：",str)
	for _,data:=range minuteDataBlock.DataReceipts{
		data.BlockID = minuteDataBlock.Header.KeyId
		values := "('"+data.CreateTimeStamp +"','"+data.EntityId+"','"+data.KeyId+"','"+strconv.FormatFloat(data.ReceiptValue,'E',-1,64)+"','"+data.Version+
			"','"+data.UserName+"','"+data.OperationType+"','"+data.DataType+"','"+data.ServiceType+"','"+data.FileName+"','"+strconv.FormatFloat(data.FileSize,'E',-1,64)+"','"+data.FileHash+
			"','"+data.Uri+"','"+data.ParentKeyId+"',"+"'a'"+",'"+data.AttachmentTotalHash+"','"+data.BlockID+"')"
		str = fmt.Sprintf("%s %s",str,values)
	}
	logger.Error("sql语句：",str)
	//res , err := db.Exec(str)
	var db *sql.DB
	url := "root:taosdata@/tcp(" + "localhost" + ":" + "6030" + ")/"
	db, err := sql.Open("taosSql", url)
	res , err := db.Exec(str)
	defer db.Close()
	if err!= nil{
		logger.Warn("插入数据错误：",err)
	}else{
		fmt.Println("插入成功")
	}
	numOfaffected , err := res.RowsAffected()
	if err!= nil{
		logger.Warn("获取插入数目错误：",err)
		return -1,err
	}else {
		fmt.Println("插入成功数量：",numOfaffected)
	}
	end := time.Now().UnixMilli()
	logger.Errorf("time need:=%d毫秒",end - start)

	return numOfaffected, nil
}
