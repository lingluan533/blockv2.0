package domain

import (
	"blockv2.0/config"
	dataStruct "blockv2.0/datastruct"
	"blockv2.0/util"
	"bufio"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"

	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)
var (
	dialTimout = 5*time.Second
	requestTimeout = 10*time.Second
	endpoints = []string{"127.0.0.1:2379"}
	globalconfig = config.Initialize()
)
const(
	VIDEO="video"
	USER="user_behaviour"
	NODE="node_credible"
	SENSOR="sensor"
	ACCESS="service_access"

	BLOCK_MINUTE="MINUTE"
	BLOCK_TENMINUTE="TENMINUTE"
	BLOCK_DAY="DAY"
)
func ScanAllBlockHeaders(client *clientv3.Client){
	indexMinInt := util.GetIndexMinInt()-2
	logger.Infof("检查1-%v是否已收到",indexMinInt)
	date := time.Now().Format("2006-01-02")
	str := "select count(*) from "+globalconfig.Block.TDengineConfig.DBName+".blockheaders where date = '"+date+"';"
	db :=config.NewDBTDengine()
	rows , err := db.DB.Query(str)
	if err != nil {
		logger.Error("err=",err)
	}
	var count int
	for rows.Next(){
		rows.Scan(&count)
	}
	if count == 5*indexMinInt {
		return
	}
	for i:=1;i<=indexMinInt;i++{
		logger.Infof("当前处理第%v个分钟块",i)
		date := time.Now().Format("2006-01-02")
		flag := date +":%:"+strconv.Itoa(i)
		str := "select count(*) from "+globalconfig.Block.TDengineConfig.DBName+".blockheaders where  keyid like '" + flag+"';"
		db :=config.NewDBTDengine()
		rows , err := db.DB.Query(str)
		if err != nil {
			logger.Error("err=",err)
		}
		var count int64
		for rows.Next(){
			rows.Scan(&count)
		}
		logger.Infof("第%v个分钟的已存入%v个区块头",i,count)
		if count == 5 {
			continue
		}else{
			//这个分钟的分钟块区块头少于五个时候，调用恢复分钟块的程序
			err := ReloadMinblock(client,i,date)
			if err !=nil {
				logger.Errorf("恢复第%v个分钟的区块时出错",i)
			}else{
				logger.Infof("恢复第%v个分钟的区块成功！",i)
			}
		}
	}

}
func ReloadMinblock(client *clientv3.Client,indexMinInt int,date string)(error) {
	var ledgers = []string{"user_behaviour", "node_credible", "video", "service_access", "sensor"}
	db := config.NewDBTDengine()
	for i := 0; i < 5; i++ {
		keyid := date + ":" + ledgers[i] + ":" + strconv.Itoa(indexMinInt)
		str := "select count(*) from " + globalconfig.Block.TDengineConfig.DBName + ".blockheaders where  keyid = '" + keyid + "';"
		rows, err := db.DB.Query(str)
		if err != nil {
			logger.Error("err=", err)
		}
		var count int64
		for rows.Next() {
			rows.Scan(&count)
		}
		if count == 1 {
			continue
		} else {
			//主动获取区块头，主动去找文件存
			etcdminblockkey := date + ":" + ledgers[i] + ":MINUTE:" + strconv.Itoa(indexMinInt)
			getblockresponse := util.GetData(client, etcdminblockkey, requestTimeout)
			switch ledgers[i] {
			case VIDEO, USER:
				var minutedatablock dataStruct.MinuteDataBlock
				for _, ev := range getblockresponse.Kvs {
					logger.Infof(" etcd获取的值：%v",string(ev.Value))
					err := json.Unmarshal(ev.Value, &minutedatablock)
					if err != nil {
						logger.Error("preMinuteMDData Unmarshal err", err)
						return err
					}
					//1.先存入区块头
					AddBlockHeaderToTdengine(date, *minutedatablock.Header)
					if minutedatablock.Header.CurrentDataCount == 0 {
						break
					} else if minutedatablock.Header.CurrentDataCount > 0 {
						//2.写区块体
						//2.1.找对应的区块文件
						//2.1 拼接出目标区块文件名字
						//2.3 读取文件转化为结构体
						//2.4 写入数据库
						const TIME_LAYOUT = "2006-01-02 15:04:05"
						time, err := time.Parse(TIME_LAYOUT, minutedatablock.Header.CreateTimestamp)
						if err != nil {
							return err
						}
						blockPath := "/root/go/hraft/scope" + "/" + strconv.FormatInt(int64(time.Year()), 10) + "-" + fmt.Sprintf("%02d", int64(time.Month())) + "-" + fmt.Sprintf("%02d", int64(time.Day())) + "/" + ledgers[i] + "/" +
							"MINUTE/" + strconv.Itoa(indexMinInt)

						logger.Info("上上分钟的目标存证区块文件目录为：", blockPath)

						//1.2打开区块文件
						var file *os.File
						var rawBlockdata []byte
						rawBlockdata = make([]byte, 99999999)
						var minuteBlockFromfile dataStruct.MinuteDataBlock
						file, err = os.OpenFile(blockPath, os.O_RDONLY, 0777)
						defer file.Close()
						if err != nil {
							logger.Error("打开文件失败,err=", err)
							return err
						}
						//2.读入区块文件到结构体数组
						reader := bufio.NewReader(file)
						n, err := reader.Read(rawBlockdata)
						if err != nil {
							return err
						}
						fmt.Println("n:=", n)
						err = json.Unmarshal(rawBlockdata[:n], &minuteBlockFromfile)
						if err != nil {
							logger.Error("从区块文件反序列化失败：", err)
							return err
						}
						logger.Info("从区块文件反序列化成功：", minuteBlockFromfile)
						//插入tdengine数据库

						AddDataBlockToTdengine(minuteBlockFromfile)
					}

				}

			case NODE, SENSOR, ACCESS:
				var minutetxblock dataStruct.MinuteTransactionBlock
				for _, ev := range getblockresponse.Kvs {
					err := json.Unmarshal(ev.Value, &minutetxblock)
					if err != nil {
						logger.Error("preMinuteMDData Unmarshal err", err)
					}
					//1.先存入区块头
					AddBlockHeaderToTdengine(date, *minutetxblock.Header)
					if minutetxblock.Header.CurrentDataCount == 0 {
						break
					} else if minutetxblock.Header.CurrentDataCount > 0 {
						//2.写区块体
						//2.1.找对应的区块文件
						//2.1 拼接出目标区块文件名字
						//2.3 读取文件转化为结构体
						//2.4 写入数据库
						const TIME_LAYOUT = "2006-01-02 15:04:05"
						time, err := time.Parse(TIME_LAYOUT, minutetxblock.Header.CreateTimestamp)
						if err != nil {
							return err
						}
						blockPath := "/root/go/hraft/scope" + "/" + strconv.FormatInt(int64(time.Year()), 10) + "-" + fmt.Sprintf("%02d", int64(time.Month())) + "-" + fmt.Sprintf("%02d", int64(time.Day())) + "/" + ledgers[i] + "/" +
							"MINUTE/" + strconv.Itoa(indexMinInt)

						logger.Info("上上分钟的目标交易区块文件目录为：", blockPath)

						//1.2打开区块文件
						var file *os.File
						var rawBlockdata []byte
						rawBlockdata = make([]byte, 99999999)
						var minuteBlockFromfile dataStruct.MinuteTransactionBlock
						file, err = os.OpenFile(blockPath, os.O_RDONLY, 0777)
						defer file.Close()
						if err != nil {
							logger.Error("打开文件失败,err=", err)
							return err
						}
						//2.读入区块文件到结构体数组
						reader := bufio.NewReader(file)
						n, err := reader.Read(rawBlockdata)
						if err != nil {
							return err
						}
						fmt.Println("n:=", n)
						err = json.Unmarshal(rawBlockdata[:n], &minuteBlockFromfile)
						if err != nil {
							logger.Error("从区块文件反序列化失败：", err)
							return err
						}
						logger.Info("从区块文件反序列化成功：", minuteBlockFromfile)
						//插入tdengine数据库

						AddTransactionBlockToTdengine(minuteBlockFromfile)
					}
				}

			}

		}

	}
return nil
}
func ScanLastBLockHeaders(client *clientv3.Client){
	//在这一分钟的0s检查上上分钟的区块头是否已经被接收到了
	indexMinInt := util.GetIndexMinInt()-2
	logger.Infof("检查%v是否已收到",indexMinInt)
	date := time.Now().Format("2006-01-02")
	flag := date +":%:"+strconv.Itoa(indexMinInt)
	str := "select count(*) from "+globalconfig.Block.TDengineConfig.DBName+".blockheaders where date = '"+date+"' and keyid like '" + flag+"';"
	db :=config.NewDBTDengine()
	rows , err := db.DB.Query(str)
	if err != nil {
		logger.Error("err=",err)
	}
	var count int64
	for rows.Next(){
		rows.Scan(&count)
	}
	if count == 5 {
		return
	}else{//这个分钟的分钟块区块头少于五个时候，调用恢复分钟块的程序
		err := ReloadMinblock(client,indexMinInt,date)
		if err !=nil {
			logger.Errorf("恢复第%v个分钟，类型为%v的区块时出错")
		}else{
			logger.Infof("恢复第%v个分钟，类型为%v的区块成功！")
		}}
}
func ScanReceivedBlockHeaders(client *clientv3.Client){
	 logger.Info("ScanReceivedBlockHeaders定时任务执行中。。。")
	start := time.Now().UnixMilli()
		//定时任务 扫描收到的所有的区块头 对比区块头中记录的记录数和表中的实际的存证记录数量是否一致，如果不一致则去找文件重新载入区块
		//var db *sql.DB
		//url := "root:taosdata@/tcp(" + "localhost" + ":" + "6030" + ")/"
		//db, err := sql.Open("taosSql", url)
		//if err != nil {
		//	fmt.Println("打开数据库错误。")
		//}
		db :=config.NewDBTDengine()

		str := "select * from "+globalconfig.Block.TDengineConfig.DBName+".blockheaders where date = '"+time.Now().Format("2006-01-02")+"';"
		rows , err := db.DB.Query(str)

		var headers []*dataStruct.BlockHeader
	if err != nil {
		logger.Error("err=",err)
	}

		for rows.Next(){
			//logger.Error("here")
			header := &dataStruct.BlockHeader{}
			rows.Scan(&header.CreateTimestamp,&header.KeyId,&header.BlockHeight,&header.DataType,&header.DataValue,
				&header.UpdateTimestamp,&header.DataHash,&header.BlockHash,&header.PreBlockHash,&header.Nonce,
				&header.Target,&header.CurrentDataCount,&header.CurrentDataSize,&header.Version,&header.BlockType,&header.LedgerType,&header.Date)
			//logger.Info(header)
			headers = append(headers,header)
		}
		for i := 0; i < len(headers); i++ {
			//如果区块头中记载的区块记录数为0 则认为该块中暂无数据
			if headers[i].CurrentDataCount==0{
				continue
			}
			switch headers[i].LedgerType {
			case VIDEO,USER:
				str := "select * from "+globalconfig.Block.TDengineConfig.DBName+".datareceipts where blockid = '"+headers[i].KeyId+"'"
				rows , err := db.DB.Query(str)
				if err != nil {
					logger.Error("查询出错,sql语句为：",str)
				}
				var count int64
				for rows.Next(){
					count ++
				}
				if count == headers[i].CurrentDataCount  {
					continue
				}
				//去找区块文件重新装载
				//KeyId:2021-11-02:service_access:1132
				keyIds := strings.Split(headers[i].KeyId,":")

				blockPath :=  "/root/go/hraft/scope"+"/"+keyIds[0]+"/"+keyIds[1]+"/MINUTE/"+keyIds[2]
				logger.Info("自动任务：目标存证区块文件目录为：",blockPath)
				//1.2打开区块文件
				var file *os.File
				var rawBlockdata []byte
				rawBlockdata = make([]byte,99999999)
				var minuteBlockFromfile dataStruct.MinuteDataBlock
				file, err = os.OpenFile(blockPath,os.O_RDONLY,0777)
				defer file.Close()
				if err != nil{
					logger.Error("打开文件失败,err=",err)
					return
				}
				//2.读入区块文件到结构体数组
				reader := bufio.NewReader(file)
				n,err := reader.Read(rawBlockdata)
				if err!=nil{
					panic(err)
				}
				fmt.Println("n:=",n)
				err = json.Unmarshal(rawBlockdata[:n], &minuteBlockFromfile)
				if err != nil{
					logger.Error("从区块文件反序列化失败：",err)
					return
				}
				logger.Info("从区块文件反序列化成功：",minuteBlockFromfile)
				//插入tdengine数据库

				AddDataBlockToTdengine(minuteBlockFromfile)
			case NODE,SENSOR,ACCESS:
				str := "select count(*) from "+globalconfig.Block.TDengineConfig.DBName+".transactions where blockid = '"+headers[i].KeyId+"'"
				rows , err := db.DB.Query(str)
				if err != nil {
					logger.Error("查询出错,sql语句为：",str)
				}
				var count int64
				for rows.Next(){
					rows.Scan(&count)
				}
				if count == headers[i].CurrentDataCount  {
					continue
				}

				//去找区块文件重新装载
				//KeyId:2021-11-02:service_access:1132
				keyIds := strings.Split(headers[i].KeyId,":")

				blockPath :=  "/root/go/hraft/scope"+"/"+keyIds[0]+"/"+keyIds[1]+"/MINUTE/"+keyIds[2]
				logger.Info("自动任务：目标交易区块文件目录为：",blockPath)

				//1.2打开区块文件
				var file *os.File
				var rawBlockdata []byte
				rawBlockdata = make([]byte,99999999)
				var minuteBlockFromfile dataStruct.MinuteTransactionBlock
				file, err = os.OpenFile(blockPath,os.O_RDONLY,0777)
				defer file.Close()
				if err != nil{
					logger.Error("打开文件失败,err=",err)
					return
				}
				//2.读入区块文件到结构体数组
				reader := bufio.NewReader(file)
				n,err := reader.Read(rawBlockdata)
				if err!=nil{
					panic(err)
				}
				fmt.Println("n:=",n)
				err = json.Unmarshal(rawBlockdata[:n], &minuteBlockFromfile)
				if err != nil{
					logger.Error("从区块文件反序列化失败：",err)
					return
				}
				logger.Info("从区块文件反序列化成功：",minuteBlockFromfile)
				//插入tdengine数据库
				AddTransactionBlockToTdengine(minuteBlockFromfile)
			}
		}
	end := time.Now().UnixMilli()
	logger.Infof("ScanReceivedBlockHeaders本次自动扫描区块任务结束，用时：%dms", end-start)



}
