package main
//
import (
	"blockv2.0/config"
	dataStruct "blockv2.0/datastruct"
	"blockv2.0/rpc"
	"blockv2.0/util"
	"bufio"
	"database/sql"
	"encoding/json"
	"os"
	"strconv"

	//"block/dal/datasource"
	//"block/dataStruct"
	//"blockv2.0/rpc"
	//"block/service"
	//	sync "block/sync"
	"context"
	//"encoding/json"
	"flag"
	"fmt"
	logger "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	//"net"
	//os "os"
	"strings"
	"time"
)

var (
	dialTimout = 5*time.Second
	requestTimeout = 10*time.Second
	endpoints = []string{"127.0.0.1:2379"}
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
//func rpcService(port string){
//	//启动RPC服务，将启动
//	ls, err := net.Listen("tcp", ":"+port)
//	if err!=nil{
//		panic(err)
//	}
//	//
//	//创建grpc服务
//	gs := grpc.NewServer()
//	//注册服务
//	//五个类型的账本 都实现了blockService，在查询的时候注册blockService即可
//	//todo：配置文件中没有配置的数据不启动
//	//rpc.RegisterAccessLedgerServiceServer(gs,&service.DefaultAccessLedgerService)
//	//rpc.RegisterNodeLedgerServiceServer(gs,&service.DefaultNodeLedgerService)
//	//rpc.RegisterSensorLedgerServiceServer(gs,&service.DefaultSensorLedgerService)
//	//rpc.RegisterVideoLedgerServiceServer(gs,&service.DefaultVideoLedgerService)
//	//rpc.RegisterUserLedgerServiceServer(gs,&service.DefaultUserService)
//	//rpc.RegisterQueryServiceServer(gs,&service.DefaultBlockService)
//	//reflection.Register(gs)
//	//启动服务
//	//fmt.Println("RPC服务 start listening .....")
//	//gs.Serve(ls)
//}

//func addDataBlock(ledgerService service.DataBlockService,keyStr[]string,value []byte){
//	//如果key的组成部分有三部分 说明检测到的是一个创世块
//	if	len(keyStr)==3&&(keyStr[2]==BLOCK_TENMINUTE||keyStr[2]==BLOCK_MINUTE||keyStr[2]==BLOCK_DAY) {
//		genesisblock:=&rpc.GenesisBlock{}
//		err:=json.Unmarshal(value,genesisblock)
//		if err!=nil{
//			logger.Error(err)
//			return
//		}
//		logger.Info("通过watch收到创世块数据:")
//		logger.Infof("%+v\n",genesisblock)
//		ledgerService.AddGenesisBlock(context.Background(),genesisblock)
//	}else if len(keyStr)>3{
//		switch keyStr[2] {
//		case BLOCK_DAY:
//			dayBlock:=&rpc.DailyDataBlock{}
//			json.Unmarshal(value,dayBlock)
//			if len(dayBlock.Blocks)==0{
//				return
//			}
//			logger.Info("通过watch收到天块数据:")
//			logger.Infof("%+v\n",dayBlock)
//			ledgerService.AddDailyBlock(context.Background(),dayBlock)
//		case BLOCK_MINUTE:
//			minuteBlock:=&rpc.MinuteDataBlock{}
//			json.Unmarshal(value,minuteBlock)
//			if len(minuteBlock.DataReceipts)==0{
//				return
//			}
//			logger.Info("通过watch收到分钟块数据:")
//			logger.Infof("%+v\n",minuteBlock)
//			ledgerService.AddMinuteBlock(context.Background(),minuteBlock)
//		case BLOCK_TENMINUTE:
//			tenMinuteBlock:=&rpc.TenMinuteDataBlock{}
//			json.Unmarshal(value,tenMinuteBlock)
//			if len(tenMinuteBlock.Blocks)==0{
//				return
//			}
//			logger.Info("通过watch收增强块数据:")
//			logger.Infof("%+v\n",tenMinuteBlock)
//			ledgerService.AddTenMinuteBlock(context.Background(),tenMinuteBlock)
//		default:
//			return
//		}
//	}
//}
//func addTxBlock(ledgerService service.TxBlockService,keyStr[]string,value []byte){
//	//如果key的组成部分有三部分 说明检测到的是一个创世块
//	if	len(keyStr)==3&&(keyStr[2]==BLOCK_TENMINUTE||keyStr[2]==BLOCK_MINUTE||keyStr[2]==BLOCK_DAY) {
//		logger.Info("通过Watch收到创世块   ..... ")
//		genesisblock:=&rpc.GenesisBlock{}
//		err:=json.Unmarshal(value,genesisblock)
//		if err!=nil{
//			logger.Error(err)
//			return
//		}
//		logger.Infof("receive data: %+v\n",genesisblock)
//		ledgerService.AddGenesisBlock(context.Background(),genesisblock)
//	}else if len(keyStr)>3{
//		//如果key的组成部分大于3部分 则说明是天块/分钟块/增强块
//		switch keyStr[2] {
//		case BLOCK_DAY:
//			dayBlock:=&rpc.DailyTxBlock{}
//			json.Unmarshal(value,dayBlock)
//			if len(dayBlock.Blocks)==0{
//				return
//			}
//			logger.Info("通过Watch收到天块   ..... ")
//			logger.Infof("%+v\n",dayBlock)
//			ledgerService.AddDailyBlock(context.Background(),dayBlock)
//		case BLOCK_MINUTE:
//			//logger.Info("通过Watch收到创世块   ..... ")
//			minuteBlock:=&rpc.MinuteTxBlock{}
//			json.Unmarshal(value,minuteBlock)
//			//如果分钟块中的交易数组为空 则返回 不处理
//			if len(minuteBlock.Transactions)==0{
//				return
//			}
//			logger.Info("通过Watch收到分钟块   ..... ")
//			logger.Infof("%+v\n",minuteBlock)
//			ledgerService.AddMinuteBlock(context.Background(),minuteBlock)
//		case BLOCK_TENMINUTE:
//			tenMinuteBlock:=&rpc.TenMinuteTxBlock{}
//			json.Unmarshal(value,tenMinuteBlock)
//			//如果增强块中的分钟块为空 则返回 不处理
//			if len(tenMinuteBlock.Blocks)==0{
//				return
//			}
//			logger.Info("通过Watch收到增强块   ..... ")
//			logger.Infof("%+v\n",tenMinuteBlock)
//			ledgerService.AddTenMinuteBlock(context.Background(),tenMinuteBlock)
//		default:
//			return
//		}
//	}
//}

//keyStr 是watch到的etcd中变化的key值
//value 是watch到的etcd中key对应的value值
//func routeHandler(keyStr[]string,value []byte){
//	//以交易类型和存证类型分成两个大类
//	switch keyStr[1]{
//	//VIDEO and USER are receipt type
//	case VIDEO:
//		ledgerService:= &service.DefaultVideoLedgerService
//		addDataBlock(ledgerService,keyStr,value)
//	case USER:
//		ledgerService:= &service.DefaultUserService
//		addDataBlock(ledgerService,keyStr,value)
//
//		//NODE and SENSOR and  ACCESS are tx type
//	case NODE:
//		ledgerService:= &service.DefaultNodeLedgerService
//		addTxBlock(ledgerService,keyStr,value)
//	case SENSOR:
//		ledgerService:= &service.DefaultSensorLedgerService
//		addTxBlock(ledgerService,keyStr,value)
//	case ACCESS:
//		ledgerService:= &service.DefaultAccessLedgerService
//		addTxBlock(ledgerService,keyStr,value)
//	default:
//	}
//}
//todo: 将本地区块文件转化成json格式
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: dialTimout,
	})
	if err != nil {
		logger.Fatal(err)
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	defer cli.Close()
	//因为天块数据和其他块数据类型的watch key值有差异，需要分开处理
	go func(){
		//启动watch机制，检测分钟块、增强块、创世块
		var now string=time.Now().Format("2006-01-02")
		for dayBlockChan := range cli.Watch(context.Background(), now, clientv3.WithPrefix()){
			for _, ev := range dayBlockChan.Events{
				keyStr := string(ev.Kv.Key)
				logger.Error("watch到变化的key为：",keyStr)
				logger.Error("watch到变化的value为：",string(ev.Kv.Value))
				if len(ev.Kv.Value)==0{
					logger.Error("value 是空的")
					break
				}
				strs := strings.Split(keyStr, ":")
				if len(strs) < 2 {
					logger.Error("watch 到数据key 格式不正确")
					break
				}
				if len(strs) == 4{
					//如果key的组成部分大于3部分 则说明是天块/分钟块/增强块
					//当检测到key的变化时去找对应的区块文件，如果文件已经收到了，那么判断文件内的数量与区块头记录的数量是否一致

					switch strs[1]{
					//VIDEO and USER are receipt type
					case VIDEO:
						logger.Info("here")
						addDataBlock(keyStr,strs,ev.Kv.Value)
						//case USER:
						//	ledgerService:= &service.DefaultUserService
						//	addDataBlock(ledgerService,keyStr,value)
						//
						//	//NODE and SENSOR and  ACCESS are tx type
						//case NODE:
						//	ledgerService:= &service.DefaultNodeLedgerService
						//	addTxBlock(ledgerService,keyStr,value)
						//case SENSOR:
						//	ledgerService:= &service.DefaultSensorLedgerService
						//	addTxBlock(ledgerService,keyStr,value)
						//case ACCESS:
						//	ledgerService:= &service.DefaultAccessLedgerService
						//	addTxBlock(ledgerService,keyStr,value)
						//default:
					}
				}
			}
			now =time.Now().Format("2006-01-02")
		}
	}()
	go func() {
		//启动watch机制，watch 天块
		var preday string=(time.Now().Add(-time.Hour*24)).Format("2006-01-02")
		for nonDayblockChan := range cli.Watch(context.Background(), preday, clientv3.WithPrefix()) {
			for _, ev := range nonDayblockChan.Events {
				keyStr := string(ev.Kv.Key)
				if len(ev.Kv.Value)==0{
					logger.Error("value 是空的")
					break
				}
				strs := strings.Split(keyStr, ":")
				if len(strs) < 2 {
					logger.Error("watch key 数据不正确")
					break
				}


				//			switch strs[2] {
				//			case BLOCK_DAY:
				//				//dayBlock:=&rpc.DailyTxBlock{}
				//				//json.Unmarshal(ev.Kv.Value,dayBlock)
				//				////
				//				//if len(dayBlock.Blocks)==0{
				//				//	return
				//				//}
				//				//logger.Info("通过Watch收到天块   ..... ")
				//				//logger.Infof("%+v\n",dayBlock)
				//				//ledgerService.AddDailyBlock(context.Background(),dayBlock)
				//			case BLOCK_MINUTE:
				//				//logger.Info("通过Watch收到创世块   ..... ")
				//				minuteBlock:=&rpc.MinuteTxBlock{}
				//				json.Unmarshal(ev.Kv.Value,minuteBlock)
				//				//如果分钟块中的交易数组为空 则返回 不处理
				//				if len(minuteBlock.Transactions)==0{
				//					return
				//				}
				//				logger.Info("通过Watch收到分钟块   ..... ")
				//				logger.Infof("%+v\n",minuteBlock)
				//				ledgerService.AddMinuteBlock(context.Background(),minuteBlock)
				//			case BLOCK_TENMINUTE:
				//				tenMinuteBlock:=&rpc.TenMinuteTxBlock{}
				//				json.Unmarshal(value,tenMinuteBlock)
				//				//如果增强块中的分钟块为空 则返回 不处理
				//				if len(tenMinuteBlock.Blocks)==0{
				//					return
				//				}
				//				logger.Info("通过Watch收到增强块   ..... ")
				//				logger.Infof("%+v\n",tenMinuteBlock)
				//}


			//	routeHandler(strs, ev.Kv.Value)
			}
			preday =(time.Now().Add(-time.Hour*24)).Format("2006-01-02")
		}
	}()
	//if datasource.ChainConfig.P2p.Local.TurnOn {
	//	go sync.P2P()
	//}
	var hostname string

	flag.StringVar(&hostname, "name", "leader", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	//for name,val:=range datasource.ChainConfig.Common.LedgerName{
	//	if val.Leader==hostname{
	//		datasource.SupportLedger[name]=true
	//	}else {
	//		for _,name:=range val.Follower{
	//			if name==hostname{
	//				datasource.SupportLedger[name]=true
	//			}
	//		}
	//	}
	//}
	//var serviceAddress string
	//for {
	//	if val, ok := datasource.ChainConfig.Consensus.EtcdGroup[hostname]; !ok {
	//		fmt.Fprintf(os.Stderr, "输入的主机名有误,请重新输入")
	//		fmt.Printf("请输入当前节点的主机名: ")
	//		fmt.Scanf("%s",&hostname)
	//	} else {
	//		serviceAddress=val.BlockAddress
	//		break
	//	}
	//}
	//创建tcp监听
	//addr:=strings.Split(serviceAddress,":")
	//if len(addr)!=2{
	//if len(addr)!=2{
	//	logger.Error("节点的block_grpc配置项有误")
	//	os.Exit(1)
	//}
	//port:=addr[1]
	//go rpcService(port)
	////go etcdDiscov()
	//go RegisterService(cli)
	select {}
}
func AddDataBlockToTdengine(minuteDataBlock *dataStruct.MinuteDataBlock)(int64,error){
	//提取区块中的存证记录，存到tdengine数据库中
	//				Tdengine为了 一设备一表的设计架构设计了超级表-子表的概念
	//我们区块链节点设计成一个设备一张表  使用超表的模板来生成每一个设备的子表
	start := time.Now().UnixMilli()
	logger.Error(minuteDataBlock.DataReceipts[0])
	//heightStr:=strconv.FormatInt(block.Header.BlockHeight,10)
	date := util.ParseDate1(minuteDataBlock.Header.CreateTimestamp)
	heightstr := strconv.FormatInt(minuteDataBlock.Header.BlockHeight,10)
	blockIdentify  :=date +"_"+minuteDataBlock.Header.LedgerType+"_"+"_"+heightstr
	//分钟快的表名： 链类型_区块类型_日期_区块高度_receipt/transaction
	blockTableName :=minuteDataBlock.Header.LedgerType+"_"+minuteDataBlock.Header.BlockType+"_"+date +"_"+heightstr+"_receipt"

	str := "insert into testdb."+blockTableName+" using testdb.mockdatareceipt tags('"+blockIdentify+"','"+date+"') values "
	logger.Error("本次要插入数据条数：",len(minuteDataBlock.DataReceipts))
	//logger.Error("sql语句：",str)
	for _,data:=range minuteDataBlock.DataReceipts{
		//写到表中
		//logger.Warn("***AddDataBlock函数中获取到的原生存证数据为:",receipt)

		//heightStr:=strconv.FormatInt(block.Header.BlockHeight,10)
		//这两个是tags要用的
		//mockReceipt.BlockIdentify=prefix+"_"+heightStr
		//mockReceipt.Date=util.ParseDate(mockReceipt.CreateTimestamp)
		//尝试一条数据的preare插入数据库
		//数据库中表为18项

		values := "('"+data.CreateTimeStamp +"','"+data.EntityId+"','"+data.KeyId+"','"+strconv.FormatFloat(data.ReceiptValue,'E',-1,64)+"','"+data.Version+
			"','"+data.UserName+"','"+data.OperationType+"','"+data.DataType+"','"+data.ServiceType+"','"+data.FileName+"','"+strconv.FormatFloat(data.FileSize,'E',-1,64)+"','"+data.FileHash+
			"','"+data.Uri+"','"+data.ParentKeyId+"',"+"'a'"+",'"+data.AttachmentTotalHash+"')"
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
//处理添加存证记录的逻辑
func addDataBlock(keyStr string,keyStrs []string, blockHeader []byte) {
	//如果key的组成部分有三部分 说明检测到的是一个创世块
	if	len(keyStrs)==3&&(string(keyStrs[2])==BLOCK_TENMINUTE||string(keyStrs[2])==BLOCK_MINUTE||string(keyStrs[2])==BLOCK_DAY) {
		//genesisblock:=&rpc.GenesisBlock{}
		//err:=json.Unmarshal(value,genesisblock)
		//if err!=nil{
		//	logger.Error(err)
		//	return
		//}
		//logger.Info("通过watch收到创世块数据:")
		//logger.Infof("%+v\n",genesisblock)
		//ledgerService.AddGenesisBlock(context.Background(),genesisblock)
	}else if len(keyStrs)>3{
		switch keyStrs[2]{
		case BLOCK_DAY:
			dayBlockHeader:=&rpc.BlockHeader{}
			json.Unmarshal(blockHeader,dayBlockHeader)

			logger.Info("通过watch收到天块数据:")
			//logger.Infof("%+v\n",dayBlock)
			//ledgerService.AddDailyBlock(context.Background(),dayBlock)
		case BLOCK_MINUTE:
			//因为暂时hraft是存储去块体的，所以传过来的数据其实是区块头jia qukuaiti 
			minuteBlock:=&dataStruct.MinuteDataBlock{}
			minuteBlockHeader := &dataStruct.BlockHeader{}
			error := json.Unmarshal(blockHeader,minuteBlockHeader)
			//把区块头赋值给区块头指针变量
			minuteBlock.Header = minuteBlockHeader
			if error!=nil{
				fmt.Println("反序列化出错：",error)
				return
			}
			logger.Info("通过watch收到分钟块数据:",string(blockHeader))
			logger.Infof("%+v\n",minuteBlock.Header)
			//如果分钟快中没有记录的话就不做处理
			if minuteBlock.Header.CurrentDataCount==0{
				return
			}else if minuteBlock.Header.CurrentDataCount>0{
				//1.找对应的区块文件
				//1.1 拼接出目标区块文件名字
				const TIME_LAYOUT = "2006-01-02 15:04:05"
				_,err:=time.Parse(TIME_LAYOUT,minuteBlock.Header.CreateTimestamp)
				if err!=nil{
					panic(err)
				}
				//blockPath :=  "/root/"+config.Initialize().ChainName+"/"+strconv.FormatInt(int64(time.Year()),10)+"-"+strconv.FormatInt(int64(time.Month()),10)+"-"+strconv.FormatInt(int64(time.Day()),10)+"/"+keyStrs[1]+"/"+
				//	keyStrs[2]+"/"+keyStrs[3]
				blockPath :=  "/root/"+config.Initialize().ChainName+"/2021-10-26/video/MINUTE/915"
				logger.Info("目标区块文件目录为：",blockPath)

				//1.2打开区块文件
				var file *os.File
				var rawBlockdata []byte
				rawBlockdata = make([]byte,10240)
				var dataReceipts []*dataStruct.DataReceipt
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
			err = json.Unmarshal(rawBlockdata[:n], &dataReceipts)
			if err != nil{
				logger.Error("从区块文件反序列化失败：",err)
				return
			}
			logger.Info("从区块文件反序列化成功：",dataReceipts)
			//插入tdengine数据库
			minuteBlock.DataReceipts = dataReceipts
				AddDataBlockToTdengine(minuteBlock)
			}

			//ledgerService.AddMinuteBlock(context.Background(),minuteBlock)
		//case BLOCK_TENMINUTE:
		//	tenMinuteBlock:=&rpc.TenMinuteDataBlock{}
		//	json.Unmarshal(value,tenMinuteBlock)
		//	if len(tenMinuteBlock.Blocks)==0{
		//		return
		//	}
		//	logger.Info("通过watch收增强块数据:")
		//	logger.Infof("%+v\n",tenMinuteBlock)
		//	ledgerService.AddTenMinuteBlock(context.Background(),tenMinuteBlock)
		//default:
		//	return
		}
	}
}
